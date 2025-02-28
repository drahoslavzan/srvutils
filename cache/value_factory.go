package cache

type (
	stringValueFactory struct{}
)

func StringValueFactory() stringValueFactory {
	return stringValueFactory{}
}

func (m stringValueFactory) FromString(v string) (*string, error) {
	return &v, nil
}
