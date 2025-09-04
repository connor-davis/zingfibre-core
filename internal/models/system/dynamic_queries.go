package system

type DynamicQueryDatabase string

const (
	ZingDynamicQuery   DynamicQueryDatabase = "zing"
	RadiusDynamicQuery DynamicQueryDatabase = "radius"
)

type DynamicQueryAggregate string

const (
	CountAggregate DynamicQueryAggregate = "COUNT"
	SumAggregate   DynamicQueryAggregate = "SUM"
	AvgAggregate   DynamicQueryAggregate = "AVG"
	MinAggregate   DynamicQueryAggregate = "MIN"
	MaxAggregate   DynamicQueryAggregate = "MAX"
)

type DynamicQueryJoinType string

const (
	OuterJoin DynamicQueryJoinType = "OUTER"
	InnerJoin DynamicQueryJoinType = "INNER"
	LeftJoin  DynamicQueryJoinType = "LEFT"
	RightJoin DynamicQueryJoinType = "RIGHT"
)

type DynamicQueryFilterType string

const (
	StringFilter  DynamicQueryFilterType = "string"
	NumberFilter  DynamicQueryFilterType = "number"
	DateFilter    DynamicQueryFilterType = "date"
	BooleanFilter DynamicQueryFilterType = "boolean"
)

type DynamicQuery struct {
	Database   DynamicQueryDatabase `json:"database"`
	Table      DynamicQueryTable    `json:"table"`
	Columns    []DynamicQueryColumn `json:"columns"`
	Joins      []DynamicQueryJoin   `json:"joins"`
	Filters    []DynamicQueryFilter `json:"filters"`
	Orders     []DynamicQueryOrder  `json:"orders"`
	SubQueries []DynamicQuery       `json:"subQueries"`
}

type DynamicQueryTable struct {
	Name      string `json:"name"`
	IsPrimary bool   `json:"isPrimary"`
}

type DynamicQueryColumn struct {
	Name      string                 `json:"name"`
	Label     string                 `json:"label"`
	Aggregate *DynamicQueryAggregate `json:"aggregate"`
}

type DynamicQueryJoin struct {
	Type              DynamicQueryJoinType `json:"type"`
	LocalDatabase     DynamicQueryDatabase `json:"localDatabase"`
	LocalTable        DynamicQueryTable    `json:"localTable"`
	LocalColumn       string               `json:"localColumn"`
	ReferenceDatabase DynamicQueryDatabase `json:"referenceDatabase"`
	ReferenceTable    DynamicQueryTable    `json:"referenceTable"`
	ReferenceColumn   string               `json:"referenceColumn"`
}

type DynamicQueryFilter struct {
	Column   string                 `json:"field"`
	Value    string                 `json:"value"`
	Operator string                 `json:"operator"`
	Type     DynamicQueryFilterType `json:"type"`
}

type DynamicQueryOrder struct {
	Database   DynamicQueryDatabase `json:"database"`
	Table      DynamicQueryTable    `json:"table"`
	Column     string               `json:"column"`
	Descending bool                 `json:"descending"`
}
