package itswizard_m_schild

import "encoding/xml"

type Group struct {
	XMLName   xml.Name `xml:"group"`
	Text      string   `xml:",chardata"`
	Recstatus string   `xml:"recstatus,attr"`
	Sourcedid struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source"`
		ID     string `xml:"id"`
	} `xml:"sourcedid"`
	Grouptype struct {
		Text      string `xml:",chardata"`
		Scheme    string `xml:"scheme"`
		Typevalue struct {
			Text  string `xml:",chardata"`
			Level string `xml:"level,attr"`
		} `xml:"typevalue"`
	} `xml:"grouptype"`
	Description struct {
		Text  string `xml:",chardata"`
		Short string `xml:"short"`
		Long  string `xml:"long"`
	} `xml:"description"`
	Extension struct {
		Text                    string `xml:",chardata"`
		XSchildnrwOwnerType     string `xml:"x-schildnrw-owner-type"`
		XSchildnrwOwnerDistrict string `xml:"x-schildnrw-owner-district"`
	} `xml:"extension"`
}
