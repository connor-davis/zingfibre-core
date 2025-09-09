package system

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
	Database   string               `json:"database"`
	IsSubQuery bool                 `json:"isSubQuery"`
	Table      DynamicQueryTable    `json:"table"`
	Columns    []DynamicQueryColumn `json:"columns"`
	Joins      []DynamicQueryJoin   `json:"joins"`
	Filters    []DynamicQueryFilter `json:"filters"`
	Orders     []DynamicQueryOrder  `json:"orders"`
	SubQueries []DynamicQuery       `json:"subQueries"`
}

type DynamicQueryTable struct {
	Table     string `json:"table"`
	IsPrimary bool   `json:"isPrimary"`
}

type DynamicQueryColumn struct {
	Table     DynamicQueryTable      `json:"table"`
	Column    string                 `json:"column"`
	Label     string                 `json:"label"`
	Aggregate *DynamicQueryAggregate `json:"aggregate"`
}

type DynamicQueryJoin struct {
	Type            DynamicQueryJoinType `json:"type"`
	LocalTable      DynamicQueryTable    `json:"localTable"`
	LocalColumn     string               `json:"localColumn"`
	ReferenceTable  DynamicQueryTable    `json:"referenceTable"`
	ReferenceColumn string               `json:"referenceColumn"`
	SubQuery        *DynamicQuery        `json:"subQuery"`
}

type DynamicQueryFilter struct {
	Table    DynamicQueryTable      `json:"table"`
	Column   string                 `json:"field"`
	Value    string                 `json:"value"`
	Operator string                 `json:"operator"`
	Type     DynamicQueryFilterType `json:"type"`
}

type DynamicQueryOrder struct {
	Table      DynamicQueryTable `json:"table"`
	Column     string            `json:"column"`
	Descending bool              `json:"descending"`
}
