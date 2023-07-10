package manticoresearch

import (
	"fmt"
	"strings"
)

// Query Options
type McQueryOptions struct {
	Match       *map[string]interface{} `json:"match,omitempty" redis:"match"`
	MatchPhrase *map[string]interface{} `json:"match_phrase,omitempty" redis:"match_phrase"`
	QueryString *string                 `json:"query_string,omitempty" redis:"query_string"`
	MatchAll    *interface{}            `json:"match_all,omitempty" redis:"match_all"`

	// Various-filters
	Equals *map[string]interface{}  `json:"equals,omitempty" redis:"equals"`
	In     *map[string]interface{}  `json:"in,omitempty" redis:"in"`
	Range  *map[string]McQueryRange `json:"range,omitempty" redis:"range"`
}

type McQueryMatchOperator struct {
	Query    string `json:"query,omitempty" redis:"query"`
	Operator string `json:"operator,omitempty" redis:"operator"`
}

type McQueryRange struct {
	Gte int `json:"gte,omitempty" redis:"gte"`
	Gt  int `json:"gt,omitempty" redis:"gt"`
	Lte int `json:"lte,omitempty" redis:"lte"`
	Lt  int `json:"lt,omitempty" redis:"lt"`
}

/*
Query/Match

"match" is a simple query that matches the specified keywords in the specified fields
Example:
```

	"query": {
		"match": {
			"field": "keyword"
			"field1,field2": "keyword"

			// Or you can use _all or * to search all fields.
			"_all": "keyword"

			// You can search all fields except one using "!field":
			"!field1": "keyword"

			// By default keywords are combined using the OR operator. However, you can change that behaviour using the "operator" clause:
			// "operator" can be set to "or" or "and".
			"content,title":{
				"query":"keyword",
				"operator":"or"
			}
		}
	}

```
*/

func NewMcQueryOptions() *McQueryOptions {
	qo := &McQueryOptions{}

	return qo
}

func (qb *McQueryOptions) AddMatch(fields []string, keyword string) *McQueryOptions {
	if qb.Match == nil {
		qb.Match = &map[string]interface{}{}
	}

	(*qb.Match)[strings.Join(fields, ",")] = keyword

	return qb
}
func (qb *McQueryOptions) AddNotMatch(field string, keyword string) *McQueryOptions {
	if qb.Match == nil {
		qb.Match = &map[string]interface{}{}
	}

	(*qb.Match)[fmt.Sprintf("!%s", field)] = keyword

	return qb
}
func (qb *McQueryOptions) AddMatchAllFields(keyword string) *McQueryOptions {
	if qb.Match == nil {
		qb.Match = &map[string]interface{}{}
	}

	(*qb.Match)["_all"] = keyword

	return qb
}
func (qb *McQueryOptions) AddOrMatch(fields []string, keyword string) *McQueryOptions {
	if qb.Match == nil {
		qb.Match = &map[string]interface{}{}
	}

	(*qb.Match)[strings.Join(fields, ",")] = McQueryMatchOperator{
		Query:    keyword,
		Operator: McSearchOperatorOR,
	}

	return qb
}
func (qb *McQueryOptions) AddAndMatch(fields []string, keyword string) *McQueryOptions {
	if qb.Match == nil {
		qb.Match = &map[string]interface{}{}
	}

	(*qb.Match)[strings.Join(fields, ",")] = McQueryMatchOperator{
		Query:    keyword,
		Operator: McSearchOperatorAND,
	}

	return qb
}

/*
Query/MatchPhrase

"match_phrase" is a query that matches the entire phrase. It is similar to a phrase operator in SQL. Here's an example:
Example:
```

	"query": {
		"match": {
			"match_phrase": { "_all" : "had grown quite" }
		}
	}

```
*/
func (qb *McQueryOptions) AddMatchPhrase(fields []string, keyword string) *McQueryOptions {
	if qb.MatchPhrase == nil {
		qb.MatchPhrase = &map[string]interface{}{}
	}

	(*qb.MatchPhrase)[strings.Join(fields, ",")] = keyword

	return qb
}
func (qb *McQueryOptions) AddNotMatchPhrase(field string, keyword string) *McQueryOptions {
	if qb.MatchPhrase == nil {
		qb.MatchPhrase = &map[string]interface{}{}
	}

	(*qb.MatchPhrase)[fmt.Sprintf("!%s", field)] = keyword

	return qb
}
func (qb *McQueryOptions) AddMatchPhraseAllFields(keyword string) *McQueryOptions {
	if qb.MatchPhrase == nil {
		qb.MatchPhrase = &map[string]interface{}{}
	}

	(*qb.MatchPhrase)["_all"] = keyword

	return qb
}

