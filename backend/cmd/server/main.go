package main

import (
	"log"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/maruks2003/sudoku/backend/internal/api"
	"github.com/maruks2003/sudoku/backend/internal/docs"
	"github.com/maruks2003/sudoku/backend/internal/handler"
	ginmw "github.com/oapi-codegen/gin-middleware"
)

func main() {
	spec, err := api.GetSpec()
	if err != nil {
		log.Fatal(err)
	}

	spec.Servers = nil // don't validate hosts

	router := gin.Default()

	if err := docs.Register(router, spec); err != nil {
		log.Fatal(err)
	}

	validated := router.Group("")

	validated.Use(ginmw.OapiRequestValidatorWithOptions(spec, &ginmw.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, // real auth later
		},
	}))

	strictHandler := api.NewStrictHandler(handler.New(), nil)
	api.RegisterHandlers(validated, strictHandler)

	log.Fatal(router.Run(":8080"))
}
