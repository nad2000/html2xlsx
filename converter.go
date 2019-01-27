package html2xlsx

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func convert(r io.Reader) {

	doc, err := html.Parse(r)
	if err != nil {
		log.WithError(err).Errorln("Failed to read the file content.")
		return
	}
	// var table Table
	// xml.Unmarshal(byteValue, &table)
	// fmt.Printf("%#v\n\n\n", table)
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "td" {
			// Do something with n...
			fmt.Printf("%#v\n", n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

}

// Convert opens zip file and goes throug all the .xls (with HTML content)
// convrerts them into .xlsx and stores into a new output zipfile
func Convert(fileName, outputFilename string) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "__MACOSX") || filepath.Ext(f.Name) != ".xls" {
			continue
		}
		log.WithField("Filename", f.Name).Infoln("Converting...")
		rc, err := f.Open()
		if err != nil {
			log.WithError(err).Errorln("Failed to open file in the archive file.")
			continue
		}
		convert(rc)

		// _, err = io.CopyN(os.Stdout, rc, 68)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		rc.Close()
		// fmt.Println()
	}
}
