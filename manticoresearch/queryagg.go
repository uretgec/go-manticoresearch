package manticoresearch

// Aggregation Options - Facet Query Options
type McAggregation struct {
	Terms McAggregationTerms `json:"terms"`
}

// Aggregation Terms
type McAggregationTerms struct {
	Field string `json:"field"` // value must contain the name of the attribute or expression we are faceting
	// By default each facet result set is limited to 20 values.
	Size int `json:"size,omitempty"` // optional size specifies the maximum number of buckets to include into the result. When not specified it inherits the main query's limit.
}

func (qb *McSearchQueryBuilder) NewAgg() *McSearchQueryBuilder {
	qb.Aggs = map[string]McAggregation{}

	return qb
}

func (qb *McSearchQueryBuilder) AddAgg(groupName string, field string, size int) *McSearchQueryBuilder {
	// check agg is set
	if len(qb.Aggs) == 0 {
		qb.Aggs = map[string]McAggregation{}
	}

	agg := McAggregationTerms{
		Field: field,
	}

	if size > 0 {
		agg.Size = size
	}

	// add new aggs terms
	qb.Aggs[groupName] = McAggregation{
		Terms: agg,
	}

	return qb
}

func (qb *McSearchQueryBuilder) AddExp(field string, exp string) *McSearchQueryBuilder {
	qb.Expressions[field] = exp

	return qb
}
