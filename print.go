package printing

import (
	"io/ioutil"

	"strings"

	"github.com/alexbrainman/printer"
)

var (
	networkNames map[string]string
)

func init() {

	networkNames = map[string]string{}

	/*	output, _ := exec.Command("powershell", "Get-Printer", "|", "select", "-exp", "name").Output()
		names := strings.Split(strings.ToLower(string(output)), "\n")
	*/

	names, _ := printer.ReadNames()

	for _, v := range names {
		v = strings.TrimSpace(v)
		if len(v) > 2 && v[0:2] == "\\\\" {
			v = strings.ToLower(v)
			parts := strings.Split(v, "\\")
			parts = strings.Fields(parts[len(parts)-1])
			networkNames[parts[0]] = v
			networkNames[strings.ToUpper(parts[0])] = v

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
	printer.ReadNames()

	return p.EndPage()
}

func PrintFile(printerName string, path string, tryNetworkQueue bool) error {
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

func PrintBytes(printerName string, document string, buf []byte, tryNetworkQueue bool) error {
	if printerName == "" || printerName == "dummy" {
		return nil
	}
	if tryNetworkQueue {
		n := networkNames[printerName]
		if n != "" {
			printerName = n
		}
	}

	err := printDocument(printerName, document, buf)
	if err != nil {
		return err
	}

	return nil
}

func NetworkPrinters() map[string]string {
	return networkNames
}
