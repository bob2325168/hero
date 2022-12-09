package demo

import (
	"github.com/bob2325168/gohero/framework"
)

type Service struct {
	c framework.Container
}

func NewService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	return &Service{c: c}, nil
}

func (s *Student) GetAllStudent() []Student {
	return []Student{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   3,
			Name: "bar",
		},
	}
}
