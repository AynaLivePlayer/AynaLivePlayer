package model

import "fmt"

type Meta struct {
	Name string
	Id   string
}

func (m Meta) String() string {
	return fmt.Sprintf("<Meta %s:%s>", m.Name, m.Id)
}

func (m Meta) Identifier() string {
	return fmt.Sprintf("%s_%s", m.Name, m.Id)
}
