package html2xlsx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nad2000/excelize"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"golang.org/x/net/html"
)

// CellAddress maps a cell coordiantes (row, column) to its address
func CellAddress(rowIndex, colIndex int) string {
	return xlsx.GetCellIDStringFromCoords(colIndex, rowIndex)
}

func convert(r io.Reader, file *excelize.File) {

	doc, err := html.Parse(r)
	if err != nil {
		log.WithError(err).Errorln("Failed to read the file content.")
		return
	}
	var f func(n *html.Node)
	var row int
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			// Do something with n...
			var col int
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Data != "th" && c.Data != "td" {
					continue
				}
				colspan := 1
				if c.Attr != nil {
					for _, a := range c.Attr {
						if a.Key == "colspan" {
							colspan, _ = strconv.Atoi(a.Val)
							break
						}
					}
				}
				// fmt.Printf("COLSPAN %d, VAL: %#v", colspan, c.FirstChild.Data)
				addr := CellAddress(row, col)
				if c.FirstChild.Data != "" {
					file.SetCellValue("Sheet1", addr, c.FirstChild.Data)
				}
				col += colspan
			}
			row++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return

}

// Convert opens zip file and goes throug all the .xls (with HTML content)
// convrerts them into .xlsx and stores into a new output zipfile
func Convert(filename, outputFilename string) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.WithError(err).Errorln("Failed to open the input file.")
	}
	defer r.Close()

	output, err := os.Create(outputFilename)
	if err != nil {
		log.WithError(err).Errorln("Failed to create output file.")
		return
	}
	defer output.Close()

	writer := zip.NewWriter(output)
	defer writer.Close()

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
		file := excelize.NewFile()
		convert(rc, file)
		of, err := writer.Create(strings.TrimSuffix(f.Name, filepath.Ext(f.Name)) + ".xlsx")
		if err != nil {
			log.WithError(err).Errorln("Failed to create a writer for a single file.")
			continue
		}
		file.Write(of)
		// file.SaveAs(strings.TrimSuffix(f.Name, filepath.Ext(f.Name)) + ".xlsx")

		rc.Close()
	}
}
