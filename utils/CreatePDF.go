// Package utils contains the common functions and structures for the application
package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/viper"
)

// RouteTemplateToPDF Convert a golang template to PDF
func RouteTemplateToPDF(routeTemplate string, data interface{}) (pdfContent string, err error) {
	viper.SetConfigFile("config.json")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	wkhtmltopdfBin := viper.GetString("Tools.WkhtmltopdfBin")
	randCryp, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return
	}

	filenamePDF := "archives/tmp/" + randCryp.String() + "_file.pdf"
	filenameHTML := "archives/tmp/" + randCryp.String() + "_file.html"

	file, err := os.Create(filenameHTML)
	if err != nil {
		return
	}

	htmlTemplate := processFile(routeTemplate, data)
	if _, err = file.WriteString(htmlTemplate); err != nil {
		return
	}

	if err = file.Close(); err != nil {
		return
	}

	args := []string{"-s", "Letter", "-O", "Portrait", filenameHTML, filenamePDF}
	cmd := exec.Command(wkhtmltopdfBin, args...)
	// vars outString, err
	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	content, err := os.ReadFile(filenamePDF)
	if err != nil {
		return
	}

	if err = os.Remove(filenamePDF); err != nil {
		return
	}
	if err = os.Remove(filenameHTML); err != nil {
		return
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

func processFile(rutaFileName string, vars interface{}) string {
	tmpl, err := template.ParseFiles(rutaFileName)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}
