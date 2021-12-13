package options

type Sort struct {
	Selector string `json:"selector"`
	Desc     bool   `json:"desc"`
}

type Field struct {
	Name   string `json:"name"`
	IsDate bool   `json:"isDate"`
	IsID   bool   `json:"isID"`
}

type LoadOptions struct {
	Skip   int64            `json:"skip"`
	Take   int64            `json:"take"`
	Filter Filter           `json:"filter"`
	Sort   []Sort           `json:"sort"`
	Field  map[string]Field `json:"field"`
}

func (m *Sort) GetField(opts *LoadOptions) string {
	if opts.Field != nil {
		if f, ok := opts.Field[m.Selector]; ok && len(f.Name) > 0 {
			return f.Name
		}
	}
	return m.Selector
}

func (m *Sort) GetOrder() int {
	if m.Desc {
		return -1
	}
	return 1
}
