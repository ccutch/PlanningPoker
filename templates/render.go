package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
)

//go:embed *.html
var templates embed.FS

var (
	t, err = template.New("").Funcs(template.FuncMap{
		"inc": func(n int) int {
			return n + 1
		},
	}).ParseFS(templates, "*.html")
	views = template.Must(t, err)
)

func Render(w io.Writer, name string, data interface{}) {
	var buf bytes.Buffer
	if err := views.ExecuteTemplate(&buf, name, data); err != nil {
		log.Println("Error rendering:", err)
		return
	}
	fmt.Fprint(w, strings.ReplaceAll(buf.String(), "\n", ""))
}
