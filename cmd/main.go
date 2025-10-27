package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Metarock/personal-database/api"
	"github.com/Metarock/personal-database/vessel"
	"github.com/labstack/echo/v4"
)

/*
*
THIS IS MY PESONAL PROEJCT, LEARNING GO LANGUAGE AND THE CORE CONCEPTS OF DATABASE
Basically a copy of firebase and supabase
*/
func main() {
	// This is a placeholder for the main function.

	// support int, string, []byte, float, ...
	// temp data
	var a any
	a = false
	b := false

	fmt.Println(a == b)
	db, err := vessel.New()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(db)

	e := echo.New()

	e.HTTPErrorHandler = func(err error, context echo.Context) {
		context.JSON(http.StatusInternalServerError, vessel.Map{"error": err.Error()})
	}

	e.HideBanner = true

	e.POST("/api/:collname", server.HandlePostInsert)
	e.GET("/api/:collname", server.HandleGetQuery)
	log.Fatal(e.Start(":7777"))
}
