package system

type ReportCustomer struct {
	FullName       string `json:"FullName,omitempty"`
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
	FullName             string `json:"FullName,omitempty"`
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
	FullName    string  `json:"FullName,omitempty"`
	ItemName    string  `json:"ItemName,omitempty"`
	Amount      float64 `json:"Amount,omitempty"`
	Method      string  `json:"Method,omitempty"`
	Successful  bool    `json:"Successful,omitempty"`
	ServiceId   int64   `json:"ServiceId,omitempty"`
	BuildName   string  `json:"BuildName,omitempty"`
	BuildType   string  `json:"BuildType,omitempty"`
}

type ReportSummary struct {
	DateCreated    string `json:"DateCreated,omitempty"`
	ItemName       string `json:"ItemName,omitempty"`
	RadiusUsername string `json:"RadiusUsername,omitempty"`
	Method         string `json:"Method,omitempty"`
	AmountGross    string `json:"AmountGross,omitempty"`
	AmountFee      string `json:"AmountFee,omitempty"`
	AmountNet      string `json:"AmountNet,omitempty"`
	CashCode       string `json:"CashCode,omitempty"`
	CashAmount     string `json:"CashAmount,omitempty"`
	ServiceId      int64  `json:"ServiceId,omitempty"`
	BuildName      string `json:"BuildName,omitempty"`
	BuildType      string `json:"BuildType,omitempty"`
}
