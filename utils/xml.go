package utils

import (
	"encoding/xml"
)

type XmlInterface map[string]interface{}

// XmlInterface marshals into XML.
func (s XmlInterface) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "xml",
	}
	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(ToString(value)), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}
