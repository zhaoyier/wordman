package wrapx

import (
	"reflect"
	"sync"
	"github.com/gin-gonic/gin"
)

// NewAPIFunc Custom context support
type NewAPIFunc2 func(*gin.Context) interface{}

var _genInfo genInfo
var _genInfoCnf genInfo
var _mp sync.Map
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// 错误消息类型
type ErrCode int

const ( //消息id定义
	NormalMessageID   ErrCode = 0 // 默认返回值
	ServerMaintenance ErrCode = 1 // 服务器维护中 请稍后再试
	AccountDisabled   ErrCode = 2 // 帐号被禁用
	AppidOverdue      ErrCode = 3 // appid过期

	UnknownError  ErrCode = 101 // 未知错误
	TokenFailure  ErrCode = 102 // token失效
	HTMLSuccess   ErrCode = 200 // 成功
	BlockingAcess ErrCode = 405 // 禁止访问

	NewReport ErrCode = 2001 // 新消息
	NewHeart  ErrCode = 2002 // 心跳

	ParameterInvalid          ErrCode = 1001 // 参数无效
	AppidParameterInvalid     ErrCode = 1002 // appid参数无效
	EncryptCheckError         ErrCode = 1003 // 密文校验失败,aa
	UserNameDoNotExist        ErrCode = 1004 // 用户名不存在或密码错误
	DuplicateKeyError         ErrCode = 1005 // 键值对重复
	InValidOp                 ErrCode = 1007 // 无效操作
	NotFindError              ErrCode = 1006 // 未找到
	InValidAuthorize          ErrCode = 1008 // 授权码错误
	HasusedError              ErrCode = 1009 // 已被使用
	HasActvError              ErrCode = 1010 // 已被激活
	ActvFailure               ErrCode = 1011 // 激活码被禁止使用
	UserExisted               ErrCode = 1012 // 用户已存在
	VerifyTimeError           ErrCode = 1013 // 验证码请求过于平凡
	MailSendFaild             ErrCode = 1014 // 邮箱发送失败
	SMSSendFaild              ErrCode = 1015 // 手机发送失败
	PhoneParameterError       ErrCode = 1016 // 手机号格式有问题
	VerifyError               ErrCode = 1017 // 验证码错误
	UserNotExisted            ErrCode = 1018 // 用户不存在
	TopicExisted              ErrCode = 1019 // topic已经存在
	TopicNotExisted           ErrCode = 1020 // topic不存在
	BundleIDNotExisted        ErrCode = 1021 // bundle_id不存在
	TopicStartFail            ErrCode = 1022 // topic开启处理失败
	TopicTypeNotExisted       ErrCode = 1023 // topic处理类型不存在
	TopicIsNotNull            ErrCode = 1024 // topic不能为空
	DeviceNotExisted          ErrCode = 1025 // 设备不存在
	StateExisted              ErrCode = 1027 // 状态已存在
	LastMenuNotExisted        ErrCode = 1028 // 上级菜单不存在
	MenuNotExisted            ErrCode = 1029 // 菜单不存在
	UserMenuNotExisted        ErrCode = 1030 // 用户权限不存在
	DeviceIDNotExisted        ErrCode = 1031 // 设备ID不存在
	GoodsDealTypeNotExisted   ErrCode = 1032 // 商品处理类型不存在
	GoodsIDNotExisted         ErrCode = 1033 // 商品不存在
	GoodsBeInDiscount         ErrCode = 1034 // 商品正在打折
	GoodsPayTypeNotExisted    ErrCode = 1035 // 商品可支付类型不存在
	GoodsIDExisted            ErrCode = 1036 // 商品已存在
	OrderIDNotExisted         ErrCode = 1043 // 订单不存在
	GoodsBeNotInDiscount      ErrCode = 1044 // 商品未打折
	NotifyIsNotMatch          ErrCode = 1045 // 会话不匹配
	GoodsIsDiscountRecovery   ErrCode = 1046 // 商品已恢复原价
	InvitationUserNotExisted  ErrCode = 1047 // 邀请用户不存在
	InvitationUserLevelIsFull ErrCode = 1048 // 邀请用户级数已满
	UserNotAuthorize          ErrCode = 1049 // 用户未授权
	ApplicantIsExisted        ErrCode = 1050 // 申请人已存在
	ApplicantNotExisted       ErrCode = 1051 // 申请人不存在
	CreditOrderNotVaild       ErrCode = 1052 // 订单无效
	RepeatWxWithdraw          ErrCode = 1053 // 微信零钱重复提现
	WxWithdrawAmountError     ErrCode = 1054 // 提现金额错误
	WxWithdrawError           ErrCode = 1055 // 微信提现失败
	RepeatSubmission          ErrCode = 1056 // 重复提交
	BundleExisted             ErrCode = 1057 // bundle已存在
	AuthExisted               ErrCode = 1058 // 权限已存在
	AuthNotExisted            ErrCode = 1059 // 权限不存在
	RoomTypeNotExisted        ErrCode = 1060 // 房型不存在
	RoomTypeExisted           ErrCode = 1061 // 房型已存在
	RoomNoNotExisted          ErrCode = 1062 // 房间不存在
	RoomNoExisted             ErrCode = 1063 // 房间已存在
	RateTypeExisted           ErrCode = 1064 // 房价代码或房价名称已存在
	RateTypeNotExisted        ErrCode = 1065 // 房价代码不存在
	FileNotExisted            ErrCode = 1066 // 文件不存在
	RoomNoInvaild             ErrCode = 1067 // 房间未启用
	ClassExisted              ErrCode = 1068 // 班次已存在
	ClassNotExisted           ErrCode = 1069 // 班次不存在
	CheckTimeError            ErrCode = 1070 // 系统时间与营业时间不匹配
	CurrentClassIsShift       ErrCode = 1071 // 当前班次已交班
	PayPriceError             ErrCode = 1072 // 支付金额错误
	StockNotEnough            ErrCode = 1073 // 存量不足
	DBSaveError               ErrCode = 1074 // 数据存储错误
	DBAddError                ErrCode = 1075 // 数据添加错误
	DBUpdateError             ErrCode = 1076 // 数据更新错误
	DBDeleteError             ErrCode = 1077 // 数据删除错误
	TimeError                 ErrCode = 1078 // 时间错误
	OrderInfoError            ErrCode = 1079 // 预定信息错误
	NotVaildError             ErrCode = 1080 // 不允许
	Overdue                   ErrCode = 1081 // 已过期
	MaxOverError              ErrCode = 1082 // 超过最大值
	MinOverError              ErrCode = 1083 // 低于最小值
	ExistedError              ErrCode = 1084 // 已存在
	NotBindError              ErrCode = 1085 // 未绑定
	BindError                 ErrCode = 1086 // 绑定失败
	CalError                  ErrCode = 1087 // 计算错误
)

