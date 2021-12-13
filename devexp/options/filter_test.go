package options

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/drahoslavzan/srvutils/devexp/options/serialize"
)

const jsonBinaryFilter = `["field", "=", 10]`
const jsonUnaryFilter = `["!", ` + jsonBinaryFilter + `]`
const jsonLoadOptions1 = `{"filter":` + jsonBinaryFilter + `}`
const jsonLoadOptions2 = `{"filter":` + jsonUnaryFilter + `}`
const jsonIDFilter = `{"filter": ["id", "=", "5cf24a90c4f75b00106a4c30"]}`
const jsonDateFilter = `{"filter": [["updatedAt", ">=", "2022-12-07T17:00:00Z"], "and", ["updatedAt", "<", "2022-12-08T17:00:00Z"]]}`

func TestFilter(t *testing.T) {
	opts := LoadOptions{}
	err := json.Unmarshal([]byte(jsonLoadOptions2), &opts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("opts: %v\n", opts.Filter)
	fmt.Printf("opts.pipeline: %v\n", opts.ParseFilter())
}

func TestFilterID(t *testing.T) {
	opts := LoadOptions{
		Field: map[string]Field{
			"id": {Name: "_id", Serialize: serialize.ID},
		},
	}
	err := json.Unmarshal([]byte(jsonIDFilter), &opts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("opts: %v\n", opts.Filter)
	fmt.Printf("opts.pipeline: %v\n", opts.ParseFilter())
}

func TestFilterDate(t *testing.T) {
	opts := LoadOptions{
		Field: map[string]Field{
			"updatedAt": {Serialize: serialize.Date},
		},
	}
	err := json.Unmarshal([]byte(jsonDateFilter), &opts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("opts: %v\n", opts.Filter)
	fmt.Printf("opts.pipeline: %v\n", opts.ParseFilter())
}
