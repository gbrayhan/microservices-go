package tools

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func GeneratePDF(data interface{}, templateName string) (pdfContent string) {
	_ = os.Setenv("WKHTMLTOPDF_PATH", "/usr/bin/")

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	resultB := ProcessFile(templateName, data)
	page := wkhtmltopdf.NewPageReader(strings.NewReader(resultB))

	pdfg.AddPage(page)
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	filename := "archives/tmp/" + strconv.Itoa(rand.Int()) + "_file.pdf"
	err = pdfg.WriteFile(filename)
	var w io.Writer
	pdfg.SetOutput(w)

	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}
	pdfContent = string(content)
	return
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

func ProcessFile(fileName string, vars interface{}) string {
	tmpl, err := template.ParseFiles(fileName)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}
