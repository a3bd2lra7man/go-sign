package roles

import (
	"github.com/a3bd2lra7man/sign/pkg/roles/internal/core"
	"github.com/a3bd2lra7man/sign/pkg/roles/internal/dao"
	"github.com/a3bd2lra7man/sign/pkg/roles/internal/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

var manager core.RolesManager

func SetUp(db *mongo.Database) {
	manager = core.New(dao.New(db))
}

func Pend(role string, identifier string) error {
	return manager.ChangeStatus(identifier, role, entities.Pending)
}

func Activate(role string, identifier string) error {
	return manager.ChangeStatus(identifier, role, entities.Active)
}

func Ban(role string, identifier string) error {
	return manager.ChangeStatus(identifier, role, entities.Banned)
}

func Delete(role string, identifier string) error {
	return manager.Delete(identifier, role)
}

func CheckStatus(role string, identifier string) error {
	status := manager.GetStatus(identifier, role)

	if status != entities.Active {
		return core.UnActive
	}

	return nil
}

func GetAllHasRole(role string) {

}
