package dbutils

import "go.mongodb.org/mongo-driver/bson"

func ToBson(s interface{}, del ...string) bson.M {
	data, err := bson.Marshal(s)
	if err != nil {
		panic(err)
	}

	var ret bson.M
	err = bson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	for _, d := range del {
		delete(ret, d)
	}

	return ret
}

func ToBsonNoID(s interface{}) bson.M {
	return ToBson(s, "_id")
}
