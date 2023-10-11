package appleiap

// NotificationType https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
type NotificationType string

const (
	NotificationConsumptionRequest     NotificationType = "CONSUMPTION_REQUEST"       //表示客户针对应用内购买的消耗品发起了退款请求，并且 App Store 正在要求您提供消耗数据。有关详细信息，请参阅发送消耗信息。
	NotificationDidChangeRenewalPref   NotificationType = "DID_CHANGE_RENEWAL_PREF"   //通知类型及其子类型指示用户对其订阅计划进行了更改。如果子类型是 UPGRADE，则用户升级了他们的订阅。升级立即生效，开始新的计费周期，用户收到前一周期未使用部分的按比例退款。如果子类型是 DOWNGRADE，则用户降级或交叉分级他们的订阅。降级在下次续订时生效。当前有效的计划不受影响。如果子类型为空，则用户将其续订偏好更改回当前订阅，从而有效地取消了降级。
	NotificationDidChangeRenewalStatus NotificationType = "DID_CHANGE_RENEWAL_STATUS" //通知类型及其子类型指示用户对订阅续订状态进行了更改。如果子类型为 AUTO_RENEW_ENABLED，则用户重新启用订阅自动续订。如果子类型为 AUTO_RENEW_DISABLED，则用户禁用订阅自动续订，或 App Store 在用户请求退款后禁用订阅自动续订。
	NotificationDidFailToRenew         NotificationType = "DID_FAIL_TO_RENEW"         //通知类型及其子类型指示由于计费问题而未能更新订阅。订阅进入计费重试周期。如果子类型为GRACE_PERIOD，在宽限期内继续提供服务。如果子类型为空，表示订阅未处于宽限期，可以停止提供订阅服务。 通知用户他们的账单信息可能有问题。App Store会在60天内继续重试计费，或者直到用户解决了他们的计费问题或取消了他们的订阅。有关更多信息，请参见减少非自愿用户流失。
	NotificationDidRenew               NotificationType = "DID_RENEW"                 //通知类型及其子类型指示订阅已成功更新。如果子类型是BILLING_RECOVERY，以前未能续订的过期订阅现在已成功续订。如果子状态为空，则活动订阅已成功自动更新为新的事务周期。向客户提供对订阅内容或服务的访问权。
	NotificationExpired                NotificationType = "EXPIRED"                   //通知类型及其子类型指示订阅已过期。如果子类型是自愿的，则订阅在用户禁用订阅续订后过期。如果子类型是BILLING_RETRY，则订阅已过期，因为计费重试周期结束时没有成功的计费事务。如果子类型是PRICE_INCREASE，则订阅过期是因为用户不同意价格上涨。
	NotificationGracePeriodExpired     NotificationType = "GRACE_PERIOD_EXPIRED"      //指示计费宽限期已结束，无需更新订阅，因此您可以关闭对服务或内容的访问。通知用户他们的账单信息可能有问题。App Store会在60天内继续重试计费，或者直到用户解决了他们的计费问题或取消了他们的订阅。有关更多信息，请参见减少非自愿用户流失。
	NotificationOfferRedeemed          NotificationType = "OFFER_REDEEMED"            //通知类型及其子类型表明用户兑换了促销优惠或优惠代码。 如果子类型是INITIAL_BUY，则用户将该优惠兑换为首次购买。如果子类型为RESUBSCRIBE，则用户赎回要约以重新订阅未激活的订阅。如果子类型为UPGRADE，则用户兑现了一个优惠以升级其活动订阅，该优惠将立即生效。如果子类型为降级，则用户兑现了一个报价，以降级其在下一次续订日期生效的活动订阅。如果用户为其活动订阅兑现了要约，您将收到一个没有子类型的offer_redemption通知类型。有关促销优惠的更多信息，请参见在您的应用程序中实现促销优惠。有关优惠代码的更多信息，请参见在您的应用程序中实现优惠代码。
	NotificationPriceIncrease          NotificationType = "PRICE_INCREASE"            //通知类型及其子类型表明系统已通知客户订阅价格上涨。如果子类型是PENDING，客户还没有对价格上涨做出回应。如果子类型为ACCEPTED，则客户已经接受了价格上涨。 有关系统如何在应用程序运行时通知客户价格上涨的信息，请参阅paymentQueueShouldShowPriceConsent(_:)。
	NotificationRefund                 NotificationType = "REFUND"                    //表示App Store成功退还了可消费的应用内购买、不可消费的应用内购买、自动更新订阅或不可更新订阅的交易。 revocationDate包含被退回事务的时间戳。originalTransactionId和productId标识原始事务和产品。revocationReason中包含了原因。要为用户请求所有退款交易的列表，请参见App Store Server API中的Get Refund History。
	NotificationRefundDeclined         NotificationType = "REFUND_DECLINED"           //表示 App Store 拒绝了应用开发者发起的退款请求。
	NotificationRenewalExtended        NotificationType = "RENEWAL_EXTENDED"          //表示 App Store 延长了开发者要求的订阅续订日期。
	NotificationRevoke                 NotificationType = "REVOKE"                    //表示用户通过“家庭分享”获得的应用内购买不再通过“分享”获得。当购买者禁止家庭共享产品，购买者(或家庭成员)离开家庭组，或者购买者要求并收到退款时，App Store会发送此通知。你的应用程序也接收一个paymentQueue(_: didrevokeentitlements sforproductidentifiers:)调用。Family Sharing适用于非消耗性应用内购买和自动更新订阅。有关家庭共享的更多信息，请参阅支持家庭共享在您的应用程序。
	NotificationSubscribed             NotificationType = "SUBSCRIBED"                //通知类型及其子类型指示用户订阅了某个产品。如果子类型是INITIAL_BUY，则用户第一次通过Family Sharing购买或接受对订阅的访问。如果子类型为RESUBSCRIBE，则用户通过“家庭共享”重新订阅或接受对同一订阅或同一订阅组内的另一个订阅的访问。
)

