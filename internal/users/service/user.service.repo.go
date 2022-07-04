package service

import (
	"context"
	"errors"
	"log"
	"time"

	database "github.com/yasensim/toyshop/internal/db"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yasensim/toyshop/internal/users"

	"golang.org/x/crypto/bcrypt"
)

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func GetUsersDataStore() users.UserDatastore {
	return &TableBasics{database.CreateLocalClient(), "users"}
}

func (basics TableBasics) CreateUser(user *users.User) error {
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return errors.New("user service repo - cannot have empty fields")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return errors.New("user service repo - password encryption failed")
	}
	user.Password = string(pass)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = basics.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}

	return err
}

func (basics TableBasics) FindUser(email, password string) (*users.User, error) {
	user := &users.User{}

	if email == "" || password == "" {
		return nil, errors.New("user service repo - cannot have empty email or password")
	}
	eml, err := attributevalue.Marshal(email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := basics.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{Key: map[string]types.AttributeValue{"email": eml}, TableName: aws.String(basics.TableName)})

	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", email, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, user)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil {
		return nil, errors.New("user service repo - invalid login credentials; please try again")
	}

	return user, nil
}
