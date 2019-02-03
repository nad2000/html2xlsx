package html2xlsx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nad2000/excelize"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"golang.org/x/net/html"
)

const sheet = "Sheet1"

// CellAddress maps a cell coordiantes (row, column) to its address
func CellAddress(rowIndex, colIndex int) string {
	return xlsx.GetCellIDStringFromCoords(colIndex, rowIndex)
}

func setStyle(file *excelize.File, addr string, format int) {
	style, _ := file.NewStyle(`{"number_format": ` + strconv.Itoa(format) + `}`)
	file.SetCellStyle(sheet, addr, addr, style)
}
func setFloat(file *excelize.File, addr, value string, format int) {
	numValue := strings.Replace(value, ",", "", -1)

	if floatVal, err := strconv.ParseFloat(numValue, 64); err == nil {
		if format == 9 || format == 10 {
			file.SetCellDefault(sheet, addr, strconv.FormatFloat(floatVal/100.0, 'f', -1, 64))
		} else {
			file.SetCellDefault(sheet, addr, numValue)
		}
		setStyle(file, addr, format)
	} else {
		file.SetCellValue(sheet, addr, value)
	}
}

func timeValToExcelTimeVal(t time.Time) string {
	return strconv.FormatFloat(float64(t.UnixNano())/8.64e13+25569.0, 'f', -1, 64)
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
				if c.FirstChild.Data != "" {
					addr := CellAddress(row, col)
					value := c.FirstChild.Data
					if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
						file.SetCellValue(sheet, addr, intVal)
					} else if strings.Contains(value, "/") {
						timeVal, err := time.Parse("01/02/2006", value)
						if err != nil {
							log.WithError(err).Errorf("Failed to parse: %q", value)
							file.SetCellValue(sheet, addr, value)

						} else {
							// need to reimplement for dates...
							file.SetCellDefault(sheet, addr, timeValToExcelTimeVal(timeVal))
							setStyle(file, addr, 14)
						}
					} else if strings.HasPrefix(value, "$") {
						setFloat(file, addr, value[1:], 165)
					} else if strings.HasSuffix(value, "%") {
						if strings.Contains(value, ".") {
							setFloat(file, addr, value[:len(value)-1], 10)
						} else {
							setFloat(file, addr, value[:len(value)-1], 9)
						}
					} else if matched, err := regexp.MatchString("^[0-9,]+\\.?\\d*$", value); matched && err == nil {
						setFloat(file, addr, value, 4)
					} else {
						file.SetCellValue(sheet, addr, value)
					}
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
		newName := strings.TrimSuffix(f.Name, filepath.Ext(f.Name)) + ".xlsx"
		of, err := writer.CreateHeader(
			&zip.FileHeader{
				Name:         newName,
				Method:       zip.Deflate,
				Modified:     f.Modified,
				ModifiedTime: f.ModifiedTime,
				ModifiedDate: f.ModifiedDate,
			})

		if err != nil {
			log.WithError(err).Errorln("Failed to create a writer for a single file.")
			continue
		}
		file.Write(of)
		// file.SaveAs(newName)

		rc.Close()
	}
}
