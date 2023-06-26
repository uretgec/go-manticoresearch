package manticoresearch

// Responses
// Search
type McSearchResponse struct {
	Took         int                   `json:"took,omitempty"`      // time in milliseconds it took to execute the search
	TimedOut     bool                  `json:"timed_out,omitempty"` // if the query timed out or not
	Aggregations interface{}           `json:"aggregations,omitempty"`
	Hits         *McSearchResponseHits `json:"hits,omitempty"` // search results. has the following properties
	Profile      *interface{}          `json:"profile,omitempty"`
	Warning      interface{}           `json:"warning,omitempty"`
}

type McSearchResponseHits struct {
	MaxScore      int                        `json:"max_score,omitempty"`
	Total         int                        `json:"total,omitempty"` // total number of matching documents
	TotalRelation string                     `json:"total_relation,omitempty"`
	Hits          []McSearchResponseHitsHits `json:"hits,omitempty"` // an array containing matches // docs are here
}

// Stupid Named Strategy
type McSearchResponseHitsHits struct {
	ID_    uint64      `json:"_id,omitempty"`     // match id
	Score_ int         `json:"_score,omitempty"`  // match weight, calculated by ranker
	Source interface{} `json:"_source,omitempty"` // an array containing the attributes of this match
}
