package options

type Group struct {
	Selector string `json:"selector"`
	Desc     bool   `json:"desc"`
}

func (m *Group) GetOrder() int {
	if m.Desc {
		return -1
	}
	return 1
}
