package itswizard_m_schild

import "encoding/xml"

type Member struct {
	XMLName   xml.Name `xml:"member"`
	Text      string   `xml:",chardata"`
	Sourcedid struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source"`
		ID     string `xml:"id"`
	} `xml:"sourcedid"`
	Idtype string `xml:"idtype"`
	Role   struct {
		Text      string `xml:",chardata"`
		Recstatus string `xml:"recstatus,attr"`
		Status    string `xml:"status"`
	} `xml:"role"`
}