// Subtype https://developer.apple.com/documentation/appstoreservernotifications/subtype
type Subtype string

const (
	SubTypeInitialBuy        Subtype = "INITIAL_BUY"         //适用于 SUBSCRIBED 通知类型。具有此子类型的通知表示用户首次购买订阅，或者用户首次通过家庭共享获得订阅访问权限。
	SubTypeResubscribe       Subtype = "RESUBSCRIBE"         //适用于 SUBSCRIBED 通知类型。具有此子类型的通知表示用户重新订阅或通过家庭共享接收对同一订阅或同一订阅组中的另一个订阅的访问权限。
	SubTypeDowngrade         Subtype = "DOWNGRADE"           //适用于 DID_CHANGE_RENEWAL_PREF 通知类型。具有此子类型的通知表明用户降级了他们的订阅。降级在下次续订时生效。
	SubTypeUpgrade           Subtype = "UPGRADE"             //适用于 DID_CHANGE_RENEWAL_PREF 通知类型。具有此子类型的通知表明用户升级了他们的订阅。升级立即生效。
	SubTypeAutoRenewEnabled  Subtype = "AUTO_RENEW_ENABLED"  //适用于 DID_CHANGE_RENEWAL_STATUS 通知类型。具有此子类型的通知表明用户启用了订阅自动续订。
	SubTypeAutoRenewDisabled Subtype = "AUTO_RENEW_DISABLED" //适用于 DID_CHANGE_RENEWAL_STATUS 通知类型。具有此子类型的通知表示用户禁用订阅自动续订，或 App Store 在用户请求退款后禁用订阅自动续订。
	SubTypeVoluntary         Subtype = "VOLUNTARY"           //适用于 EXPIRED 通知类型。具有此子类型的通知指示订阅在用户禁用订阅自动续订后过期。
	SubTypeBillingRetry      Subtype = "BILLING_RETRY"       //适用于 EXPIRED 通知类型。具有此子类型的通知指示订阅已过期，因为订阅未能在计费重试期结束之前续订。
	SubTypePriceIncrease     Subtype = "PRICE_INCREASE"      //适用于 EXPIRED 通知类型。具有此子类型的通知表明订阅已过期，因为用户不同意涨价。
	SubTypeGracePeriod       Subtype = "GRACE_PERIOD"        //适用于 DID_FAIL_TO_RENEW 通知类型。带有此子类型的通知表明订阅由于计费问题而无法续订；在宽限期内继续提供对订阅的访问。
	SubTypeBillingRecovery   Subtype = "BILLING_RECOVERY"    //适用于 DID_RENEW 通知类型。具有此子类型的通知表示以前无法续订的过期订阅现在已成功续订。
	SubTypePending           Subtype = "PENDING"             //适用于 PRICE_INCREASE 通知类型。此子类型的通知表示系统通知用户订阅价格上涨，但用户尚未接受。
	SubTypeAccepted          Subtype = "ACCEPTED"            //适用于 PRICE_INCREASE 通知类型。具有此子类型的通知表示用户接受了订阅价格上涨。

)

