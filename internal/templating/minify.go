package templating

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func CompileTemplates(pattern string, funcMap template.FuncMap) (*template.Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}

	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: false,
		KeepWhitespace:      false,
		KeepDocumentTags:    false,
	})

	var tmpl *template.Template
	for _, filename := range filenames {
		name := filepath.Base(filename)
		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}
		tmpl.Funcs(funcMap)
		tmpl, err = tmpl.Parse(string(mb))
		if err != nil {
			fmt.Println(err)
		}
	}
	return tmpl, nil
}
