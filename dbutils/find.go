package dbutils

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[T any](logger log.Logger, col *mongo.Collection, filter bson.M) *T {
	res := col.FindOne(context.Background(), filter)
	return SingleResultValue[T](logger, res)
}

func UpdateOne[T any](logger log.Logger, col *mongo.Collection, filter, update bson.M) *T {
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	res := col.FindOneAndUpdate(context.Background(), filter, update, opts)
	return SingleResultValue[T](logger, res)
}

func SingleResultValue[T any](logger log.Logger, res *mongo.SingleResult) *T {
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		logger.Panic(err)
	}

	var ret T
	if err := res.Decode(&ret); err != nil {
		logger.Panic(err)
	}

	return &ret
}
