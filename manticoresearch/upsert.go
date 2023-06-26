package manticoresearch

import "encoding/json"

// Upsert Document
type MCDocumentUpsertRequest struct {
	// Name of the index
	Index string `json:"index"`

	// cluster name
	Cluster string `json:"cluster,omitempty"`

	// Document ID
	Id uint64 `json:"id,omitempty"`

	// Object with document data
	Doc interface{} `json:"doc"`
}

func (mc *MCDocumentUpsertRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(mc)
}

func (mc *MCDocumentUpsertRequest) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &mc); err != nil {
		return err
	}

	return nil
}
