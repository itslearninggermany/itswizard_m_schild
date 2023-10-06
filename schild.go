package itswizard_m_schild

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/alexmullins/zip"
	itswizard_basic "github.com/itslearninggermany/itswizard_m_basic"
	normalisation "github.com/itslearninggermany/itswizard_m_normalisation"
	itswizard_s3Bucket "github.com/itslearninggermany/itswizard_m_s3bucket"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"strings"
)

type Schild struct {
	gorm.Model
	FileDirContent     string `gorm:"mediumtext"`
	FileDirPersons     string `gorm:"mediumtext"`
	FileDirGroups      string `gorm:"mediumtext"`
	FileDirMemberships string `gorm:"mediumtext"`
	OrganisationID     uint
	InstitutionID      uint
	s3bucket           *itswizard_s3Bucket.Bucket `gorm:"-"`
}

/*
Creates a newSChild struct for all Operations.
*/
func NewSchild(filename, password string, organisationID, institutionID uint, bucketname string, bucketRegion string) (out *Schild, err error) {
	out = new(Schild)
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}
		r, err := f.Open()
		if err != nil {
			return nil, err
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		buf = normalisation.Normalise(buf)

		out.s3bucket, err = itswizard_s3Bucket.NewBucket(bucketname, bucketRegion)
		if err != nil {
			return nil, err
		}

		filedir := fmt.Sprint(uuid.Must(uuid.NewV4(), err).String(), ".xml")

		/*
			TODO: Add Normalisation
		*/

		out.FileDirContent = filedir
		err = out.s3bucket.ContentUpload(filedir, buf)
		if err != nil {
			return nil, err
		}
	}

	out.OrganisationID = organisationID
	out.InstitutionID = institutionID
	return out, nil
}

func (p *Schild) PrintOutAsXML() (out string, err error) {
	buf, err := p.s3bucket.DownloadContent(p.FileDirContent)
	if err != nil {
		return "", err
	}
	return string(buf), err
}

