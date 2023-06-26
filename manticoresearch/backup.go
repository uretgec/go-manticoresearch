package manticoresearch

// Backup Document
type MCBackupRequest struct {
	Tables  []string       `json:"tables"`
	Options MCBackupOption `json:"options"`
	Path    string         `json:"path"`
}

type MCBackupOption struct {
	Async    bool `json:"async"`
	Compress bool `json:"compress"`
}
