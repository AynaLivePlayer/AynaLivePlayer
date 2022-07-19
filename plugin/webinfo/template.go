package webinfo

import (
	"AynaLivePlayer/util"
	"encoding/json"
	"io/ioutil"
)

const WebTemplateStorePath = "./webtemplates.json"

type WebTemplate struct {
	Name     string
	Template string
}

type TemplateStore struct {
	Templates map[string]*WebTemplate
}

func newTemplateStore(filename string) *TemplateStore {
	s := &TemplateStore{Templates: map[string]*WebTemplate{}}
	var templates []WebTemplate
	file, err := ioutil.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(file, &templates)
	}
	for _, tmpl := range templates {
		s.Templates[tmpl.Name] = &tmpl
	}
	return s
}

func (s *TemplateStore) Save(filename string) {
	templates := make([]WebTemplate, 0)
	for _, tmp := range s.Templates {
		templates = append(templates, *tmp)
	}
	unescape, err := util.MarshalIndentUnescape(templates, "", "    ")
	if err != nil {
		lg.Warnf("save web templates to %s failed: %s", filename, err)
		return
	}
	if err := ioutil.WriteFile(filename, []byte(unescape), 0666); err != nil {
		lg.Warnf("save web templates to %s failed: %s", filename, err)
		return
	}
}

func (s *TemplateStore) Get(name string) *WebTemplate {
	if t, ok := s.Templates[name]; ok {
		return t
	}
	t := &WebTemplate{Name: name, Template: "<p>Empty</p>"}
	s.Templates[name] = t
	return t
}

func (s *TemplateStore) Modify(name string, content string) {
	if _, ok := s.Templates[name]; ok {
		s.Templates[name].Template = content
		return
	}
}

func (s *TemplateStore) List() []string {
	names := make([]string, 0)
	for name, _ := range s.Templates {
		names = append(names, name)
	}
	return names
}

func (s *TemplateStore) Delete(name string) {
	delete(s.Templates, name)
}
