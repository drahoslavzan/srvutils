package dbutils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Distinct[T any](mc *mongo.Collection, col string) []T {
	return DistinctFilter[T](mc, col, bson.M{})
}

func DistinctFilter[T any](mc *mongo.Collection, col string, filter interface{}) []T {
	res, err := mc.Distinct(context.Background(), col, filter)
	if err != nil {
		panic(err)
	}

	// NOTE: MongoDB can return unefined value if the array value for the indexed field is empty
	vals := make([]T, 0, len(res))
	for _, l := range res {
		if v, ok := l.(T); ok {
			vals = append(vals, v)
		}
	}

	return vals
}
