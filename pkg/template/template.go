// Copies from https://github.com/ray-g/go-bindata-template/blob/master/template.go
package template

import (
	"html/template"
)

// BinData caches the functions of generated bindata.go needed by this package
type BinData struct {
	// Asset returns the asset content in bytes
	Asset func(path string) ([]byte, error)
	// AssetDir returns the file names below a certain path
	AssetDir func(path string) ([]string, error)
	// AssetNames returns the names of the assets
	AssetNames func() []string
}

// Template extends Golang's html/template with go-bindata
type Template struct {
	*BinData
	*template.Template
}

// Must wraps Golang's html/template Must
func Must(t *template.Template, err error) *template.Template {
	return template.Must(t, err)
}

// New creates a new Template
func New(name string, data *BinData) *Template {
	return &Template{data, template.New(name)}
}

// Parse loads the given file and parse the content of it
func (t *Template) Parse(filename string) (*template.Template, error) {
	loaded, err := t.load(filename)
	if err != nil {
		return nil, err
	}
	tmpl, err := t.Template.Parse(string(loaded))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// ParseFiles loads of all given files, then parse concatenated contents of these files
func (t *Template) ParseFiles(filenames ...string) (*template.Template, error) {
	content := []byte{}
	for _, filename := range filenames {
		loaded, err := t.load(filename)
		if err != nil {
			return nil, err
		}
		content = append(content, loaded...)
	}
	tmpl, err := t.Template.Parse(string(content))

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// ParseDir loads all files under the dir, then parse concatenated contents of these files
func (t *Template) ParseDir(dir string) (*template.Template, error) {
	filenames, err := t.AssetDir(dir)
	if err != nil {
		return nil, err
	}

	var filepaths []string
	for _, filename := range filenames {
		filepaths = append(filepaths, dir+"/"+filename)
	}

	return t.ParseFiles(filepaths...)
}

// ParseAll loads all files in the embeded asset, then parse concatenated contents of these files
func (t *Template) ParseAll() (*template.Template, error) {
	assets := t.AssetNames()
	return t.ParseFiles(assets...)
}

// load loads the file content
func (t *Template) load(file string) ([]byte, error) {
	bytes, err := t.Asset(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
