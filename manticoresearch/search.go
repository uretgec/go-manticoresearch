package manticoresearch

// Responses
// Search
type McSearchResponse struct {
	Took         int                     `json:"took,omitempty"`      // time in milliseconds it took to execute the search
	TimedOut     bool                    `json:"timed_out,omitempty"` // if the query timed out or not
	Aggregations *McAggregationsResponse `json:"aggregations,omitempty"`
	Hits         *McSearchResponseHits   `json:"hits,omitempty"` // search results. has the following properties
	Profile      *interface{}            `json:"profile,omitempty"`
	Warning      interface{}             `json:"warning,omitempty"`
}

type McSearchResponseHits struct {
	MaxScore      int                        `json:"max_score,omitempty"`
	Total         int                        `json:"total,omitempty"` // total number of matching documents
	TotalRelation string                     `json:"total_relation,omitempty"`
	Hits          []McSearchResponseHitsHits `json:"hits,omitempty"` // an array containing matches // docs are here
}

// Stupid Named Strategy
type McSearchResponseHitsHits struct {
	ID_    string                 `json:"_id,omitempty"`     // match id
	Score_ int                    `json:"_score,omitempty"`  // match weight, calculated by ranker
	Source McSearchResponseSource `json:"_source,omitempty"` // an array containing the attributes of this match
}

// Source Struct
type McSearchResponseSource struct {

	// mc models
	Collection int `json:"collection,omitempty" redis:"collection"`
	Bucket     int `json:"bucket,omitempty" redis:"bucket"`

	Shourl  string `json:"shourl,omitempty" redis:"shourl"`
	Title   string `json:"title,omitempty" redis:"title"`
	Slug    string `json:"slug,omitempty" redis:"slug"`
	Url     string `json:"url,omitempty" redis:"url"`
	Content string `json:"content,omitempty" redis:"content"`
	Terms   string `json:"terms,omitempty" redis:"terms"` // all term contents in here (cat,tag,badge)

	CreatedAt int64 `json:"created_at,omitempty" redis:"created_at"`
	UpdatedAt int64 `json:"updated_at,omitempty" redis:"updated_at"`

	Status int `json:"status,omitempty" redis:"status"`

	MakerID   uint64 `json:"maker_id,omitempty" redis:"maker_id"`
	WebsiteID uint64 `json:"website_id,omitempty" redis:"website_id"`
	LicenseID uint64 `json:"license_id,omitempty" redis:"license_id"`

	Maker   string `json:"maker,omitempty" redis:"maker"`
	Website string `json:"website,omitempty" redis:"website"`

	Score    int `json:"score,omitempty" redis:"score"`
	Priority int `json:"priority,omitempty" redis:"priority"`
	Price    int `json:"price,omitempty" redis:"price"`

	Likes     uint64 `json:"likes,omitempty" redis:"likes"`
	Downloads uint64 `json:"downloads,omitempty" redis:"downloads"`
	Views     uint64 `json:"views,omitempty" redis:"views"`

	Cats   []uint64 `json:"cats,omitempty" redis:"cats"`
	Tags   []uint64 `json:"tags,omitempty" redis:"tags"`
	Badges []uint64 `json:"badges,omitempty" redis:"badges"`

	// Videos
	ChannelName string `json:"channel_name,omitempty" redis:"channel_name"`
	ChannelID   string `json:"channel_id,omitempty" redis:"channel_id"`

	// Models
	Cat string `json:"cat,omitempty" redis:"cat"`

	// Medias
	Thumb bool `json:"thumb,omitempty" redis:"thumb"`
}

// Aggregations Struct
type McAggregationsResponse struct {
	// Search: Models
	Badges   *McAggregationsSource `json:"badges"`
	Sources  *McAggregationsSource `json:"sources"`
	Websites *McAggregationsSource `json:"websites"`
}

// Aggregations Sources Struct
type McAggregationsSource struct {
	Buckets []McAggregationsBuckets `json:"buckets"`
}

// Aggregations Bucket Struct
type McAggregationsBuckets struct {
	Key      uint64 `json:"key"`
	DocCount int `json:"doc_count"`
}
