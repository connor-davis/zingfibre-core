package system

type ReportCustomer struct {
	FirstName      string `json:"FirstName"`
	Surname        string `json:"Surname"`
	Email          string `json:"Email"`
	PhoneNumber    string `json:"PhoneNumber"`
	RadiusUsername string `json:"RadiusUsername"`
}
