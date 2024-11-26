package models

import (
	"html/template"
	"io"
	"log"
	"time"
)

type TemplateBlog struct {
	*template.Template
}

type HtmlTemplate struct {
	Index      TemplateBlog
	Category   TemplateBlog
	Custom     TemplateBlog
	Detail     TemplateBlog
	Login      TemplateBlog
	Pigeonhole TemplateBlog
	Writing    TemplateBlog
}

func (t *TemplateBlog) WriteData(w io.Writer, data interface{}) {
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte("write template data error"))
	}
}

func InitTemplate(templateDir string) (HtmlTemplate, error) {

	tp, err := readTemplate(
		[]string{"index", "category", "custom", "detail", "login", "pigeonhole", "writing"},
		templateDir,
	)
	var htmlTemplate HtmlTemplate

	if err != nil {
		return htmlTemplate, err
	}
	htmlTemplate.Index = tp[0]
	htmlTemplate.Category = tp[1]
	htmlTemplate.Custom = tp[2]
	htmlTemplate.Detail = tp[3]
	htmlTemplate.Login = tp[4]
	htmlTemplate.Pigeonhole = tp[5]
	htmlTemplate.Writing = tp[6]

	return htmlTemplate, nil
}
func IsODD(num int) bool {
	return num%2 == 0
}
func getNextName(strs []string, index int) string {
	return strs[index+1]
}
func Date(layout string) string {
	return time.Now().Format(layout)
}
func DateDay(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}
func readTemplate(templates []string, templateDir string) ([]TemplateBlog, error) {
	var templateBlogs []TemplateBlog
	for _, view := range templates {
		viewName := view + ".html"
		t := template.New(viewName)
		index := templateDir + viewName
		home := templateDir + "home.html"
		header := templateDir + "layout/header.html"
		footer := templateDir + "layout/footer.html"
		pagination := templateDir + "layout/pagination.html"
		personal := templateDir + "layout/personal.html"
		postList := templateDir + "layout/post-list.html"
		t.Funcs(template.FuncMap{
			"isODD":       IsODD,
			"getNextName": getNextName,
			"date":        Date,
			"dateDay":     DateDay,
		})
		files, err := t.ParseFiles(index, home, header, footer, pagination, personal, postList)
		if err != nil {
			log.Printf("Template Parse Error: %s", err)
			return nil, err
		}
		var tb TemplateBlog
		tb.Template = files
		templateBlogs = append(templateBlogs, tb)
	}
	return templateBlogs, nil
}
