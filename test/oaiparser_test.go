package test

import (
	"github.com/getkin/kin-openapi/openapi3"
	"testing"
)

func testParseYmlSpec(t *testing.T) {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("swagger.json")

	if err != nil {
		t.Error(err)
	}

	if(swagger.Info.Title != "Swagger Petstore") {
		t.Fail()
	}
}
