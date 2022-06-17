package model

// 对应角色表数据库的字段
type RoleModel struct {
	RoleID     string `json:"role_id" db:"role_id"`
	RoleName   string `json:"role_name" db:"role_name"`
	IsAdmin    bool   `json:"is_admin" db:"is_admin"`
	IsDel      bool   `json:"-" db:"is_del"`
	Remark     string `json:"remark" db:"remark"`
	CreateBy   string `json:"create_by" db:"create_by"`
	CreateTime string `json:"create_time" db:"create_time"`
	DeleteTime string `json:"-" db:"delete_time"`
}

// 创建角色的请求
type ParamCreateRole struct {
	RoleName string `json:"role_name" validate:"required,min=2,max=32"`
	CreateBy string
}

// ---------------------用户表-----------------------------
// 对应用户表数据库的字段
type UserModel struct {
	UserID     string `json:"user_id" db:"user_id"`
	UserName   string `json:"user_name" db:"user_name"`
	UserPwd    string `json:"-" db:"user_pwd"`
	RoleID     string `json:"role_id" db:"role_id"`
	RoleName   string `json:"role_name" db:"role_name"`
	Remark     string `json:"remark" db:"remark"`
	IsDel      bool   `json:"-" db:"is_del"`
	CreateTime string `json:"create_time" db:"create_time"`
	DeleteTime string `json:"-" db:"delete_time"`
}

// 登录请求
type ParamLogin struct {
	UserName string `json:"user_name" validate:"required,min=2,max=32"`
	UserPwd  string `json:"user_pwd" validate:"required"`
}

// 更新自己密码的请求
type ParamUpdateMyPwd struct {
	UserPwd string `json:"user_pwd" validate:"required"`
	UserID  string
}

// ----------------------------权限表---------------------------------------
// 对应权限表数据库的字段
type PrivilegeModel struct {
	MenuID     string `json:"menu_id" db:"menu_id"`
	MenuName   string `json:"menu_name" db:"menu_name"`
	MenuPath   string `json:"menu_path" db:"menu_path"`
	IsDel      bool   `json:"-" db:"is_del"`
	CreateTime string `json:"create_time" db:"create_time"`
	DeleteTime string `json:"-" db:"delete_time"`
}

// 获取单角色权限的响应
type PrivilegeRow struct {
	RoleID   string `json:"role_id,omitempty" db:"role_id"`
	RoleName string `json:"role_name,omitempty" db:"role_name"`
	MenuID   string `json:"menu_id,omitempty" db:"menu_id"`
	MenuName string `json:"menu_name,omitempty" db:"menu_name"`
	MenuPath string `json:"menu_path,omitempty" db:"menu_path"`
	IsAdmin  bool   `json:"is_admin,omitempty" db:"is_admin"`
}
