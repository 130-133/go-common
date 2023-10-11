package appleiap

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt/v4"
)

type JwsParse struct {
	publicCert string
}

// NewJws jws解析
func NewJws(publicCert string) *JwsParse {
	if publicCert == "" {
		publicCert = "AppleRootCA-G3.pem"
	}
	return &JwsParse{
		publicCert: publicCert,
	}
}

// IapParseNotify 解析服务器通知信息
func (n *JwsParse) IapParseNotify(signedPayload string) (*IapData, error) {
	claims, err := n.IapParse(signedPayload)
	if err != nil {
		return nil, err
	}
	notificationType, _ := claims["notificationType"].(string)
	subtype, _ := claims["subtype"].(string)

	data := claims["data"].(map[string]interface{})
	signedTransactionToken, err := jwt.Parse(data["signedTransactionInfo"].(string), n.checkToken)
	if err != nil {
		return nil, err
	}
	signedRenewalToken, err := jwt.Parse(data["signedRenewalInfo"].(string), n.checkToken)
	if err != nil {
		return nil, err
	}
	mapTransaction, ok := signedTransactionToken.Claims.(jwt.MapClaims)
	if !(ok && signedTransactionToken.Valid) {
		return nil, errors.New("无效token")
	}
	transaction := TransactionClaims{}
	if by, err := json.Marshal(mapTransaction); err == nil {
		_ = json.Unmarshal(by, &transaction)
	}

	mapRenewal, ok := signedRenewalToken.Claims.(jwt.MapClaims)
	if !(ok && signedRenewalToken.Valid) {
		return nil, errors.New("无效token")
	}
	renewal := RenewalClaims{}
	if by, err := json.Marshal(mapRenewal); err == nil {
		_ = json.Unmarshal(by, &renewal)
	}
	return &IapData{
		NotificationType: NotificationType(notificationType),
		Subtype:          Subtype(subtype),
		Transaction:      &transaction,
		Renewal:          &renewal,
	}, nil
}

// IapParseLookUp 解析查询用户订单
func (n *JwsParse) IapParseLookUp(signedPayload string) (*LookUpData, error) {
	claims, err := n.IapParse(signedPayload)
	if err != nil {
		return nil, err
	}
	data := LookUpData{}
	if by, err := json.Marshal(claims); err == nil {
		_ = json.Unmarshal(by, &data)
	}
	return &data, nil
}

// IapParse 解析jws
func (n *JwsParse) IapParse(signedPayload string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedPayload, n.checkToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("无效token")
	}
	return claims, nil
}

func (n *JwsParse) checkToken(token *jwt.Token) (interface{}, error) {
	// header alg: ES256
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
	}
	// header x5c: ["A","B","C"]
	x5c, ok := token.Header["x5c"]
	if !ok {
		return nil, errors.New("header x5c need set")
	}

	x5cArr, ok := x5c.([]interface{})
	if !ok {
		return nil, errors.New("header x5c is array")
	}

	var x5cArrStr []string
	for _, v := range x5cArr {
		s, ok := v.(string)
		if !ok {
			return nil, errors.New("header x5c is string array")
		}
		x5cArrStr = append(x5cArrStr, s)
	}
	// 校验证书有效性，并返回公钥
	publicKey, err := n.checkCerts(x5cArrStr)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (n *JwsParse) checkCerts(x5cField []string) (interface{}, error) {
	var pems []string
	for _, x5c := range x5cField {
		pemData := "-----BEGIN CERTIFICATE-----\n"
		for i := 0; i < len(x5c); i += 64 {
			end := i + 64
			if end > len(x5c) {
				end = len(x5c)
			}
			pemData += x5c[i:end] + "\n"
		}
		pemData += "-----END CERTIFICATE-----"
		pems = append(pems, pemData)
	}

	var certs []*x509.Certificate
	// https://www.apple.com/certificateauthority/AppleRootCA-G3.cer
	// openssl x509 -inform der -in AppleRootCA-G3.cer -out AppleRootCA-G3.pem
	rootPem, err := ioutil.ReadFile(n.publicCert)
	if err != nil {
		return nil, err
	}
	pems = append(pems, string(rootPem))

	for _, pemData := range pems {
		// Parse PEM block
		var block *pem.Block
		if block, _ = pem.Decode([]byte(pemData)); block == nil {
			return nil, errors.New("invalid pem format")
		}

		// Parse the key
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		return nil, errors.New("invalid x5c")
	}
	// 校验证书链
	for i := 0; i+1 < len(certs); i++ {
		if err := certs[i].CheckSignatureFrom(certs[i+1]); err != nil {
			return nil, err
		}
	}

	/*publicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(cer))
	  if err != nil {
	      return err
	  }*/
	return certs[0].PublicKey, nil
}
