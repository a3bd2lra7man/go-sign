//go:build integration
// +build integration

package dao

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/a3bd2lra7man/go-sign/pkg/sign/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection
var dao signDao

func init() {
	uri := "mongodb://127.0.0.1:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database("test-db")
	db.Collection("users").Drop(context.Background())

	coll = db.Collection("users")
	dao = New(db)
}

func TestCreateUser(t *testing.T) {
	//giver
	var identifier = strconv.Itoa(rand.Int())

	// when
	dao.CreateUser(identifier, []byte("pass"))

	// then
	res := coll.FindOne(context.Background(), bson.D{{Key: "identifier", Value: identifier}})

	var user entities.User = entities.User{}

	if res.Err() != nil {
		t.Fatalf("%v", res.Err().Error())
	}

	res.Decode(&user)

	if user.Allowed != true {
		t.Fatal("allowed is not inserted in the database")
	}

	if user.Password != "pass" {
		t.Fatal("password is not inserted in the database")
	}

	if user.Identifier != identifier {
		t.Fatal("password is not inserted in the database")
	}

	coll.Drop(context.Background())
}

func TestIsUserExist(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		isNotExist bool
		preRun     func()
	}{
		{
			name:       "case exist",
			identifier: "identifier@gmail.com",
			isNotExist: false,
			preRun: func() {
				var d = bson.D{{Key: "identifier", Value: "identifier@gmail.com"}, {Key: "password", Value: []byte("password")}, {Key: "allowed", Value: true}}
				coll.InsertOne(context.Background(), d)
			},
		},

		{
			name:       "case not exist",
			identifier: "identifier@gmail.com",
			isNotExist: true,
			preRun: func() {
				coll.Drop(context.Background())
			},
		},
	}

	for _, test := range tests {
		test.preRun()
		isNotExist := dao.IsNotExist(test.identifier)
		if isNotExist != test.isNotExist {
			t.Fatalf("%v : isNotExist expected to be : %v found %v", test.name, test.isNotExist, isNotExist)
		}

	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		expectErr  bool
		user       entities.User
		preRun     func()
	}{
		{
			name:       "case exist",
			identifier: "identifier@gmail.com",
			expectErr:  false,
			preRun: func() {
				var d = bson.D{{Key: "identifier", Value: "identifier@gmail.com"}, {Key: "password", Value: []byte("password")}, {Key: "allowed", Value: true}}
				coll.InsertOne(context.Background(), d)
			},
			user: entities.User{Identifier: "identifier@gmail.com", Password: "password", Allowed: true},
		},

		{
			name:       "case not exist",
			identifier: "identifier@gmail.com",
			expectErr:  true,
			preRun: func() {
				coll.Drop(context.Background())
			},
			user: entities.User{},
		},
	}

	for _, test := range tests {
		test.preRun()
		user, err := dao.GetUser(test.identifier)
		if (err == nil) == test.expectErr {
			t.Fatalf("%v : GetUser expected err : %v found : %v", test.name, test.expectErr, err)
		}
		userInDb := test.user
		if user.Identifier != userInDb.Identifier || user.Password != userInDb.Password || user.Allowed != userInDb.Allowed {
			t.Fatalf("%v : GetUser expected user : %v found %v", test.name, test.user, user)
		}

	}
}
