package itswizard_m_schild

import "encoding/xml"

type Membership struct {
	XMLName   xml.Name `xml:"membership"`
	Text      string   `xml:",chardata"`
	Sourcedid struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source"`
		ID     string `xml:"id"`
	} `xml:"sourcedid"`
	Member []Member `xml:"member"`
}
