package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"os"
	"path/filepath"
)

const (
	header = ` <!DOCTYPE html> <html lang="en"> <head> <meta charset="UTF-8"> <meta name="viewport" content="width=device-width, initial-scale=1"> <title>Markdown Preview Tool</title> </head> <body>`
	footer = `</body></html>`
)

func main() {
	// parse flags
	filename := flag.String("file", "", "Markdown file to Preview")
	flag.Parse()
	// si el usuario no provee un archivo de entrada, mostrar ayuda
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	// lee los datos del archivo de entrada
	// y verifica los errores
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	// parsear el archivo markdown con
	// blackfriday y bluemonday para generar
	// un archivo HTML seguro
	ouput := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(ouput)

	// crear un buffer de bytes para escribir el contenido
	// del archivo
	var buffer bytes.Buffer

	// escribir html al buffer de bytes
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFname string, data []byte) error {
	// write the bytes to the file
	return os.WriteFile(outFname, data, 0644)
}
