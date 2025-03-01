package cache

type (
	stringValueFactory struct{}
)

func StringValueFactory() ValueFactory[string] {
	return stringValueFactory{}
}

func (m stringValueFactory) FromString(v string) (string, error) {
	return v, nil
}

func (m stringValueFactory) ToString(v string) string {
	return v
}
