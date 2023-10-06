package itswizard_m_schild

import "encoding/xml"

type DataFromSchild struct {
	XMLName    xml.Name       `xml:"enterprise"`
	Text       string         `xml:",chardata"`
	Xmlns      string         `xml:"xmlns,attr"`
	Properties Properties     `xml:"properties"`
	Person     []schildPerson `xml:"person"`
	Group      []Group        `xml:"group"`
	Membership []Membership   `xml:"membership"`
}

type Properties struct {
	Text       string `xml:",chardata"`
	Datasource string `xml:"datasource"`
	Target     string `xml:"target"`
	Type       string `xml:"type"`
	Datetime   string `xml:"datetime"`
}
