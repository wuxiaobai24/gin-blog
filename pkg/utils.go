package pkg

import (
	"fmt"
	"html/template"
	"time"

	"github.com/russross/blackfriday"
)

// ToDate template func
func ToDate(t time.Time) string {
	return fmt.Sprintf("%v-%v-%v", t.Year(), t.Month(), t.Day())
}

// Markdown2HTML template func
func Markdown2HTML(content string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(content)))
}