const (
	_ErrCode_name_0 = "默认返回值服务器维护中 请稍后再试帐号被禁用appid过期"
	_ErrCode_name_1 = "未知错误token失效"
	_ErrCode_name_2 = "成功"
	_ErrCode_name_3 = "禁止访问"
	_ErrCode_name_4 = "参数无效appid参数无效密文校验失败,aa用户名不存在或密码错误键值对重复未找到无效操作授权码错误已被使用已被激活激活码被禁止使用用户已存在验证码请求过于平凡邮箱发送失败手机发送失败手机号格式有问题验证码错误用户不存在topic已经存在topic不存在bundle_id不存在topic开启处理失败topic处理类型不存在topic不能为空设备不存在"
	_ErrCode_name_5 = "状态已存在上级菜单不存在菜单不存在用户权限不存在设备ID不存在商品处理类型不存在商品不存在商品正在打折商品可支付类型不存在商品已存在"
	_ErrCode_name_6 = "订单不存在商品未打折会话不匹配商品已恢复原价邀请用户不存在邀请用户级数已满用户未授权申请人已存在申请人不存在订单无效微信零钱重复提现提现金额错误微信提现失败重复提交bundle已存在权限已存在权限不存在房型不存在房型已存在房间不存在房间已存在房价代码或房价名称已存在房价代码不存在文件不存在房间未启用班次已存在班次不存在系统时间与营业时间不匹配当前班次已交班支付金额错误存量不足数据存储错误数据添加错误数据更新错误数据删除错误时间错误预定信息错误不允许已过期超过最大值低于最小值已存在未绑定绑定失败计算错误"
	_ErrCode_name_7 = "新消息心跳"
)

var (
	_ErrCode_index_0 = [...]uint8{0, 15, 49, 64, 75}
	_ErrCode_index_1 = [...]uint8{0, 12, 23}
	_ErrCode_index_4 = [...]uint16{0, 12, 29, 50, 83, 98, 107, 119, 134, 146, 158, 182, 197, 224, 242, 260, 284, 299, 314, 331, 345, 363, 386, 412, 429, 444}
	_ErrCode_index_5 = [...]uint8{0, 15, 36, 51, 72, 89, 116, 131, 149, 179, 194}
	_ErrCode_index_6 = [...]uint16{0, 15, 30, 45, 66, 87, 111, 126, 144, 162, 174, 198, 216, 234, 246, 261, 276, 291, 306, 321, 336, 351, 387, 408, 423, 438, 453, 468, 504, 525, 543, 555, 573, 591, 609, 627, 639, 657, 666, 675, 690, 705, 714, 723, 735, 747}
	_ErrCode_index_7 = [...]uint8{0, 9, 15}
)
