package wsdl

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	gowsdl "github.com/hooklift/gowsdl"
)

// LoadWSDL reads a WSDL file from a local path or URL and unmarshals it
// into the gowsdl.WSDL structure.
func LoadWSDL(location string) (*gowsdl.WSDL, error) {
	loc, err := gowsdl.ParseLocation(location)
	if err != nil {
		return nil, err
	}

	locStr := loc.String()
	var data []byte
	if strings.HasPrefix(locStr, "http://") || strings.HasPrefix(locStr, "https://") {
		resp, err := http.Get(locStr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch WSDL '%s': %w", locStr, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to fetch WSDL '%s': status %d, body: %s", locStr, resp.StatusCode, string(body))
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		b, err := os.ReadFile(locStr)
		if err != nil {
			return nil, err
		}
		data = b
	}

	ws := new(gowsdl.WSDL)
	if err := xml.Unmarshal(data, ws); err != nil {
		return nil, err
	}
	return ws, nil
}

// BuildSOAPEnvelope creates a simple SOAP 1.1 envelope with the given
// operation name and parameters. The namespace will be used as the
// prefix 'tns' if provided.
func BuildSOAPEnvelope(namespace, operation string, params map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	sb.WriteString(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"`)
	if namespace != "" {
		sb.WriteString(` xmlns:tns="` + namespace + `"`)
	}
	sb.WriteString(`>`)
	sb.WriteString(`<soapenv:Body>`)
	if namespace != "" {
		sb.WriteString(`<tns:` + operation + `>`)
	} else {
		sb.WriteString(`<` + operation + `>`)
	}
	for k, v := range params {
		sb.WriteString("<" + k + ">" + fmt.Sprintf("%v", v) + "</" + k + ">")
	}
	if namespace != "" {
		sb.WriteString(`</tns:` + operation + `>`)
	} else {
		sb.WriteString(`</` + operation + `>`)
	}
	sb.WriteString(`</soapenv:Body></soapenv:Envelope>`)
	return sb.String()
}
