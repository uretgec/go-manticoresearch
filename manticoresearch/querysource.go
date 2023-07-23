package manticoresearch

// Source Selection
type McSourceMultiOptions struct {
	Includes []string `json:"includes,omitempty" redis:"includes"`
	Excludes []string `json:"excludes,omitempty" redis:"excludes"`
}

func NewMcSourceMultiOptions() McSourceMultiOptions {
	return McSourceMultiOptions{}
}

func (qb McSourceMultiOptions) AddIncludesAttr(attr ...string) McSourceMultiOptions {
	qb.Includes = append(qb.Includes, attr...)

	return qb
}

func (qb McSourceMultiOptions) AddExcludesAttr(attr ...string) McSourceMultiOptions {
	qb.Excludes = append(qb.Excludes, attr...)

	return qb
}
