package dbutils

import "go.mongodb.org/mongo-driver/bson"

func ToBson(s interface{}) bson.M {
	data, err := bson.Marshal(s)
	if err != nil {
		panic(err)
	}

	var replace bson.M
	err = bson.Unmarshal(data, &replace)
	if err != nil {
		panic(err)
	}

	return replace
}
