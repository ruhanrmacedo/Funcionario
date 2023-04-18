package main

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type Funcionario struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CPF  int    `json:"cpf"`
}

func InsertProduct(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var funcionario Funcionario
	err := json.Unmarshal([]byte(request.Body), &funcionario)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	funcionario.ID = uuid.New().String()

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Funcionarios"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(funcionario.ID),
			},
			"name": {
				S: aws.String(funcionario.Name),
			},
			"cpf": {
				N: aws.String(strconv.Itoa(funcionario.CPF)),
			},
		},
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(funcionario)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil

}

func main() {
	lambda.Start(InsertProduct)
}
