package itswizard_m_schild

import "encoding/xml"

/*
Diese Person funtkioniert
*/
type schildPerson struct {
	XMLName   xml.Name `xml:"person"`
	Text      string   `xml:",chardata"`
	Recstatus string   `xml:"recstatus,attr"`
	Sourcedid struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source"`
		ID     string `xml:"id"`
	} `xml:"sourcedid"`
	Name struct {
		Text string `xml:",chardata"`
		Fn   string `xml:"fn"`
		N    struct {
			Text   string `xml:",chardata"`
			Family string `xml:"family"`
			Given  string `xml:"given"`
		} `xml:"n"`
	} `xml:"name"`
	Demographics struct {
		Text string `xml:",chardata"`
		Bday string `xml:"bday"`
	} `xml:"demographics"`
	Email      string `xml:"email"`
	Systemrole struct {
		Text           string `xml:",chardata"`
		Systemroletype string `xml:"systemroletype,attr"`
	} `xml:"systemrole"`
	Institutionrole struct {
		Text                string `xml:",chardata"`
		Institutionroletype string `xml:"institutionroletype,attr"`
	} `xml:"institutionrole"`
	Extension struct {
		Text                  string `xml:",chardata"`
		XSchildnrwPersonState string `xml:"x-schildnrw-person-state"`
		XSchildnrwGrade       string `xml:"x-schildnrw-grade"`
	} `xml:"extension"`
}
