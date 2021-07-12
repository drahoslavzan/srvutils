package dbutils

import "go.mongodb.org/mongo-driver/bson"

func ToBson(s interface{}) bson.M {
	data, err := bson.Marshal(s)
	if err != nil {
		panic(err)
	}

	var ret bson.M
	err = bson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return ret
}

func ToBsonNoID(s interface{}) bson.M {
	ret := ToBson(s)
	delete(ret, "_id")
	return ret
}
