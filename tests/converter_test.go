package main

import "testing"
import "github.com/nad2000/html2xlsx"

func TestHelloWorld(t *testing.T) {
	html2xlsx.Convert("jan.zip", "output.zip")
}
