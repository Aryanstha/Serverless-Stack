package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/byron/serverless/pkg/users"
)

var ErrorMethodNotAllowed = "method not allowed "

type ErrorBody struct {
	ErrorMessage *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error) {
   email := req.QueryStringParameters["email"]
   if len(email)>0{
	  result , err :=  users.FetchUser(email,tablename,dynaClient)

	  if err != nil {
		  return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	  }
	  return apiResponse(http.StatusOK, result)
   }
   result , err := users.FetchUsers(tablename,dynaClient)

   if err != nil {
	return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
   }
   return apiResponse(http.StatusOK,result)

}

func CreateUser(req events.APIGatewayProxyRequest, tablename string , dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse,error) {
  result, err := users.CreateUser(req,tablename,dynaClient)
  if err != nil {
	  return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error()),})
  }
  return apiResponse(http.StatusCreated, result)
}

func UpdateUser(req events.APIGatewayProxyRequest,tablename string , dynaClient dynamodbiface.DynamoDBAPI)(
     *events.APIGatewayProxyResponse, error) {
 result , err := users.UpdateUser(req ,tablename, dynaClient)

 if err != nil { 
	 return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
 }
 return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest,tablename string , dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error) {
err := users.DeleteUser(req, tablename, dynaClient)
if err != nil{
	return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
}
return apiResponse(http.StatusOK, nil )
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
