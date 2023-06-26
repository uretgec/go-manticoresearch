package manticoresearch

// Source Selection
type McSourceMultiOptions struct {
	Includes []string `json:"includes,omitempty" redis:"includes"`
	Excludes []string `json:"excludes,omitempty" redis:"excludes"`
}
