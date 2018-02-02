package app

type PrintValue struct {
	Label string
	Value interface{}
}

func (s PrintValue) GetLabel() string {
	return s.Label
}

func (s PrintValue) GetValue() interface{} {
	return s.Value
}
