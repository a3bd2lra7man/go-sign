//go:build integration
// +build integration

package dao

import (
	"context"
	"testing"

	"github.com/a3bd2lra7man/go-sign/pkg/otp/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection
var dao otpDao

func init() {
	uri := "mongodb://127.0.0.1:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database("test-db")
	db.Collection("otp_code").Drop(context.Background())

	coll = db.Collection("otp_code")
	dao = New(db)
}

func TestSave(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		code       string
		role       string
	}{
		{
			name:       "success",
			identifier: "identifier@gmail.com",
			code:       "random_code",
			role:       "admin",
		},
	}

	for _, test := range tests {
		err := dao.Save(test.identifier, test.code, test.role)
		if err != nil {
			t.Fatalf("%v : err expected to be : %v found %v", test.name, nil, err)
		}
		var otpCode entities.OtpCode
		coll.FindOne(context.Background(), bson.D{{Key: "identifier", Value: test.identifier}, {Key: "role", Value: test.role}}).Decode(&otpCode)

		if otpCode.Identifier != test.identifier {
			t.Fatalf("%v : err expected identifier to be : %v found %v", test.name, otpCode.Identifier, test.identifier)
		}
		if otpCode.Code != test.code {
			t.Fatalf("%v : err expected code to be : %v found %v", test.name, otpCode.Code, test.code)
		}
		if otpCode.Role != test.role {
			t.Fatalf("%v : err expected role to be : %v found %v", test.name, otpCode.Role, test.role)
		}

	}
}

func TestGetCode(t *testing.T) {
	var d = bson.D{{Key: "$set", Value: bson.D{{Key: "identifier", Value: "code@gmail.com"}, {Key: "role", Value: "admin"}, {Key: "code", Value: "value"}}}}
	opts := options.Update().SetUpsert(true)
	coll.UpdateOne(context.Background(), bson.D{{Key: "identifier", Value: "code@gmail.com"}, {Key: "role", Value: "admin"}}, d, opts)

	tests := []struct {
		name          string
		identifier    string
		code          string
		role          string
		expectToFound bool
	}{
		{
			name:          "different email",
			identifier:    "different@gmail.com",
			role:          "admin",
			expectToFound: false,
		},
		{
			name:          "different role",
			identifier:    "code@gmail.com",
			role:          "role",
			expectToFound: false,
		},
		{
			name:          "success",
			identifier:    "code@gmail.com",
			role:          "admin",
			code:          "value",
			expectToFound: true,
		},
	}

	for _, test := range tests {
		code, err := dao.GetCode(test.identifier, test.role)
		if !test.expectToFound {
			if err != UnFound {
				t.Fatalf("%v : expected code to be : %v found %v", test.name, test.code, code)
			}
		} else {

			var otpCode entities.OtpCode
			coll.FindOne(context.Background(), bson.D{{Key: "identifier", Value: test.identifier}, {Key: "role", Value: test.role}}).Decode(&otpCode)

			if otpCode.Identifier != test.identifier {
				t.Fatalf("%v : err expected identifier to be : %v found : %v", test.name, otpCode.Identifier, test.identifier)
			}
			if otpCode.Code != test.code {
				t.Fatalf("%v : err expected code to be : %v found : %v", test.name, otpCode.Code, test.code)
			}
			if otpCode.Role != test.role {
				t.Fatalf("%v : err expected role to be : %v found : %v", test.name, otpCode.Role, test.role)
			}
		}

	}
}