/*
Query/QueryString

"match_phrase" is a query that matches the entire phrase. It is similar to a phrase operator in SQL. Here's an example:
Example:
```

	"query": {
		"match": {
			"query_string": "Church NOTNEAR/3 street"
			"query_string": "@comment_text \"find joe fast \"/2"
		}
	}

```
*/
func (qb *McQueryOptions) AddQueryString(queryString string) *McQueryOptions {
	(*qb.QueryString) = queryString

	return qb
}

/*
Query/Equal

Equality filters are the simplest filters that work with integer, float and string attributes.
Via: https://manual.manticoresearch.com/Searching/Filters#Various-filters

Example:
```

	 "query": {
	    "equals": { "price": 500 }
		"equals": { "any(price)": 100 }
		"in": {
	      "price": [1,10,100]
	    }
		"in": {
	      "all(price)": [1,10]
	    }
		"range": {
	      "price": {
	        "gte": 500,
	        "lte": 1000
	      }
	    }
		"geo_distance": {
	      "location_anchor": {"lat":49, "lon":15},
	      "location_source": {"attr_lat, attr_lon"},
	      "distance_type": "adaptive",
	      "distance":"100 km"
	    }
	  }

```
*/
func (qb *McQueryOptions) AddEquals(k string, v interface{}) *McQueryOptions {
	if qb.Equals == nil {
		qb.Equals = &map[string]interface{}{}
	}

	(*qb.Equals)[k] = v

	return qb
}

// any() which will be positive if the attribute has at least one value which equals to the queried value;
func (qb *McQueryOptions) AddEqualsAny(k string, v interface{}) *McQueryOptions {
	if qb.Equals == nil {
		qb.Equals = &map[string]interface{}{}
	}

	(*qb.Equals)[fmt.Sprintf("any(%s)", k)] = v

	return qb
}

// all() which will be positive if the attribute has a single value and it equals to the queried value
func (qb *McQueryOptions) AddEqualsAll(k string, v interface{}) *McQueryOptions {
	if qb.Equals == nil {
		qb.Equals = &map[string]interface{}{}
	}

	(*qb.Equals)[fmt.Sprintf("all(%s)", k)] = v

	return qb
}

func (qb *McQueryOptions) AddIn(k string, v []int) *McQueryOptions {
	if qb.In == nil {
		qb.In = &map[string]interface{}{}
	}

	(*qb.In)[k] = v

	return qb
}

// any() (equivalent to no function) which will be positive if there's at least one match between the attribute values and the queried values;
func (qb *McQueryOptions) AddInAny(k string, v []int) *McQueryOptions {
	if qb.In == nil {
		qb.In = &map[string]interface{}{}
	}

	(*qb.In)[fmt.Sprintf("any(%s)", k)] = v

	return qb
}

// all() which will be positive if all the attribute values are in the queried set
func (qb *McQueryOptions) AddInAll(k string, v []int) *McQueryOptions {
	if qb.In == nil {
		qb.In = &map[string]interface{}{}
	}

	(*qb.In)[fmt.Sprintf("all(%s)", k)] = v

	return qb
}

/*
Range
gte: greater than or equal to
gt: greater than
lte: less than or equal to
lt: less than
*/
func (qb *McQueryOptions) AddRange(k string, v McQueryRange) *McQueryOptions {
	if qb.Range == nil {
		qb.Range = &map[string]McQueryRange{}
	}

	(*qb.Range)[k] = v

	return qb
}

// TODO: Geo Distance - https://manual.manticoresearch.com/Searching/Filters#Geo-distance-filters
