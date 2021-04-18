package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Resource interface {
	GetClient() *dynamodb.DynamoDB
	PutItem(tableName string, item map[string]*dynamodb.AttributeValue) error
	GetItem(tableName string, keyID string, keyValue string, target interface{}) error
}

type resource struct {
	db *dynamodb.DynamoDB
}

func (r resource) PutItem(tableName string, item map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}
	_, err := r.db.PutItem(input)
	return err
}

func (r resource) GetItem(tableName string, keyID string, keyValue string, target interface{}) error {
	key := map[string]*dynamodb.AttributeValue{
		keyID: {
			S: aws.String(keyValue),
		},
	}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}
	result, err := r.db.GetItem(input)
	if err != nil {
		return err
	}
	if result.Item == nil {
		return nil
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, target)
	if err != nil {
		return err
	}
	return nil
}

func (r resource) GetClient() *dynamodb.DynamoDB {
	return r.db
}

func New() Resource {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &resource{
		db: dynamodb.New(session),
	}
}
