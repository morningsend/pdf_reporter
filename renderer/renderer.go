package renderer

import (
	"bytes"
	"html/template"
	"io"
)

type Renderer struct {
	//tempplate    *template.Template
	templateFile string
}

func NewRenderer(templateFile string) (*Renderer, error) {

	return &Renderer{
		templateFile: templateFile,
	}, nil
}

func (r *Renderer) getTemplate() (*template.Template, error) {
	return template.ParseFiles(r.templateFile)
}

func (r *Renderer) RenderToString(data interface{}) (string, error) {
	temp, err := r.getTemplate()
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBuffer(nil)
	if err := temp.Execute(buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (r *Renderer) Render(writer io.Writer, data interface{}) error {
	temp, err := r.getTemplate()
	if err != nil {
		return err
	}
	return temp.Execute(writer, data)
}
