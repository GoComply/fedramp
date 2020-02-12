package template

import (
	"github.com/GoComply/fedramp/bundled"
	"github.com/GoComply/fedramp/pkg/fedramp/common"
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

	return &t, nil
}
