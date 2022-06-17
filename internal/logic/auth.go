package logic

import (
	"myGin/internal/dao/mysql"
	"myGin/internal/model"
	"myGin/pkg/snowflake"
)

// 创建角色
func (l *Logic) CreateRole(param *model.ParamCreateRole) error {
	roleId := snowflake.GenID()
	return mysql.CreateRole(roleId, param.RoleName, param.CreateBy)
}

// 获取所有的角色
func (l *Logic) AllRole() ([]*model.RoleModel, error) {
	return mysql.AllRole()
}
