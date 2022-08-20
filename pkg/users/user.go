package users

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/byron/serverless/pkg/validators"
)
var (
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnMashalRecord = "failed to unmarsha record "
	ErrorInvalidUserData = "invalid user data "
	ErrorInvalidEmail = "Invalid email "
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotMarshalItem = "could not marshal item "
	ErrorCouldNOtDynamoPutItem = "could not dynamo put item"
	ErrorUserDoesNotExist = "user doesnt not exist"
	ErrorUserAlreadyExist = "user alredy exists"
)


type User struct{
	Email     string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string  `json:"lastName"`
}


func FetchUser(email string , tablename string , dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{S: aws.String(email)},
		},
		TableName: aws.String(tablename),
	}
	result , err := dynaClient.GetItem(input)

	if err != nil {
		return nil , errors.New(ErrorFailedToFetchRecord)
	}
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil , errors.New(ErrorFailedToUnMashalRecord)
	}
	return item , nil 
}

func FetchUsers(tablename string , dynaClient dynamodbiface.DynamoDBAPI)(*[]User,error)  {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}
	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil , errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items,item)

	return item , nil 

}
func CreateUser(req events.APIGatewayProxyRequest, tablename string , dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
 if err := 	json.Unmarshal([]byte(req.Body),&u); err != nil {
	 return nil , errors.New(ErrorInvalidUserData)
 }
 if !validators.IsEmailValid(u.Email){
	 return nil , errors.New(ErrorInvalidEmail)
 }
 currentUser,_ := FetchUser(u.Email,tablename,dynaClient)
 if currentUser != nil && len(currentUser.Email) != 0 {
	 return nil , errors.New(ErrorUserAlreadyExist)
 }
 av, err := dynamodbattribute.MarshalMap(u)
 if err != nil {
	 return nil , errors.New(ErrorCouldNotMarshalItem)
 }
 input:= &dynamodb.PutItemInput{
	 Item: av,
	 TableName: aws.String(tablename),
 }
 _, err = dynaClient.PutItem(input)
if err != nil {
	return nil , errors.New(ErrorCouldNOtDynamoPutItem)
}
return &u , nil 

 }






func UpdateUser(req events.APIGatewayProxyRequest, tablename string , dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User

    if err := json.Unmarshal([]byte(req.Body),&u);err != nil {
		return nil , errors.New(ErrorInvalidEmail)
	}
	currentUser, _:= FetchUser(u.Email,tablename,dynaClient)
	if currentUser !=nil && len(currentUser.Email)==0 {
		return nil , errors.New(ErrorUserDoesNotExist)
	}
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil , errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tablename),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil , errors.New(ErrorCouldNOtDynamoPutItem)
	}
	return &u, nil
}











func DeleteUser(req events.APIGatewayProxyRequest,tablename string , dynaClient dynamodbiface.DynamoDBAPI) error {
	email:= req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return  errors.New(ErrorCouldNotDeleteItem)
	}
	return nil 
}