type IapData struct {
	NotificationType NotificationType   `json:"notification_type"`
	Subtype          Subtype            `json:"subtype"`
	Transaction      *TransactionClaims `json:"transaction"`
	Renewal          *RenewalClaims     `json:"renewal"`
}

// TransactionClaims https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type TransactionClaims struct {
	BundleID                    string             `json:"bundleId"`
	Environment                 Environment        `json:"environment"`
	ExpiresDate                 int64              `json:"expiresDate"`
	InAppOwnershipType          InAppOwnershipType `json:"inAppOwnershipType"`
	OriginalPurchaseDate        int64              `json:"originalPurchaseDate"`
	OriginalTransactionId       string             `json:"originalTransactionId"`
	ProductId                   string             `json:"productId"`
	PurchaseDate                int64              `json:"purchaseDate"`
	Quantity                    int64              `json:"quantity"`
	SignedDate                  int64              `json:"signedDate"` //App Store 对 JSON Web 签名数据进行签名的 UNIX 时间（以毫秒为单位）
	SubscriptionGroupIdentifier string             `json:"subscriptionGroupIdentifier"`
	TransactionId               string             `json:"transactionId"`
	Type                        ProductType        `json:"type"`
	WebOrderLineItemId          string             `json:"webOrderLineItemId"`
	//以下支付回调无法获得
	IsUpgraded       bool             `json:"isUpgraded"`
	OfferIdentifier  string           `json:"offerIdentifier"`  //包含促销代码或促销优惠标识符的标识符。
	OfferType        OfferType        `json:"offerType"`        //订阅优惠的类型。
	RevocationDate   int64            `json:"revocationDate"`   //Apple 支持退还交易的 UNIX 时间（以毫秒为单位）。
	RevocationReason RevocationReason `json:"revocationReason"` //退款交易的原因
}

// RenewalClaims https://developer.apple.com/documentation/appstoreservernotifications/jwsrenewalinfodecodedpayload
type RenewalClaims struct {
	OriginalTransactionId string           `json:"originalTransactionId"` //购买的原始交易标识符。
	AutoRenewProductId    string           `json:"autoRenewProductId"`    //在下一个计费周期续订的产品的产品标识符。
	ProductId             string           `json:"productId"`             //应用内购买的产品标识符。
	AutoRenewStatus       AutoRenewStatus  `json:"autoRenewStatus"`       //自动续订订阅的续订状态。
	SignedDate            int64            `json:"signedDate"`            //App Store 对 JSON Web 签名数据进行签名的 UNIX 时间（以毫秒为单位）
	Environment           Environment      `json:"environment"`           //服务器环境，沙盒或生产环境
	ExpirationIntent      ExpirationIntent `json:"expirationIntent"`      //订阅过期的原因
	//以下支付回调无法获得
	GracePeriodExpiresDate string              `json:"gracePeriodExpiresDate"` //订阅续订的计费宽限期到期的时间
	IsInBillingRetryPeriod bool                `json:"isInBillingRetryPeriod"` //指示 App Store 是否正在尝试自动续订因计费问题而过期的订阅。
	OfferIdentifier        string              `json:"offerIdentifier"`        //包含促销代码或促销优惠标识符的标识符。
	OfferType              OfferType           `json:"offerType"`              //订阅优惠的类型。
	PriceIncreaseStatus    PriceIncreaseStatus `json:"priceIncreaseStatus"`    //指示自动续订订阅是否会涨价的状态。
}
