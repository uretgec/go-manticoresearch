package manticoresearch

// Delete Document
type MCDocumentDeleteRequest struct {
	// Index name
	Index string `json:"index"`
	// cluster name
	Cluster string `json:"cluster,omitempty"`
	// Document ID
	Id uint64 `json:"id,omitempty"`
	// Query tree object
	Query *interface{} `json:"query,omitempty"`
}