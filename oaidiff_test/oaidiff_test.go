package oaidiff_test

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"github.com/wwerner/oaidiff/oaidiff"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
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

func TestCustomReporterWithNewVersion(t *testing.T) {
	swaggerLhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic.yml")
	swaggerRhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic_new-version.yml")

	require.NoError(t, err)
	testCustomOpenAPI3Reporter(t, swaggerLhs, swaggerRhs)
}

func TestCustomReporterWithAddedModelProp(t *testing.T) {
	swaggerLhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic.yml")
	swaggerRhs, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("resources/oai-basic_add-model-prop.yml")

	require.NoError(t, err)
	testCustomOpenAPI3Reporter(t, swaggerLhs, swaggerRhs)
}

func TestCustomReporterWithGoSwaggerAndAddedModelProp(t *testing.T) {
	swaggerLhs, err := loads.Spec("resources/oai-basic.yml")
	swaggerRhs, err := loads.Spec("resources/oai-basic_add-model-prop.yml")

	require.NoError(t, err)
	testCustomGoSwaggerReporter(t, swaggerLhs, swaggerRhs)
}

func TestCustomReporterWithNakedYmlAndAddedModelProp(t *testing.T) {
	swaggerLhs := map[string]interface{}{}
	swaggerRhs := map[string]interface{}{}

	yamlFile, err := ioutil.ReadFile("resources/oai-basic.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &swaggerLhs)
	if err != nil {
		panic(err)
	}

	yamlFile, err = ioutil.ReadFile("resources/oai-basic_add-operation.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &swaggerRhs)
	if err != nil {
		panic(err)
	}

	require.NoError(t, err)
	testCustomNakedYamlReporter(t, swaggerLhs, swaggerRhs)
}

func testCustomOpenAPI3Reporter(t *testing.T, old, new *openapi3.Swagger) {
	opts := []cmp.Option{
		cmpopts.IgnoreTypes(
			openapi3.ExtensionProps{},
			openapi3.ExtensionProps{}.Extensions,

		),
		cmpopts.IgnoreUnexported(
			openapi3.Schema{},

		),
		cmp.AllowUnexported(
			openapi3.Schema{},
		),
		cmp.Transformer("RawJSONFilter", func(rm json.RawMessage) string {
			var r, _ = json.Marshal(rm)
			return string(r)
		}),
		cmp.Transformer("UnwrapValue", func(v reflect.Value) string {
			return fmt.Sprintf("%#v", v)
		}),
	}
	if diff, changes := oaidiff.Diff(old, new, opts...); diff != "" {
		spew.Dump(changes)
		t.Errorf("%s", diff)
	}
}

func testCustomNakedYamlReporter(t *testing.T, old, new map[string]interface{}) {
	opts := []cmp.Option{
		cmpopts.IgnoreTypes(
		),
		cmpopts.IgnoreUnexported(
		),
		//cmpopts.IgnoreFields(analysis.Spec{}, "spec" ),
		cmp.AllowUnexported(
		),
	}
	o := map[string]string{}
	for k, v := range old {
		flatten2(k, v, o)
	}
	n := map[string]string{}
	for k, v := range new {
		flatten2(k, v, n)
	}

	if diff, changes := oaidiff.Diff(o, n, opts...); diff != "" {
		spew.Dump(changes)
		t.Errorf("%s", diff)
	}
}

func testCustomGoSwaggerReporter(t *testing.T, old, new *loads.Document) {
	opts := []cmp.Option{
		cmpopts.IgnoreTypes(
			//spec.Ref{}.Ref,
			//spec.Refable{},
			//loads.Document{}.Analyzer,
			//analysis.Spec{},
			// spec.Refable{},
			// spec.Ref{},
			json.RawMessage{},
		),
		cmpopts.IgnoreUnexported(
			spec.Ref{}.Ref,
		),
		//cmpopts.IgnoreFields(analysis.Spec{}, "spec" ),
		cmp.AllowUnexported(
			loads.Document{},
			analysis.Spec{},
			//loads.Document{},
		),
		cmp.Transformer("UnwrapValue", func(r *analysis.Spec) interface{} {
			return r.AllDefinitions()
		}),
	}
	var o = flatten(old)
	var n = flatten(new)

	if diff, changes := oaidiff.Diff(o, n, opts...); diff != "" {
		spew.Dump(changes)
		t.Errorf("%s", diff)
	}
}

func flatten(specDoc *loads.Document) *loads.Document {
	flattenOpts := &analysis.FlattenOpts{
		// defaults
		Minimal:      false,
		Verbose:      true,
		Expand:       true,
		RemoveUnused: false,
	}
	flattenOpts.BasePath = specDoc.SpecFilePath()
	flattenOpts.Spec = analysis.New(specDoc.Spec())
	if err := analysis.Flatten(*flattenOpts); err != nil {
		panic(err)
	}
	return specDoc
}

func flatten2(prefix string, value interface{}, flatmap map[string]string) {
	submap, ok := value.(map[interface{}]interface{})
	if ok {
		for k, v := range submap {
			flatten2(prefix+"."+k.(string), v, flatmap)
		}
		return
	}
	stringlist, ok := value.([]interface{})
	if ok {
		flatten2(fmt.Sprintf("%s.size", prefix), len(stringlist), flatmap)
		for i, v := range stringlist {
			flatten2(fmt.Sprintf("%s.%d", prefix, i), v, flatmap)
		}
		return
	}
	flatmap[prefix] = fmt.Sprintf("%v", value)
}
