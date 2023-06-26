package manticoresearch

import "time"

// Response Models
type MCDocumentResponse struct {
	// Success
	Index string `json:"_index,omitempty"`
	Id    uint64 `json:"_id,omitempty"`

	Created bool `json:"created,omitempty"`
	Deleted int  `json:"deleted,omitempty"`
	Updated int  `json:"updated,omitempty"`

	Result string `json:"result,omitempty"`
	Status int    `json:"status,omitempty"`
	Found  bool   `json:"found,omitempty"`

	Error MCBulkError `json:"error,omitempty"`
}

// Main Response
type MCDocumentMainResponse []struct {
	Total   int    `json:"total,omitempty"`
	Warning string `json:"warning,omitempty"`
	Error   string `json:"error,omitempty"`

	Columns interface{} `json:"columns,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
type MCDocumentErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// Percolate
type MCPercolateResponse struct {
	Index  string `json:"index"`
	Type   string `json:"type"`
	ID     string `json:"_id"`
	Result string `json:"result"`
}

// Info Response
type McInfoResponse struct {
	ClusterName string          `json:"cluster_name"`
	ClusterUUID string          `json:"cluster_uuid"`
	Name        string          `json:"name"`
	Tagline     string          `json:"tagline"`
	Version     McServerVersion `json:"version"`
}

type McServerVersion struct {
	BuildDate                        time.Time `json:"build_date"`
	BuildFlavor                      string    `json:"build_flavor"`
	BuildHash                        string    `json:"build_hash"`
	BuildSnapshot                    bool      `json:"build_snapshot"`
	BuildType                        string    `json:"build_type"`
	LuceneVersion                    string    `json:"lucene_version"`
	MinimumIndexCompatibilityVersion string    `json:"minimum_index_compatibility_version"`
	MinimumWireCompatibilityVersion  string    `json:"minimum_wire_compatibility_version"`
	Number                           string    `json:"number"`
}
