package dao

import "github.com/a3bd2lra7man/sign/pkg/roles/internal/entities"

type IRoleDao interface {
	ChangeStatus(id string, role string, status entities.RoleStatus) error
	Delete(id string, role string) error
	GetStatus(id string, role string) entities.RoleStatus
}
