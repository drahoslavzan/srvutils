package options

type Sort struct {
	Selector string `json:"selector"`
	Desc     bool   `json:"desc"`
}

type LoadOptions struct {
	Skip   int64  `json:"skip"`
	Take   int64  `json:"take"`
	Filter Filter `json:"filter"`
	Sort   []Sort `json:"sort"`
}

func (m *Sort) GetField() string {
	if m.Selector == "id" {
		return "_id"
	}
	return m.Selector
}

func (m *Sort) GetOrder() int {
	if m.Desc {
		return -1
	}
	return 1
}
