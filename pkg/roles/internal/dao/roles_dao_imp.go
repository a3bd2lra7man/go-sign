package dao

import (
	"context"

	errs "github.com/a3bd2lra7man/sign/pkg/err"
	"github.com/a3bd2lra7man/sign/pkg/roles/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rolesDao struct {
	mongo.Collection
}

func New(db *mongo.Database) rolesDao {
	return rolesDao{*db.Collection("roles_status")}
}

func (d rolesDao) ChangeStatus(id string, role string, status entities.RoleStatus) error {
	var doc = bson.D{{Key: "$set", Value: bson.D{{Key: "identifier", Value: id}, {Key: "status", Value: status}, {Key: "role", Value: role}}}}
	_, err := d.UpdateOne(context.Background(), bson.D{{Key: "identifier", Value: id}, {Key: "role", Value: role}}, doc, options.Update().SetUpsert(true))

	if err != nil {
		return errs.UnExpectedError{}
	}

	return nil
}

func (d rolesDao) Delete(id string, role string) error {
	_, err := d.DeleteOne(context.Background(), bson.D{{Key: "identifier", Value: id}, {Key: "role", Value: role}})

	if err != nil {
		return errs.UnExpectedError{}
	}

	return nil
}

func (d rolesDao) GetStatus(id string, role string) entities.RoleStatus {

	var roleDoc entities.UserRole
	err := d.FindOne(context.Background(), bson.D{{Key: "identifier", Value: id}, {Key: "role", Value: role}}).Decode(&roleDoc)

	if err != nil {
		return entities.UnFound
	}
	return roleDoc.Status
}

var _ IRoleDao = rolesDao{}
