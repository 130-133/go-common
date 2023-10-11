package appleiap

type LookUpResp struct {
	Status             int32    `json:"status"`
	SignedTransactions []string `json:"signedTransactions"`
}

type LookUpData struct {
	TransactionId         string             `json:"transactionId"`
	OriginalTransactionId string             `json:"originalTransactionId"`
	BundleId              string             `json:"bundleId"`
	ProductId             string             `json:"productId"`
	PurchaseDate          int64              `json:"purchaseDate"`
	OriginalPurchaseDate  int64              `json:"originalPurchaseDate"`
	Quantity              int32              `json:"quantity"`
	Type                  string             `json:"type"`
	InAppOwnershipType    InAppOwnershipType `json:"inAppOwnershipType"`
	SignedDate            int64              `json:"signedDate"`
}
