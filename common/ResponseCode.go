package common

var (
	OK  = response(200, "ok")    // 通用成功
	Err = response(500, "服务器错误") // 通用错误

	// ErrParam 1-服务级 01-模块级 01-具体错误码 10101
	// ErrParam 服务级错误码
	ErrParam       = response(10001, "参数有误")
	ErrSignParam   = response(10002, "签名参数有误")
	ErrIdBlank     = response(10003, "Id不能为空")
	ErrDataNoExist = response(10004, "数据不存在")

	// ErrUserService 模块级错误码 - 用户模块
	ErrUserService = response(20100, "用户服务异常")
	ErrUserPhone   = response(20101, "用户手机号不合法")
	ErrUserCaptcha = response(20102, "用户验证码有误")

	// ErrOrderService 库存模块
	ErrOrderService = response(20200, "订单服务异常")
	ErrOrderOutTime = response(20201, "订单超时")

	// ErrAdminService 模块级错误码 - 管理员模块
	ErrAdminService = response(20300, "管理员服务异常")
	ErrAdminPhone   = response(20301, "用户手机号不合法")
	ErrAdminCaptcha = response(20302, "用户验证码有误")
)
