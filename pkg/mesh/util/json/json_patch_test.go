package json

import (
	"fmt"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
)

func TestName(t *testing.T) {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	patchJSON := []byte(`[
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "remove", "path": "/height"}
	]`)

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original document: %s\n", original)
	fmt.Printf("Modified document: %s\n", modified)
}

func TestGet(t *testing.T) {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21, "test": {"value": "abcd"}}`)
	raw := json.Get(original, "test", "value")
	fmt.Printf("%+v", raw)
}

func TestGetNotExistsKey(t *testing.T) {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21, "test": {"value": "abcd"}}`)
	raw := Get(original, "test", "not_exists_key")
	fmt.Printf("%+v", raw.LastError())
	fmt.Printf("%+v", raw.ToString())
}

func TestGetNotJsonString(t *testing.T) {
	original := []byte(`404 not found`)
	raw := Get(original, "test")
	fmt.Printf("%+v", raw.LastError())
	fmt.Printf("%+v", raw.ToString())
}

func TestJsonString(t *testing.T) {
	name := "abc"
	patches := []Patch{
		{
			Op:    Replace,
			Path:  "name",
			Value: "hello_" + name,
		},
	}
	bs, err := Marshal(patches)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bs))
}
