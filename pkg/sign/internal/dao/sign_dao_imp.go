package dao

import (
	"context"

	errs "github.com/a3bd2lra7man/go-sign/pkg/err"
	"github.com/a3bd2lra7man/go-sign/pkg/sign/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signDao struct {
	mongo.Collection
}

func New(db *mongo.Database) signDao {
	return signDao{*db.Collection("users")}
}

func (s signDao) CreateUser(identifier string, password []byte) error {
	var d = bson.D{{Key: "identifier", Value: identifier}, {Key: "password", Value: password}, {Key: "allowed", Value: true}}
	_, err := s.InsertOne(context.Background(), d)
	if err != nil {
		return errs.UnExpectedError{}
	}

	return nil
}

func (s signDao) GetUser(identifier string) (entities.User, error) {
	var user entities.User
	err := s.FindOne(context.Background(), bson.D{{Key: "identifier", Value: identifier}}).Decode(&user)

	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (s signDao) IsNotExist(identifier string) bool {
	opts := options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}})

	res := s.FindOne(context.Background(), bson.D{{Key: "identifier", Value: identifier}}, opts)

	return res.Err() == mongo.ErrNoDocuments
}

var _ ISingDao = signDao{}
