package webinfo

import (
	"fmt"
	"testing"
)

func TestTemplateStore_Create(t *testing.T) {
	s := newTemplateStore(WebTemplateStorePath)
	s.Get("A")
	s.Get("B")
	s.Modify("A", "33333")
	s.Save(WebTemplateStorePath)
}

func TestTemplateStore_Load(t *testing.T) {
	s := newTemplateStore(WebTemplateStorePath)
	fmt.Println(s.List())
	for name, tmpl := range s.Templates {
		fmt.Println(name, tmpl.Template)
	}
}
