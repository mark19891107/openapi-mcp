package wsdl

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Definitions represents the root of a WSDL document.
type Definitions struct {
	XMLName         xml.Name   `xml:"definitions"`
	TargetNamespace string     `xml:"targetNamespace,attr"`
	Services        []Service  `xml:"service"`
	PortTypes       []PortType `xml:"portType"`
	Bindings        []Binding  `xml:"binding"`
	Messages        []Message  `xml:"message"`
}

type Service struct {
	Name  string `xml:"name,attr"`
	Ports []Port `xml:"port"`
}

type Port struct {
	Name    string  `xml:"name,attr"`
	Binding string  `xml:"binding,attr"`
	Address Address `xml:"address"`
}

type Address struct {
	Location string `xml:"location,attr"`
}

type PortType struct {
	Name       string              `xml:"name,attr"`
	Operations []PortTypeOperation `xml:"operation"`
}

type PortTypeOperation struct {
	Name   string            `xml:"name,attr"`
	Input  *OperationMessage `xml:"input"`
	Output *OperationMessage `xml:"output"`
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
	Name   string        `xml:"name,attr"`
	SOAPOp SOAPOperation `xml:"operation"`
}

type SOAPOperation struct {
	SOAPAction string `xml:"soapAction,attr"`
	Style      string `xml:"style,attr"`
}

type Message struct {
	Name  string `xml:"name,attr"`
	Parts []Part `xml:"part"`
}

type Part struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Element string `xml:"element,attr"`
}

// LoadWSDL reads a WSDL from a file path or URL and parses it.
func LoadWSDL(location string) (*Definitions, error) {
	var reader io.ReadCloser
	var err error
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		resp, err := http.Get(location)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch WSDL: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch WSDL: status %d, body %s", resp.StatusCode, string(body))
		}
		reader = resp.Body
	} else {
		file, err := os.Open(location)
		if err != nil {
			return nil, fmt.Errorf("failed to open WSDL file: %w", err)
		}
		reader = file
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ParseWSDL(data)
}

// ParseWSDL parses WSDL XML bytes into Definitions.
func ParseWSDL(data []byte) (*Definitions, error) {
	var defs Definitions
	if err := xml.Unmarshal(data, &defs); err != nil {
		return nil, err
	}
	return &defs, nil
}
