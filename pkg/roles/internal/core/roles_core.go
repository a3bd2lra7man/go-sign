package core

import (
	"github.com/a3bd2lra7man/go-sign/pkg/roles/internal/dao"
	"github.com/a3bd2lra7man/go-sign/pkg/roles/internal/entities"
)

type RolesManager struct {
	dao dao.IRoleDao
}

func New(dao dao.IRoleDao) RolesManager {
	return RolesManager{dao: dao}
}

func (m RolesManager) ChangeStatus(identifier string, role string, status entities.RoleStatus) error {
	return m.dao.ChangeStatus(identifier, role, status)
}

func (m RolesManager) Delete(identifier string, role string) error {
	return m.dao.Delete(identifier, role)
}

func (m RolesManager) GetStatus(identifier string, role string) entities.RoleStatus {
	return m.dao.GetStatus(identifier, role)
}
