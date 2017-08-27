package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadConfig(target interface{}, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func LoadTemplates(tplRoot string, fmap template.FuncMap) (*template.Template, error) {
	files := []string{}
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".html") {
			files = append(files, path)
		}
		return nil
	}

	err := filepath.Walk(tplRoot, walk)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no .html templates found in %s", tplRoot)
	}

	tpls := template.Must(template.New("").Funcs(fmap).ParseFiles(files...))

	return tpls, nil
}
