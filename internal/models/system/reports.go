package system

type ReportCustomer struct {
	FirstName      string `json:"FirstName,omitempty"`
	Surname        string `json:"Surname,omitempty"`
	Email          string `json:"Email,omitempty"`
	PhoneNumber    string `json:"PhoneNumber,omitempty"`
	RadiusUsername string `json:"RadiusUsername,omitempty"`
}

type ReportRechargeTypeCount struct {
	RechargeName    string `json:"RechargeName,omitempty"`
	RechargeCount   int    `json:"RechargeCount,omitempty"`
	RechargePeriod  string `json:"RechargePeriod,omitempty"`
	RechargeMaxDate string `json:"RechargeMaxDate,omitempty"`
}

type ReportExpiringCustomer struct {
	FirstName            string `json:"FirstName,omitempty"`
	Surname              string `json:"Surname,omitempty"`
	Email                string `json:"Email,omitempty"`
	PhoneNumber          string `json:"PhoneNumber,omitempty"`
	RadiusUsername       string `json:"RadiusUsername,omitempty"`
	LastPurchaseDuration string `json:"LastPurchaseDuration,omitempty"`
	LastPurchaseSpeed    string `json:"LastPurchaseSpeed,omitempty"`
	Expiration           string `json:"Expiration,omitempty"`
	Address              string `json:"Address,omitempty"`
}

type ReportRecharge struct {
	DateCreated string  `json:"DateCreated,omitempty"`
	Email       string  `json:"Email,omitempty"`
	FirstName   string  `json:"FirstName,omitempty"`
	Surname     string  `json:"Surname,omitempty"`
	ItemName    string  `json:"ItemName,omitempty"`
	Amount      float64 `json:"Amount,omitempty"`
	Successful  bool    `json:"Successful,omitempty"`
	ServiceId   int64   `json:"ServiceId,omitempty"`
	BuildName   string  `json:"BuildName,omitempty"`
	BuildType   string  `json:"BuildType,omitempty"`
}

type ReportSummary struct {
	DateCreated    string  `json:"DateCreated,omitempty"`
	ItemName       string  `json:"ItemName,omitempty"`
	RadiusUsername string  `json:"RadiusUsername,omitempty"`
	AmountGross    string  `json:"AmountGross,omitempty"`
	AmountFee      string  `json:"AmountFee,omitempty"`
	AmountNet      string  `json:"AmountNet,omitempty"`
	CashCode       string  `json:"CashCode,omitempty"`
	CashAmount     float64 `json:"CashAmount,omitempty"`
	ServiceId      int64   `json:"ServiceId,omitempty"`
	BuildName      string  `json:"BuildName,omitempty"`
	BuildType      string  `json:"BuildType,omitempty"`
}
