package system

type MonthlyStatistics struct {
	Revenue                 int64   `json:"Revenue"`
	RevenueGrowth           int64   `json:"RevenueGrowth"`
	RevenueGrowthPercentage float64 `json:"RevenueGrowthPercentage"`
	UniquePurchasers        int64   `json:"UniquePurchasers"`
}
