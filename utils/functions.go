package utils

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// original post: https://stackoverflow.com/a/38688083/16158590
func parseTemplates() *template.Template {
	t := template.New("")
	err := filepath.Walk("./public/views/", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = t.ParseFiles(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return err
	})
	if err != nil {
		panic(err)
	}

	return t
}

var ParsedTemplates = parseTemplates()

func DecodeJSON(res *http.Response) map[string]interface{} {
	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	return response
}
