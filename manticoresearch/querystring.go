package manticoresearch

import (
	"fmt"
	"strings"
)

/*
QueryString Builder

The query string can contain certain operators that allow telling the conditions of how the words from the query string should be matched.
AND Operator: There always is implicit AND operator, so "hello world" means that both "hello" and "world" must be present in matching document.
OR Operator: OR operator precedence is higher than AND, so looking for cat | dog | mouse means looking for ( cat | dog | mouse ) and not (looking for cat) | dog | mouse.
MAYBE Operator: MAYBE operator works much like operator | but doesn't return documents which match only right subtree expression.
NEGATION Operator: The negation operator enforces a rule for a word to not exist.
- Queries having only negations are not supported by default in Manticore Search.
- There's the server option not_terms_only_allowed to enable it.

Field Operators:
@title hello @body world
@@relaxed @nosuchfield my query -> Field limit operator limits subsequent searching to a given field. Normally, query will fail with an error message if given field name does not exist in the searched table. However, that can be suppressed by specifying @@relaxed option at the very beginning of the query
@body[50] hello -> Field position limit additionally restricts the searching to first N position within given field (or fields). For example, @body [50] hello will not match the documents where the keyword hello occurs at position 51 and below in the body.
@(title,body) hello world
@!title hello world -> ignore filed
@!(title,body) hello world -> ignore multiple fields

Phrase Operator:
"exact * phrase * * for terms" -> The phrase operator requires the words to be next to each other.

Proximity Operator:
"hello world"~10 -> Proximity distance is specified in words, adjusted for word count, and applies to all words within quotes. For instance, "cat dog mouse"~5 query means that there must be less than 8-word span which contains all 3 words, ie. CAT aaa bbb ccc DOG eee fff MOUSE document will not match this query, because this span is exactly 8 words long.

Quorum Operator:
"the world is a wonderful place"/3 -> Quorum matching operator introduces a kind of fuzzy matching. It will only match those documents that pass a given threshold of given words. The example above ("the world is a wonderful place"/3) will match all documents that have at least 3 of the 6 specified words. Operator is limited to 255 keywords
"the world is a wonderful place"/0.5 -> Instead of an absolute number, you can also specify a number between 0.0 and 1.0 (standing for 0% and 100%), and Manticore will match only documents with at least the specified percentage of given words. The same example above could also have been written "the world is a wonderful place"/0.5 and it would match documents with at least 50% of the 6 words.

Exact Operator: Exact form keyword modifier will match the document only if the keyword occurred in exactly the specified form.
raining =cats and =dogs
="exact phrase"

* Exact form operator requires index_exact_words option to be enabled.

Wildcard Operator:
nation* *nation* *national -> Requires min_infix_len for prefix (expansion in trail) and/or sufix (expansion in head). If only prefixing is wanted, min_prefix_len can be used instead.

In addition, the following inline wildcard operators are supported:

? can match any(one) character: t?st will match test, but not teast
% can match zero or one character : tes% will match tes or test, but not testing
The inline operators require dict=keywords and infixing enabled.

Example:
```
"hello world" @title "example program"~5 @body python -(php|perl) @* code
```

The full meaning of this search is:

Find the words 'hello' and 'world' adjacently in any field in a document;
Additionally, the same document must also contain the words 'example' and 'program' in the title field, with up to, but not including, 5 words between the words in question; (E.g. "example PHP program" would be matched however "example script to introduce outside data into the correct context for your program" would not because two terms have 5 or more words between them)
Additionally, the same document must contain the word 'python' in the body field, but not contain either 'php' or 'perl';
Additionally, the same document must contain the word 'code' in any field.
OR operator precedence is higher than AND, so "looking for cat | dog | mouse" means "looking for ( cat | dog | mouse )" and not "(looking for cat) | dog | mouse".

To understand how a query will be executed, Manticore Search offer query profile tooling for viewing the query tree created by a query expression.
*/
type McQueryStringBuilder struct {
	QueryString string
}

func (qb McQueryStringBuilder) And(keywords ...string) string {
	return strings.Join(keywords, "   ")
}
func (qb McQueryStringBuilder) Or(keywords ...string) string {
	return strings.Join(keywords, " | ")
}
func (qb McQueryStringBuilder) Maybe(keywords ...string) string {
	return strings.Join(keywords, " MAYBE ")
}
func (qb McQueryStringBuilder) Negation(keywords ...string) string {
	return strings.Join(keywords, " !")
}
func (qb McQueryStringBuilder) NegationOr(keywords ...string) string {
	return fmt.Sprintf("-(%s)", strings.Join(keywords, "|"))
}
func (qb McQueryStringBuilder) Field(field string, keyword string) string {
	return fmt.Sprintf("@%s %s", field, keyword)
}
func (qb McQueryStringBuilder) FieldNot(field string, keyword string) string {
	return fmt.Sprintf("@!%s %s", field, keyword)
}
func (qb McQueryStringBuilder) Fields(fields []string, keyword string) string {
	return fmt.Sprintf("@(%s) %s", strings.Join(fields, ","), keyword)
}
func (qb McQueryStringBuilder) FieldsNot(fields []string, keyword string) string {
	return fmt.Sprintf("@!(%s) %s", strings.Join(fields, ","), keyword)
}
func (qb McQueryStringBuilder) FieldNoError(field string, keyword string) string {
	return fmt.Sprintf("@@%s %s", field, keyword)
}
func (qb McQueryStringBuilder) FieldAll(keyword string) string {
	return fmt.Sprintf("@* %s", keyword)
}
func (qb McQueryStringBuilder) Phrase(keywords []string) string {
	return fmt.Sprintf("\"%s\"", strings.Join(keywords, " "))
}
func (qb McQueryStringBuilder) Proximity(keywords []string, proximity int) string {
	return fmt.Sprintf("\"%s\"~%d", strings.Join(keywords, " "), proximity)
}
func (qb McQueryStringBuilder) Quorum(keywords []string, quorum int) string {
	return fmt.Sprintf("\"%s\"/%d", strings.Join(keywords, " "), quorum)
}
func (qb McQueryStringBuilder) QuorumFloat(keywords []string, quorum float32) string {
	return fmt.Sprintf("\"%s\"/%.1f", strings.Join(keywords, " "), quorum)
}
func (qb McQueryStringBuilder) Exact(keyword string) string {
	return fmt.Sprintf("=%s", keyword)
}
