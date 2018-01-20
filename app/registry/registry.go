package registry

import "errors"

func Create() *Registry {
	return &Registry{
		items: make(map[string]interface{}),
	}
}

type Registry struct {
	items map[string]interface{}
}

func (s *Registry) Add(name string, c interface{}) {
	if _, exists := s.items[name]; exists == true {
		panic(errors.New("object is already registered: " + name))
	}
	s.items[name] = c
}

func (s *Registry) GetByName(name string) interface{} {
	item, ok := s.items[name]
	if ok == false {
		panic(errors.New("object is not registered: " + name))
	}
	return item
}
