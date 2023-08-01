package manticoresearch

import "encoding/json"

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

func (mc *MCDocumentDeleteRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(mc)
}
func (mc *MCDocumentDeleteRequest) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &mc); err != nil {
		return err
	}

	return nil
}
