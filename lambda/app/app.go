package app

import (
	"go-melkey-lambda/api"
	"go-melkey-lambda/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	// initialize dbstore
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)

	return App{
		ApiHandler: apiHandler,
	}
}
