package template

import (
	"github.com/GoComply/fedramp/bundled"
	"github.com/GoComply/fedramp/pkg/fedramp/common"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
	"github.com/opencontrol/doc-template/docx"
	"io"
	"io/ioutil"
	"os"
)

type Template struct {
	wordDoc *docx.Docx
	xmlDoc  *xml.XmlDocument
}

func NewTemplate(level common.BaselineLevel) (*Template, error) {
	in, err := bundled.TemplateDOCX(level)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	out, err := ioutil.TempFile("/tmp", "FedRAMP-"+level.Name())
	if err != nil {
		return nil, err
	}
	defer out.Close()
	defer os.Remove(out.Name())

	_, err = io.Copy(out, in)
	if err != nil {
		return nil, err
	}

	return NewTemplateFile(out.Name())
}

func NewTemplateFile(filePath string) (*Template, error) {
	var t Template
	t.wordDoc = new(docx.Docx)
	err := t.wordDoc.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = t.parseDocxXml()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Template) parseDocxXml() error {
	var err error
	t.xmlDoc, err = gokogiri.ParseXml([]byte(t.wordDoc.GetContent()))
	if err != nil {
		return err
	}
	xp := t.xmlDoc.DocXPathCtx()
	xp.RegisterNamespace("w", "http://schemas.openxmlformats.org/wordprocessingml/2006/main")
	xp.RegisterNamespace("w14", "http://schemas.microsoft.com/office/word/2010/wordml")
	return nil
}
