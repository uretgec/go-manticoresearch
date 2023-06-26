package manticoresearch

import "encoding/json"

// Bulk Upsert Document
type MCDocumentBulkUpsertRequest struct {
	Insert  MCDocumentUpsertRequest `json:"insert,omitempty"`
	Replace MCDocumentUpsertRequest `json:"replace,omitempty"`
	Update  MCDocumentUpsertRequest `json:"update,omitempty"`
}

func (mc *MCDocumentBulkUpsertRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(mc)
}

func (mc *MCDocumentBulkUpsertRequest) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &mc); err != nil {
		return err
	}

	return nil
}

// Responses
type MCDocumentBulkResponse struct {
	Items  []MCBulk `json:"items"`
	Errors bool     `json:"errors"`
}

type MCBulk struct {
	Bulk    *MCDocumentResponse `json:"bulk,omitempty"`
	Replace *MCDocumentResponse `json:"replace,omitempty"`
	Update  *MCDocumentResponse `json:"update,omitempty"`
}

type MCBulkError struct {
	Type  string `json:"type,omitempty"`
	Index string `json:"index,omitempty"`
}
