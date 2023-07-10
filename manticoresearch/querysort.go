package manticoresearch

// Sorting Options
/*
"sort": [ "_score", "id" ]
"sort": [{ "price":"asc" },"id"] // asc|desc
"sort": [{ "gid": { "order":"desc" } }]
"sort": [{ "attr_mva": { "order":"desc", "mode":"max" } }] // max|min
*/
const (
	MCSortOrderASC  = "asc"
	MCSortOrderDESC = "desc"
	// MCSortModeMIN  = "min" // Only use json query
	// MCSortModeMAX  = "max" // Only use json query
)

type McSortOptions struct {
	Sorts []interface{}
}

func NewMcSortOptions() McSortOptions {
	return McSortOptions{
		Sorts: []interface{}{},
	}
}

func (qb McSortOptions) MultiField(fields ...string) McSortOptions {
	qb.Sorts = append(qb.Sorts, fields)

	return qb
}
func (qb McSortOptions) SingleField(field string, order string) McSortOptions {
	qb.Sorts = append(qb.Sorts, map[string]string{field: order})

	return qb
}
func (qb McSortOptions) SingleFieldOrder(field string, order string) McSortOptions {
	qb.Sorts = append(qb.Sorts, map[string]interface{}{
		field: map[string]string{"order": order},
	})

	return qb
}
