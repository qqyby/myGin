package mysql

import (
	"myGin/internal/model"
	"myGin/pkg/utils"
)

// ---------------------------权限--------------------------------
// 获取普通角色的权限列表
func GetPrivilegesByRoleID(roleID string) ([]*model.PrivilegeRow, error) {
	var privileges []*model.PrivilegeRow
	query := "SELECT r.role_id AS role_id, r.is_admin AS is_admin, r.role_name AS role_name, m.menu_id AS menu_id, " +
		"m.menu_name AS menu_name, m.menu_path AS menu_path FROM bzr_privilege AS p " +
		"LEFT JOIN bzr_role AS r ON p.role_id=r.role_id " +
		"LEFT JOIN bzr_menu_list AS m ON p.menu_id=m.menu_id " +
		"WHERE p.role_id = ? AND r.is_del=0 ORDER BY id"
	err := db.Select(&privileges, query, roleID)
	return privileges, err
}

// 获取超级管理员的权限列表
func AllMenuList() ([]*model.PrivilegeModel, error) {
	var privileges []*model.PrivilegeModel
	query := "SELECT menu_id, menu_name, menu_path FROM bzr_menu_list WHERE is_del=0 ORDER BY id"
	err := db.Select(&privileges, query)
	return privileges, err
}

// ---------------------------角色-----------------------------------
// 通过角色ID获取角色的信息
func GetRoleByRoleID(roleID string) (*model.RoleModel, error) {
	var r model.RoleModel
	query := "SELECT `role_id`, `role_name`, `is_admin`, `remark`, `create_time` " +
		"FROM bzr_role WHERE role_id=? AND is_del=0"
	err := db.Get(&r, query, roleID)
	return &r, err
}

// 创建角色
func CreateRole(roleID, roleName, createBy string) error {
	query := "INSERT INTO bzr_role (`role_id`, `role_name`, `create_time`, `create_by`) VALUES(?, ?, ?, ?)"
	_, err := db.Exec(query, roleID, roleName, utils.NowDateTimeStr(), createBy)
	return err
}

// 获取所有的角色
func AllRole() ([]*model.RoleModel, error) {
	var roles = make([]*model.RoleModel, 0)
	query := "SELECT `role_id`, `role_name`, `is_admin`, `create_time`, `create_by` FROM bzr_role " +
		"WHERE is_del=0 ORDER BY id"
	err := db.Select(&roles, query)
	return roles, err
}

// ---------------------------用户-----------------------------------
// 通过角色名称获取角色信息
func GetUserByName(userName string) (*model.UserModel, error) {
	var u model.UserModel
	query := "SELECT `user_id`, `user_name`, `user_pwd`, `role_id`, `role_name`, `remark`, `create_time` " +
		"FROM bzr_user WHERE user_name=? AND is_del=0"
	err := db.Get(&u, query, userName)
	return &u, err
}

// 更新角色的密码
func UpdateUserPwd(userID, userPwd string) error {
	query := "UPDATE bzr_user SET `user_pwd`=? WHERE user_id=?"
	_, err := db.Exec(query, userPwd, userID)
	return err
}
