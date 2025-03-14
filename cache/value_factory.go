package cache

import "encoding/json"

type (
	stringValueFactory struct{}

	jsonValueFactory[T any] struct{}
)

func StringValueFactory() ValueFactory[string] {
	return stringValueFactory{}
}

func (m stringValueFactory) FromString(v string) (string, error) {
	return v, nil
}

func (m stringValueFactory) ToString(v string) (string, error) {
	return v, nil
}

func JSONValueFactory[T any]() ValueFactory[T] {
	return jsonValueFactory[T]{}
}

func (m jsonValueFactory[T]) FromString(v string) (T, error) {
	var value T
	err := json.Unmarshal([]byte(v), &value)
	return value, err
}

func (m jsonValueFactory[T]) ToString(v T) (string, error) {
	value, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(value), nil
}
