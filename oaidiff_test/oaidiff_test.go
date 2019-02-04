package oaidiff_test

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddModelProp(t *testing.T) {
	swaggerLhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic.yml")
	swaggerRhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic_add-model.yml")

	require.NoError(t, err)
	require.Equal(t, swaggerLhs, swaggerRhs)
}
