package prompts

import (
	"bytes"
	"text/template"
)

type PromptTemplate struct {
	Template *template.Template
}

// creates a new prompt template from a string
func NewPromptTemplate(templateString string) (*PromptTemplate, error) {
	tmpl, err := template.New("prompt").Parse(templateString)
	if err != nil {
		return nil, err
	}
	return &PromptTemplate{Template: tmpl}, nil
}

// populates the template with the given data.
func (pt *PromptTemplate) Format(data any) (string, error) {
	var buf bytes.Buffer
	if err := pt.Template.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}