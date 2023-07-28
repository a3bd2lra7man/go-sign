package dao

import (
	"context"
	"errors"
	"time"

	errs "github.com/a3bd2lra7man/go-sign/pkg/err"
	"github.com/a3bd2lra7man/go-sign/pkg/otp/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type otpDao struct {
	mongo.Collection
}

func New(db *mongo.Database) otpDao {
	col := *db.Collection("otp_code")

	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(time.Now().Add(time.Hour).Unix())), // Will be removed after 24 Hours.
	})
	return otpDao{col}
}

func (d otpDao) Save(identifier, code, role string) error {
	var doc = bson.D{{Key: "$set", Value: bson.D{{Key: "identifier", Value: identifier}, {Key: "role", Value: role}, {Key: "code", Value: code}, {Key: "expireTime", Value: time.Now().Add(time.Minute * 5)}}}}
	_, err := d.UpdateOne(context.Background(), bson.D{{Key: "identifier", Value: identifier}, {Key: "role", Value: role}}, doc, options.Update().SetUpsert(true))

	if err != nil {
		return errs.UnExpectedError{}
	}

	return nil
}

type OtpError uint8

const (
	UnFound OtpError = 0
)

func (err OtpError) Error() string {
	switch err {
	case UnFound:
		return "UnFound"

	}
	return errors.New("error").Error()
}

func (d otpDao) GetCode(identifier, role string) (entities.OtpCode, error) {
	var otpCode entities.OtpCode
	err := d.FindOne(context.Background(), bson.D{{Key: "identifier", Value: identifier}, {Key: "role", Value: role}}).Decode(&otpCode)

	if err == mongo.ErrNoDocuments {
		return entities.OtpCode{}, UnFound
	} else if err != nil {
		return entities.OtpCode{}, errs.UnExpectedError{}
	}

	return otpCode, nil

}

var _ IOtpDao = otpDao{}
