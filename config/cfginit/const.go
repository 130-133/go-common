package cfginit

const (
	VIP_TYPE_FOREVER_CARD    = iota + 1 //终身会员
	VIP_TYPE_YEAR_CARD                  //年会员
	VIP_TYPE_EXPERIENCE_CARD            //体验卡
	VIP_TYPE_MONTH_CARD                 //月会员
	VIP_TYPE_SUBMONTH_CARD              //订阅月会员
)

const (
	GOODS_TYPE_COURSE_PACKAGE = iota + 4 //课包商品
	GOODS_TYPE_WORK_ADAPT                //作品改编
	GOODS_TYPE_SKIN_MODEL                //皮肤模块
)

const (
	GOODS_REF_GROUP_SPECIAL    = iota + 1 //专项课
	GOODS_REF_GROUP_SYSTEM                //系统课
	GOODS_REF_GROUP_CPA                   //CPA
	GOODS_REF_GROUP_WORK_ADAPT            //作品改编
	GOODS_REF_GROUP_SKIN_MODEL            //皮肤模块
)

// 商品ID分类
const (
	GOODS_REF_TYPE_SPECIAL_SUBJECT = iota + 1 //专项课课程
	GOODS_REF_TYPE_SYSTEM_SUBJECT             //系统课课程
	GOODS_REF_TYPE_SPECIAL_THEME              //专项课主题
	GOODS_REF_TYPE_SYSTEM_THEME               //系统课主题
	GOODS_REF_TYPE_CPA                        //cpa表
	GOODS_REF_TYPE_LESSON                     //课时
	GOODS_REF_TYPE_WORK                       //作品
	GOODS_REF_TYPE_MODEL                      //皮肤
)

const (
	SHOP_TRADE_STAGE_WAIT_PAY = iota + 1
	SHOP_TRADE_STAGE_PAID     = 2
	SHOP_TRADE_STAGE_FAILD
)