/*
Makes all Persons, Groups and Membership ready to use
*/
func (p *Schild) Initial() error {
	out, err := p.s3bucket.DownloadContent(p.FileDirContent)
	if err != nil {
		return err
	}

	var schildData DataFromSchild
	err = xml.Unmarshal(out, &schildData)
	if err != nil {
		return err
	}

	var persons []itswizard_basic.DbPerson15
	var groups []itswizard_basic.DbGroup15
	var memberships []itswizard_basic.DbGroupMembership15

	//1. Persons
	for i := 0; i < len(schildData.Person); i++ {
		var firstname string
		if schildData.Person[i].Name.N.Given != "" {
			firstname = schildData.Person[i].Name.N.Given
		} else {
			firstname = "nn"
		}

		var lastname string
		if schildData.Person[i].Name.N.Family != "" {
			lastname = schildData.Person[i].Name.N.Family
		} else {
			lastname = "nn"
		}

		name := fmt.Sprint(firstname, " ", lastname)

		ausgabe := strings.Split(name, " ")

		var username string
		if len(ausgabe) > 1 {
			username = fmt.Sprint(ausgabe[0], ".", ausgabe[len(ausgabe)-1])
		}

		//profile
		var profile string
		switch schildData.Person[i].Institutionrole.Institutionroletype {
		case "Student":
			profile = "Student"
		case "Faculty":
			profile = "Staff"
		default:
			profile = "Guest"

		}

		username = ""

		persons = append(persons, itswizard_basic.DbPerson15{
			ID:                 fmt.Sprint(p.InstitutionID, "++", p.OrganisationID, "++", schildData.Person[i].Sourcedid.ID),
			SyncPersonKey:      schildData.Person[i].Sourcedid.ID,
			FirstName:          firstname,
			LastName:           lastname,
			Username:           username,
			Profile:            profile,
			Email:              schildData.Person[i].Email,
			DbOrganisation15ID: p.OrganisationID,
			DbInstitution15ID:  p.InstitutionID,
		})

	}

	m := make(map[string]string)

	//2. Groups
	for i := 0; i < len(schildData.Group); i++ {
		// Level 1 is the "Schulträger" and Level 2 is the school
		if schildData.Group[i].Grouptype.Typevalue.Level == "1" {
			continue
		}
		if schildData.Group[i].Grouptype.Typevalue.Level == "2" {
			continue
		}

		exist := false
		if strings.Contains(schildData.Group[i].Description.Long, "Klasse") {
			exist = true
		}
		if strings.Contains(schildData.Group[i].Description.Long, "Alle") {
			exist = true
		}

		if exist {
			var groupname string
			groupname = schildData.Group[i].Description.Short

			/*
				if schildData.Group[i].Description.Long == "" {
					groupname = schildData.Group[i].Description.Short
				}
				if schildData.Group[i].Description.Long == "" && schildData.Group[i].Description.Short == "" {
					groupname = uuid.Must(uuid.NewV4()).String()
				}
			*/

			groups = append(groups, itswizard_basic.DbGroup15{
				ID:                 fmt.Sprint(p.InstitutionID, "++", p.OrganisationID, "++", schildData.Group[i].Sourcedid.ID),
				SyncID:             schildData.Group[i].Sourcedid.ID + "haha",
				Name:               groupname,
				ParentGroupID:      "rootPointer",
				Level:              1,
				IsCourse:           false,
				DbInstitution15ID:  p.InstitutionID,
				DbOrganisation15ID: p.OrganisationID,
			})
			m[schildData.Group[i].Sourcedid.ID] = groupname
		}
	}

	//3. Membership
	/*
		Alg. durch alle Gruppen gehen und die Mitglieder aufnehmen
	*/
	for groupid, groupname := range m {
		for i := 0; i < len(schildData.Membership); i++ {
			if schildData.Membership[i].Sourcedid.ID == groupid {
				for is := 0; is < len(schildData.Membership[i].Member); is++ {
					memberships = append(memberships, itswizard_basic.DbGroupMembership15{
						ID:                 fmt.Sprint(p.InstitutionID, "++", p.OrganisationID, "++", fmt.Sprint(p.InstitutionID, "++", p.OrganisationID, "++", groupname), "++", schildData.Membership[i].Member[is].Sourcedid.ID),
						PersonSyncKey:      schildData.Membership[i].Member[is].Sourcedid.ID,
						GroupName:          fmt.Sprint(p.InstitutionID, "++", p.OrganisationID, "++", groupname),
						DbInstitution15ID:  p.InstitutionID,
						DbOrganisation15ID: p.OrganisationID,
					})
				}
			}
		}
	}

	bufPersons, err := json.Marshal(persons)
	if err != nil {
		return err
	}
	err = nil
	fileDir := fmt.Sprint(uuid.Must(uuid.NewV4(), err).String(), ".json")
	p.FileDirPersons = fileDir
	err = p.s3bucket.ContentUpload(fileDir, bufPersons)
	if err != nil {
		return err
	}

	bufGroups, err := json.Marshal(groups)
	if err != nil {
		return err
	}

	err = nil
	fileDir = fmt.Sprint(uuid.Must(uuid.NewV4(), err).String(), ".json")
	p.FileDirGroups = fileDir
	err = p.s3bucket.ContentUpload(fileDir, bufGroups)
	if err != nil {
		return err
	}

	err = nil
	fileDir = fmt.Sprint(uuid.Must(uuid.NewV4(), err).String(), ".json")
	bufMemberships, err := json.Marshal(memberships)
	if err != nil {
		return err
	}
	p.FileDirMemberships = fileDir
	err = p.s3bucket.ContentUpload(fileDir, bufMemberships)
	if err != nil {
		return err
	}
	return nil
}

func (p *Schild) GetPersons() (out []itswizard_basic.DbPerson15, err error) {
	ctx, err := p.s3bucket.DownloadContent(p.FileDirPersons)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal([]byte(ctx), &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (p *Schild) GetGroups() (out []itswizard_basic.DbGroup15, err error) {
	ctx, err := p.s3bucket.DownloadContent(p.FileDirGroups)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal([]byte(ctx), &out)
	if err != nil {
		return out, err
	}
	return out, nil

}

func (p *Schild) GetMemberships() (out []itswizard_basic.DbGroupMembership15, err error) {
	ctx, err := p.s3bucket.DownloadContent(p.FileDirGroups)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal([]byte(ctx), &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (p *Schild) GetPersonsAsString() (out []byte, err error) {
	return p.s3bucket.DownloadContent(p.FileDirPersons)
}

func (p *Schild) GetGroupsAsString() (out []byte, err error) {
	return p.s3bucket.DownloadContent(p.FileDirGroups)
}

func (p *Schild) GetMembershipsAsString() (out []byte, err error) {
	return p.s3bucket.DownloadContent(p.FileDirMemberships)
}
