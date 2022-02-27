package dbutils

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

const DupErrCode = 11000

func IsMongoDupWriteErr(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == DupErrCode {
				return true
			}
		}
	}
	return false
}

func IsMongoDupBulkWriteErr(err error) bool {
	var e mongo.BulkWriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == DupErrCode {
				return true
			}
		}
	}
	return false
}
