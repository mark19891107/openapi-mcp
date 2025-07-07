package wsdl

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Definitions struct {
	XMLName         xml.Name   `xml:"definitions"`
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
	Name   string   `xml:"name,attr"`
	Input  ParamRef `xml:"input"`
	Output ParamRef `xml:"output"`
}

type ParamRef struct {
	Message string `xml:"message,attr"`
}

type Binding struct {
	Name       string             `xml:"name,attr"`
	Type       string             `xml:"type,attr"`
	Operations []BindingOperation `xml:"operation"`
}

type BindingOperation struct {
	Name   string        `xml:"name,attr"`
	Action SOAPOperation `xml:"operation"`
}

type SOAPOperation struct {
	SOAPAction string `xml:"soapAction,attr"`
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

func Load(location string) (*Definitions, error) {
	var reader io.ReadCloser
	var err error
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		resp, err := http.Get(location)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch %s: %s", location, string(body))
		}
		reader = resp.Body
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
	return Parse(data)
}

func Parse(data []byte) (*Definitions, error) {
	var defs Definitions
	if err := xml.Unmarshal(data, &defs); err != nil {
		return nil, err
	}
	return &defs, nil
}
