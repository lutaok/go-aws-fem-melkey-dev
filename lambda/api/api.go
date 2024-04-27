package api

import (
	"encoding/json"
	"fmt"
	"go-melkey-lambda/database"
	"go-melkey-lambda/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (h ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "request not valid",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "request not valid",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("request not valid, empty parameters")
	}

	// Check user existence
	userExists, err := h.dbStore.DoesUserExist(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error while checking user existence - %s", err.Error())
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, fmt.Errorf("error while trying to insert an existing user")
	}

	user, err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error while trying to hash password")
	}

	err = h.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error while registering the user - %s", err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       "successfully registered user",
		StatusCode: http.StatusOK,
	}, nil
}

func (h ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginUser types.LoginUser
	err := json.Unmarshal([]byte(request.Body), &loginUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Request not valid",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := h.dbStore.GetUser(loginUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error while trying to get the user - %w", err)
	}

	if err := types.ValidatePassword(user.PasswordHash, loginUser.Password); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("invalid credentials provided - %w", err)
	}

	accessToken, err := types.CreateToken(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("could not create token - %w", err)
	}

	successMsg := fmt.Sprintf(`{ "access_token": "%s" }`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}
