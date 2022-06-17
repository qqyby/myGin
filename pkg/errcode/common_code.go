package errcode

var (
	Success             = NewError(0, "成功")
	OperateError        = NewError(10000000, "操作失败")
	UploadError         = NewError(10000001, "上传失败")
	GetError            = NewError(10000002, "上传失败")
	LicenseExpiredError = NewError(10000003, "License验证失败或已过期")

	ServerError               = NewError(20000000, "服务器内部错误")
	InvalidParams             = NewError(20000001, "入参错误")
	NotFound                  = NewError(20000002, "找不到")
	UnauthorizedAuthFailed    = NewError(20000003, "鉴权失败, 用户名或密码错误")
	UnauthorizedTokenError    = NewError(20000004, "鉴权失败, Token错误")
	UnauthorizedTokenTimeout  = NewError(20000005, "鉴权失败, Token超时")
	UnauthorizedTokenGenerate = NewError(20000006, "鉴权失败, Token生成失败")
	UnauthorizedToAccess      = NewError(20000007, "无权访问")
	TooManyRequests           = NewError(20000008, "请求过多")
	TaskStatusError           = NewError(30000001, "任务当前状态无法执行该操作")
)
