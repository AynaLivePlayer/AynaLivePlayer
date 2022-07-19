package webinfo

import (
	"fmt"
	"testing"
)

func TestTemplateStore_Create(t *testing.T) {
	s := newTemplateStore(WebTemplateStorePath)
	s.Get("A")
	s.Get("B")
	s.Modify("A", "123123")
	s.Save(WebTemplateStorePath)
}

func TestTemplateStore_Load(t *testing.T) {
	s := newTemplateStore(WebTemplateStorePath)
	for name, tmpl := range s.Templates {
		fmt.Println(name, tmpl.Template)
	}
}
