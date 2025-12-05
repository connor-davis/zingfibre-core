package system

type DynamicQueryResult struct {
	Columns []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Label string `json:"label"`
	} `json:"columns"`
	Data []map[string]any `json:"data"`
}
