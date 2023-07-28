//go:build integration
// +build integration

package dao

import (
	"context"
	"testing"

	"github.com/a3bd2lra7man/sign/pkg/roles/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection
var dao rolesDao

func init() {
	uri := "mongodb://127.0.0.1:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database("test-db")
	db.Collection("roles_status").Drop(context.Background())

	coll = db.Collection("roles_status")
	dao = New(db)
}

func TestChangeStatus(t *testing.T) {
	var d = bson.D{{Key: "$set", Value: bson.D{{Key: "identifier", Value: "identifier@gmail.com"}, {Key: "status", Value: entities.Pending}, {Key: "role", Value: "admin"}}}}
	opts := options.Update().SetUpsert(true)
	coll.UpdateOne(context.Background(), bson.D{{Key: "identifier", Value: "identifier@gmail.com"}, {Key: "role", Value: "admin"}}, d, opts)

	tests := []struct {
		name           string
		identifier     string
		expectedStatus entities.RoleStatus
	}{
		{
			name:           "change status to active",
			identifier:     "identifier@gmail.com",
			expectedStatus: entities.Active,
		},
		{
			name:           "change status to banned",
			identifier:     "identifier@gmail.com",
			expectedStatus: entities.Banned,
		},
		{
			name:           "change status to active",
			identifier:     "identifier@gmail.com",
			expectedStatus: entities.Pending,
		},
	}

	for _, test := range tests {
		err := dao.ChangeStatus(test.identifier, "admin", test.expectedStatus)
		if err != nil {
			t.Fatalf("%v : err expected to be : %v found %v", test.name, nil, err)
		}
		var role entities.UserRole
		coll.FindOne(context.Background(), bson.D{{Key: "identifier", Value: test.identifier}, {Key: "role", Value: "admin"}}).Decode(&role)

		if role.Status != test.expectedStatus {
			t.Fatalf("%v : err expected status to be : %v found %v", test.name, role.Status, test.expectedStatus)
		}

	}
}

func TestDelete(t *testing.T) { // s
	var d = bson.D{{Key: "$set", Value: bson.D{{Key: "identifier", Value: "delete@gmail.com"}, {Key: "status", Value: entities.Pending}, {Key: "role", Value: "admin"}}}}
	opts := options.Update().SetUpsert(true)
	coll.UpdateOne(context.Background(), bson.D{{Key: "identifier", Value: "delete@gmail.com"}, {Key: "role", Value: "admin"}}, d, opts)

	tests := []struct {
		name        string
		identifier  string
		expectError bool
	}{
		{
			name:        "success",
			identifier:  "delete@gmail.com",
			expectError: false,
		},
	}

	for _, test := range tests {
		err := dao.Delete(test.identifier, "admin")
		if (err == nil) == test.expectError {
			t.Fatalf("%v : err expected to be : %v found %v", test.name, nil, err)
		}

		res := coll.FindOne(context.Background(), bson.D{{Key: "identifier", Value: test.identifier}, {Key: "role", Value: "admin"}})

		if res.Err() != mongo.ErrNoDocuments {
			t.Fatalf("err expected document not to be found")
		}
	}
}

func TestGetStatus(t *testing.T) {
	var d = bson.D{{Key: "$set", Value: bson.D{{Key: "identifier", Value: "status@gmail.com"}, {Key: "status", Value: entities.Active}, {Key: "role", Value: "admin"}}}}
	opts := options.Update().SetUpsert(true)
	coll.UpdateOne(context.Background(), bson.D{{Key: "identifier", Value: "status@gmail.com"}, {Key: "role", Value: "admin"}}, d, opts)

	tests := []struct {
		name           string
		identifier     string
		expectedStatus entities.RoleStatus
	}{
		{
			name:           "success",
			identifier:     "status@gmail.com",
			expectedStatus: entities.Active,
		},
	}

	for _, test := range tests {
		status := dao.GetStatus(test.identifier, "admin")
		if status != test.expectedStatus {
			t.Fatalf("%v : err expected status to be : %v found %v", test.name, test.expectedStatus, status)
		}
	}
}
