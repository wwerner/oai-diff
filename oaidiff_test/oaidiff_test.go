package oaidiff_test

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"github.com/wwerner/oaidiff/oaidiff"
	"testing"
)

func TestAddModelProp(t *testing.T) {
	swaggerLhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic.yml")
	swaggerRhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic_add-model-prop.yml")

	require.NoError(t, err)

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(
			openapi3.Schema{},
		),
		cmp.Transformer("RawJSONFilter", func(rm json.RawMessage) string {
			return ""
		}),

	}
	if diff := cmp.Diff(swaggerLhs, swaggerRhs, opts...); diff != "" {
		t.Errorf("%s", diff)
	}
}

func TestCustomReporter(t *testing.T) {
	swaggerLhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic.yml")
	swaggerRhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic_add-model-prop.yml")

	require.NoError(t, err)

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(
			openapi3.Schema{},
		),
		cmp.Transformer("RawJSONFilter", func(rm json.RawMessage) string {
			return ""
		}),

	}
	if diff := oaidiff.Diff(swaggerLhs, swaggerRhs, opts...); diff != "" {
		t.Errorf("%s", diff)
	}
}
