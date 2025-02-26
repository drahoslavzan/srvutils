package cache

type (
	StringValueFactory struct{}
)

func NewStringValueFactory() *StringValueFactory {
	return &StringValueFactory{}
}

func (m *StringValueFactory) FromString(v string) (*string, error) {
	return &v, nil
}
