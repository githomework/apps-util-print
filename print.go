package print

import (
	"io/ioutil"

	"os/exec"
	"strings"

	"github.com/alexbrainman/printer"
)

var (
	networkNames map[string]string
)

func init() {
	networkNames = map[string]string{}

	output, _ := exec.Command("powershell", "Get-Printer", "|", "select", "-exp", "name").Output()

	names := strings.Split(strings.ToLower(string(output)), "\n")
	for _, v := range names {
		if len(v) > 2 && v[0:2] == "\\\\" {

			parts := strings.Split(v, "\\")
			parts = strings.Fields(parts[len(parts)-1])
			networkNames[parts[0]] = strings.TrimSpace(v)
			networkNames[strings.ToUpper(parts[0])] = strings.TrimSpace(v)

		}

	}

}

func printDocument(printerName, documentName string, output []byte) error {
	p, err := printer.Open(printerName)
	if err != nil {
		return err
	}
	defer p.Close()

	err = p.StartRawDocument(documentName)
	if err != nil {
		return err
	}
	defer p.EndDocument()

	err = p.StartPage()
	if err != nil {
		return err
	}

	p.Write(output)

	return p.EndPage()
}

func Print(printerName string, path string, tryNetworkQueue bool) error {
	if printerName == "" || printerName == "dummy" {
		return nil
	}
	if tryNetworkQueue {
		n := networkNames[printerName]
		if n != "" {
			printerName = n
		}
	}
	output, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = printDocument(printerName, path, output)
	if err != nil {
		return err
	}

	return nil
}

func NetworkPrinters() map[string]string {
	return networkNames
}
