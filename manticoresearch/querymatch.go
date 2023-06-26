package manticoresearch

import (
	"fmt"
	"strings"
)

// Query Options
type McQueryOptions struct {
	Match       map[string]interface{} `json:"match,omitempty" redis:"match"`
	MatchPhrase map[string]interface{} `json:"match_phrase,omitempty" redis:"match_phrase"`
	QueryString string                 `json:"query_string,omitempty" redis:"query_string"`
	MatchAll    interface{}            `json:"match_all,omitempty" redis:"match_all"`
}

type McQueryMatchOperator struct {
	Query    string `json:"query,omitempty" redis:"query"`
	Operator string `json:"operator,omitempty" redis:"operator"`
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

func NewMcQueryOptions() McQueryOptions {
	return McQueryOptions{}
}

func (qb McQueryOptions) AddMatch(fields []string, keyword string) McQueryOptions {
	qb.Match[strings.Join(fields, ",")] = keyword

	return qb
}
func (qb McQueryOptions) AddNotMatch(field string, keyword string) McQueryOptions {
	qb.Match[fmt.Sprintf("!%s", field)] = keyword

	return qb
}
func (qb McQueryOptions) AddMatchAllFields(keyword string) McQueryOptions {
	qb.Match["_all"] = keyword

	return qb
}
func (qb McQueryOptions) AddOrMatch(fields []string, keyword string) McQueryOptions {
	qb.Match[strings.Join(fields, ",")] = McQueryMatchOperator{
		Query:    keyword,
		Operator: McSearchOperatorOR,
	}

	return qb
}
func (qb McQueryOptions) AddAndMatch(fields []string, keyword string) McQueryOptions {
	qb.Match[strings.Join(fields, ",")] = McQueryMatchOperator{
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
func (qb McQueryOptions) AddMatchPhrase(fields []string, keyword string) McQueryOptions {
	qb.MatchPhrase[strings.Join(fields, ",")] = keyword

	return qb
}
func (qb McQueryOptions) AddNotMatchPhrase(field string, keyword string) McQueryOptions {
	qb.MatchPhrase[fmt.Sprintf("!%s", field)] = keyword

	return qb
}
func (qb McQueryOptions) AddMatchPhraseAllFields(keyword string) McQueryOptions {
	qb.MatchPhrase["_all"] = keyword

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
func (qb McQueryOptions) AddQueryString(queryString string) McQueryOptions {
	qb.QueryString = queryString

	return qb
}
