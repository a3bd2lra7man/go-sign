package dao

import "github.com/a3bd2lra7man/go-sign/pkg/roles/internal/entities"

type IRoleDao interface {
	ChangeStatus(id string, role string, status entities.RoleStatus) error
	Delete(id string, role string) error
	GetStatus(id string, roles ...string) entities.RoleStatus
}
