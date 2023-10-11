package appleiap

type ProductType string

const (
	Consumable                ProductType = "Consumable"                  //消耗型产品
	NonConsumable             ProductType = "Non-Consumable"              //非消耗型产品
	AutoRenewableSubscription ProductType = "Auto-Renewable Subscription" //连续订阅产品
	NonRenewingSubscription   ProductType = "Non-Renewing Subscription"   //非连续订阅产品
)

type InAppOwnershipType string

const (
	FamilyShared InAppOwnershipType = "AMILY_SHARED" //家庭共享
	Purchased    InAppOwnershipType = "PURCHASED"    //仅限本人
)

type Environment string

const (
	Sandbox    Environment = "Sandbox"    //沙箱
	Production Environment = "Production" // 生产环境
)

type AutoRenewStatus int

const (
	AutoRenewOff AutoRenewStatus = iota //自动续费状态关
	AutoRenewOn                         //自动续费状态开
)

// ExpirationIntent https://developer.apple.com/documentation/appstoreservernotifications/expirationintent
type ExpirationIntent int

const (
	CustomerCanceled    ExpirationIntent = iota + 1 //客户取消了订阅。
	BillingError                                    //计费错误；例如，客户的付款信息不再有效。
	CustomerRejected                                //客户不同意需要客户同意的自动续订订阅价格上涨，从而导致订阅过期。
	ProductNotAvailable                             //该产品在续订时无法购买。
)

// OfferType https://developer.apple.com/documentation/appstoreservernotifications/offertype
type OfferType int

const (
	IntroductoryOffer OfferType = iota + 1 //介绍性报价
	PromotionalOffer                       //促销优惠
	OfferCode                              //带有订阅优惠代码的优惠
)

// PriceIncreaseStatus https://developer.apple.com/documentation/appstoreservernotifications/priceincreasestatus
type PriceIncreaseStatus int

const (
	CustomerUncertain PriceIncreaseStatus = iota //客户尚未对需要客户同意的自动续订订阅价格上涨做出回应。
	CustomerConsented                            //客户同意自动续订订阅价格上涨需要客户同意，或者 App Store 已通知客户自动续订订阅价格上涨不需要客户同意。

)

// RevocationReason https://developer.apple.com/documentation/appstoreserverapi/revocationreason
type RevocationReason int

const (
	CustomerReasons RevocationReason = iota //客户原因
	AppHasIssue                             //app存在问题
)
