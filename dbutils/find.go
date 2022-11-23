package dbutils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[T any](col *mongo.Collection, filter bson.M) *T {
	res := col.FindOne(context.Background(), filter)
	return SingleResultValue[T](res)
}

func UpdateOne[T any](col *mongo.Collection, filter, update bson.M) *T {
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	res := col.FindOneAndUpdate(context.Background(), filter, update, opts)
	return SingleResultValue[T](res)
}

func SingleResultValue[T any](res *mongo.SingleResult) *T {
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		panic(err)
	}

	var ret T
	if err := res.Decode(&ret); err != nil {
		panic(err)
	}

	return &ret
}
