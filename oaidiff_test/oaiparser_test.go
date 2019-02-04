package oaidiff_test

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseYmlSpec(t *testing.T) {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic.yml")

	require.NoError(t, err)
	require.Equal(t, "Swagger Petstore", swagger.Info.Title)
}
