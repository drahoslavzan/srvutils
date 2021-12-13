package serialize

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ID(val interface{}) interface{} {
	sv, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("invalid object id: %v", val))
	}
	ret, err := primitive.ObjectIDFromHex(sv)
	if err != nil {
		panic(fmt.Errorf("invalid object id: %v", sv))
	}
	return ret
}

func Date(val interface{}) interface{} {
	sv, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("invalid date provided: %v", val))
	}
	ret, err := time.Parse(time.RFC3339, sv)
	if err != nil {
		panic(fmt.Errorf("invalid date provided: %v", sv))
	}
	return ret
}
