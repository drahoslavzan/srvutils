package options

import (
	"encoding/json"
	"fmt"
	"testing"
)

const jsonBinaryFilter = `["field", "=", 10]`
const jsonUnaryFilter = `["!", ` + jsonBinaryFilter + `]`
const jsonLoadOptions1 = `{"filter":` + jsonBinaryFilter + `}`
const jsonLoadOptions2 = `{"filter":` + jsonUnaryFilter + `}`

// TODO: implement

func TestFilter(t *testing.T) {
	tmp := LoadOptions{}
	err := json.Unmarshal([]byte(jsonLoadOptions2), &tmp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("opts = %v\n", tmp.Filter)
	fmt.Printf("opts.pipeline = %v\n", tmp.ParseFilter())
}
