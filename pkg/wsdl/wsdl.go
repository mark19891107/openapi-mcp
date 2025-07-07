package wsdl

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Definitions represents the top-level WSDL definitions element.
type Definitions struct {
	XMLName         xml.Name   `xml:"definitions"`
	Name            string     `xml:"name,attr"`
	TargetNamespace string     `xml:"targetNamespace,attr"`
	Messages        []Message  `xml:"message"`
	PortTypes       []PortType `xml:"portType"`
	Bindings        []Binding  `xml:"binding"`
	Services        []Service  `xml:"service"`
}

type Message struct {
	Name  string `xml:"name,attr"`
	Parts []Part `xml:"part"`
}

type Part struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}

type PortType struct {
	Name       string              `xml:"name,attr"`
	Operations []PortTypeOperation `xml:"operation"`
}

type PortTypeOperation struct {
	Name          string           `xml:"name,attr"`
	Documentation string           `xml:"documentation"`
	Input         OperationMessage `xml:"input"`
	Output        OperationMessage `xml:"output"`
}

type OperationMessage struct {
	Message string `xml:"message,attr"`
}

type Binding struct {
	Name       string             `xml:"name,attr"`
	Type       string             `xml:"type,attr"`
	Operations []BindingOperation `xml:"operation"`
}

type BindingOperation struct {
	Name       string `xml:"name,attr"`
	SoapAction string `xml:"operation>soapAction,attr"`
}

type Service struct {
	Name  string        `xml:"name,attr"`
	Ports []ServicePort `xml:"port"`
}

type ServicePort struct {
	Name    string  `xml:"name,attr"`
	Binding string  `xml:"binding,attr"`
	Address Address `xml:"address"`
}

type Address struct {
	Location string `xml:"location,attr"`
}

// LoadWSDL loads a WSDL file from the provided location which can be a local path or a URL.
func LoadWSDL(location string) (*Definitions, error) {
	var reader io.ReadCloser
	var err error
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		u, errURL := http.Get(location)
		if errURL != nil {
			return nil, errURL
		}
		reader = u.Body
	} else {
		reader, err = os.Open(location)
		if err != nil {
			return nil, err
		}
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var defs Definitions
	if err := xml.Unmarshal(data, &defs); err != nil {
		return nil, err
	}
	return &defs, nil
}

// BuildSOAPEnvelope creates a basic SOAP 1.1 envelope for the given operation and parameters.
func BuildSOAPEnvelope(operation string, params map[string]interface{}) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	b.WriteString(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">`)
	b.WriteString(`<soap:Body>`)
	b.WriteString(`<` + operation + `>`)
	for k, v := range params {
		b.WriteString(`<` + k + `>` + fmt.Sprintf("%v", v) + `</` + k + `>`)
	}
	b.WriteString(`</` + operation + `>`)
	b.WriteString(`</soap:Body></soap:Envelope>`)
	return b.String()
}
