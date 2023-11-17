package wechat

import (
	"crypto/sha1"
	"fmt"
)

func JSTicketSignature(ts, nonceStr, url string, ticket JSTicket) (string, error) {
	s1 := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, ts, url)
	h := sha1.New()
	if _, err := h.Write([]byte(s1)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
