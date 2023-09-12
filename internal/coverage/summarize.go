package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mlctrez/godom"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) == 1 {
		log.Fatal("provide coverage.html file as first arg")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("error reading %s : %s", os.Args[1], err)
	}

	onlyBody := &bytes.Buffer{}
	scanner := bufio.NewScanner(file)
	var inBody bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "<body>" {
			inBody = true
		}
		if inBody {
			onlyBody.WriteString(line)
			onlyBody.WriteString("\n")
		}
		if strings.TrimSpace(line) == "</body>" {
			inBody = false
		}
	}

	document := godom.Document()
	api := document.DocApi()
	h := api.H(onlyBody.String())
	for _, element := range h.GetElementsByTagName("option") {
		optionText := element.ChildNodes()[0].String()
		if !strings.HasSuffix(optionText, "(100.0%)") {
			fmt.Println(optionText)
		}
	}

}
