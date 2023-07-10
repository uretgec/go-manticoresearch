package manticoresearch

import "encoding/json"

// TODO: keisnlikle gerekli
// Faced search te desteklemeli
// TODO: validator eklenerek veriler kontrol edilmeli.
/*
Search is a key functionality of Manticore Search. You can:

- Do full-text search and search results highlighting
- Do non-full-text filtering
- Use expressions for filtering
- Use various search options
- Use multi-queries and sub-selects
- Do aggreations and faceting of search results
- And many more
*/

const (
	McSearchOperatorOR  = "or"
	McSearchOperatorAND = "and"
)

type McSearchQueryBuilder struct {
	// Profile: get full-text query tree structure
	Profile bool `json:"profile,omitempty" redis:"profile"`

	// index name or table name
	Index   string     `json:"index" redis:"index"`
	Options *McOptions `json:"options,omitempty" redis:"options"`

	// Query Options
	Query *McQueryOptions `json:"query,omitempty" redis:"query"`

	// Source
	// Each entry can be an attribute name or a wildcard (*, % and ? symbols are supported)
	// "_source":"attr*",
	// "_source": [ "attr1", "attri*" ]"
	// "_source": {"includes": [ "attr1", "attri*" ],"excludes": [ "*desc*" ]}
	Source interface{} `json:"_source,omitempty" redis:"_source"`

	// Highlight
	// Highlight interface{} `json:"highlight,omitempty" redis:"highlight"`

	// Sorting
	Sort interface{} `json:"sort,omitempty" redis:"sort"`
	// Sorting With: You can enable weight calculation by setting the track_scores property to true
	TrackScores bool `json:"track_scores,omitempty" redis:"track_scores"`

	// Group: Aggregation and Facet
	// "aggs": {"<aggr_name>": {"terms": {"field": "<attribute>","size": <int value>}}
	Aggs map[string]McAggregationTerms `json:"aggs,omitempty" redis:"aggs"`

	// Facets can aggregate over expressions.
	// "expressions": {"price_range": "INTERVAL(price,200,400,600,800)"}
	// price_range is the agg field name for facet search
	Expressions map[string]string `json:"expressions,omitempty" redis:"expressions"`

	// Limitation Options
	Size       int64 `json:"size,omitempty" redis:"size"`
	From       int64 `json:"from,omitempty" redis:"from"`
	Offset     int64 `json:"offset,omitempty" redis:"offset"`
	Limit      int64 `json:"limit,omitempty" redis:"limit"`
	MaxMatches int64 `json:"max_matches,omitempty" redis:"max_matches"` // same as limit. Look at: max_matches setting (Default: 1000)
}

func NewMCSearchQueryBuilder(indexName string) *McSearchQueryBuilder {
	return &McSearchQueryBuilder{
		Index: indexName,
	}
}

func (mc *McSearchQueryBuilder) MarshalBinary() ([]byte, error) {
	return json.Marshal(mc)
}
func (mc *McSearchQueryBuilder) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &mc); err != nil {
		return err
	}

	return nil
}

// Profile: get full-text query tree structure
func (qb *McSearchQueryBuilder) SetProfile(profile bool) *McSearchQueryBuilder {
	qb.Profile = profile

	return qb
}

// index name or table name
func (qb *McSearchQueryBuilder) SetIndex(indexName string) *McSearchQueryBuilder {
	qb.Index = indexName

	return qb
}

// SELECT * FROM test WHERE MATCH('@title hello @body world') OPTION ranker=bm25, max_matches=3000,
func (qb *McSearchQueryBuilder) SetOptions(options *McOptions) *McSearchQueryBuilder {
	qb.Options = options

	return qb
}
func (qb *McSearchQueryBuilder) SetQuery(query *McQueryOptions) *McSearchQueryBuilder {
	qb.Query = query

	return qb
}

// Sorting
func (qb *McSearchQueryBuilder) SetSort(opt McSortOptions) *McSearchQueryBuilder {
	qb.Sort = opt.Sorts

	return qb
}

func (qb *McSearchQueryBuilder) SetLimit(limit int64) *McSearchQueryBuilder {
	qb.Limit = limit

	return qb
}
func (qb *McSearchQueryBuilder) SetOffset(offset int64) *McSearchQueryBuilder {
	qb.Offset = offset

	return qb
}
func (qb *McSearchQueryBuilder) SetSize(size int64) *McSearchQueryBuilder {
	qb.Size = size

	return qb
}
func (qb *McSearchQueryBuilder) SetFrom(from int64) *McSearchQueryBuilder {
	qb.From = from

	return qb
}
func (qb *McSearchQueryBuilder) SetMaxMatches(limit int64) *McSearchQueryBuilder {
	qb.MaxMatches = limit

	return qb
}

// Aggregation Options - Facet Query Options
type McAggregationTerms struct {
	Field string `json:"field"` // value must contain the name of the attribute or expression we are faceting
	// By default each facet result set is limited to 20 values.
	Size int `json:"size,omitempty"` // optional size specifies the maximum number of buckets to include into the result. When not specified it inherits the main query's limit.
}

func (qb *McSearchQueryBuilder) AddAgg(aggName string, field string, size int) *McSearchQueryBuilder {
	agg := McAggregationTerms{
		Field: field,
	}

	if size > 0 {
		agg.Size = size
	}

	qb.Aggs[aggName] = agg

	return qb
}

func (qb *McSearchQueryBuilder) AddExp(field string, exp string) *McSearchQueryBuilder {
	qb.Expressions[field] = exp

	return qb
}

// Query Match Options
func (qb *McSearchQueryBuilder) AddQueryMatchOptions(opt *McQueryOptions) *McSearchQueryBuilder {
	qb.Query = opt

	return qb
}
