package manticoresearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Manticore Client Options
type MCOption func(*ManticoreClient)

func RegisterMCDefaultHttpClient() MCOption {
	return func(m *ManticoreClient) {
		m.client = &HttpClient{
			name:    DefaultUHCUserAgent,
			timeout: DefaultUHCTimeout,
			reuse:   false,
			debug:   false,
		}
	}
}

func RegisterMCHttpClient(name string, timeout int64, reuse bool, debug bool) MCOption {
	return func(m *ManticoreClient) {
		if name == "" {
			name = DefaultMCName
		}

		if timeout == 0 {
			timeout = DefaultUHCTimeout
		}

		m.client = &HttpClient{
			name:    name,
			timeout: timeout,
			reuse:   reuse,
			debug:   debug,
		}
	}
}

func RegisterMCApiSettings(url string, readOnly bool) MCOption {
	return func(m *ManticoreClient) {
		m.url = url
		m.readOnly = readOnly
	}
}

// Manticore Client Constants
const DefaultMCName = "MyManticoreBot"
const (
	MCApiRouteSql     = "sql"
	MCApiRouteCli     = "cli"
	MCApiRouteBulk    = "bulk"
	MCApiRouteInsert  = "insert"
	MCApiRouteUpdate  = "update"
	MCApiRouteReplace = "replace"
	MCApiRouteDelete  = "delete"
	MCApiRouteSearch  = "search"
)

type ManticoreClient struct {
	// schema://host:port
	url string

	readOnly bool

	client *HttpClient
}

func NewManticoreClient(options ...MCOption) *ManticoreClient {
	a := &ManticoreClient{}

	for _, opt := range options {
		opt(a)
	}

	return a
}

func (m *ManticoreClient) DebugMode(status bool) {
	m.client.debug = status
}

func (m *ManticoreClient) ReadOnlyMode(status bool) {
	m.readOnly = status
}

func (m *ManticoreClient) IsReadOnly() bool {
	return m.readOnly
}

func (m ManticoreClient) generateUrl(args []string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(m.url, "/"), strings.Join(args, "/"))
}

// Info
func (m *ManticoreClient) Info() (resp *McInfoResponse, err error) {
	code, body, errs := m.client.Get(m.generateUrl([]string{}))
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if m.client.debug {
		fmt.Printf("\nBody: %s - Status: %d\n", string(body), code)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		errs = append(errs, err)
	}

	return resp, errs[0]
}

/*
SQL over HTTP
Endpoints /sql and /cli allow running SQL queries via HTTP.

/*
Endpoint: POST /sql

/sql endpoint accepts only SELECT statements and returns the response in HTTP JSON format. The query parameter should be URL-encoded.
The /sql?mode=raw endpoint accepts any SQL query and returns the response in raw format, similar to what you would receive via mysql. The query parameter should also be URL-encoded.
*/
func (m *ManticoreClient) RunSql(payload string) (resp *MCDocumentMainResponse, err error) {
	return nil, errors.New("not implemented yet")
}

/*
Endpoint: POST /cli

The /cli endpoint accepts any SQL query and returns the response in raw format, similar to what you would receive via mysql. Unlike the /sql and /sql?mode=raw endpoints, the query parameter should not be URL-encoded. This endpoint is intended for manual actions using a browser or command line HTTP clients such as curl. It is not recommended to use the /cli endpoint in scripts.
*/
func (m *ManticoreClient) RunCli(payload []byte) (resp *MCDocumentMainResponse, err error) {
	code, body, errs := m.client.Post(m.generateUrl([]string{MCApiRouteCli}), payload)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if m.client.debug {
		fmt.Printf("\nBody: %s - Status: %d\n", string(body), code)
	}

	// catch error json
	singleError := MCDocumentErrorResponse{}
	err = json.Unmarshal(body, &singleError)
	if err == nil {
		return nil, errors.New(singleError.Error)
	}

	// catch main response json
	err = json.Unmarshal(body, &resp)
	if err == nil {
		if m.client.debug {
			fmt.Printf("Resp: %#v\n", resp)
		}

		return resp, nil
	}

	return nil, err
}

// return only unknown interface
func (m *ManticoreClient) RunCliRaw(payload []byte) (resp *interface{}, err error) {
	code, body, errs := m.client.Post(m.generateUrl([]string{MCApiRouteCli}), payload)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if m.client.debug {
		fmt.Printf("\nBody: %s - Status: %d\n", string(body), code)
	}

	// catch main response json
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if m.client.debug {
		fmt.Printf("Resp: %#v\n", resp)
	}

	return resp, nil
}

// Quick query for CLI
func (m *ManticoreClient) ShowThreads() (resp *MCDocumentMainResponse, err error) {
	return m.RunCli([]byte("SHOW THREADS"))
}
func (m *ManticoreClient) ShowTables() (resp *MCDocumentMainResponse, err error) {
	return m.RunCli([]byte("SHOW TABLES"))
}
func (m *ManticoreClient) ShowTableStatus(tableName string) (resp *MCDocumentMainResponse, err error) {
	return m.RunCli([]byte(fmt.Sprintf("SHOW TABLE %s STATUS", tableName)))
}
func (m *ManticoreClient) DescTable(tableName string) (resp *MCDocumentMainResponse, err error) {
	return m.RunCli([]byte(fmt.Sprintf("DESC %s", tableName)))
}
func (m *ManticoreClient) DescPQTable(tableName string) (resp *MCDocumentMainResponse, err error) {
	return m.RunCli([]byte(fmt.Sprintf("DESC %s TABLE", tableName)))
}
func (m *ManticoreClient) DropTable(tableName string) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCli([]byte(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)))
}
func (m *ManticoreClient) TruncateTable(tableName string) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCli([]byte(fmt.Sprintf("TRUNCATE TABLE %s with reconfigure", tableName)))
}

// Queries and kill switch - stupid response return text but content type json?
func (m *ManticoreClient) ShowQueries() (resp *interface{}, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCliRaw([]byte("SHOW QUERIES"))
}

func (m *ManticoreClient) KillQuery(id int) (resp *interface{}, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCliRaw([]byte(fmt.Sprintf("KILL %d", id)))
}

// Flushes all in-memory attribute updates in all the active disk tables to disk. Returns a tag that identifies the result on-disk state (basically, a number of actual disk attribute saves performed since the server startup).
// Look at: attr_flush_period setting (attr_flush_period = 900 # persist updates to disk every 15 minutes)
func (m *ManticoreClient) FlushAttributes(tableName string) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCli([]byte("FLUSH ATTRIBUTES"))
}

// OPTIMIZE merges the RT table's disk chunks down to the number which equals to # of CPU cores * 2 by default. The number of optimized disk chunks can be controlled with option cutoff.
// - server setting optimize_cutoff for overriding the above threshold
// - per-table setting optimize_cutoff
// If OPTION sync=1 is used (0 by default), the command will wait until the optimization process is done (in case the connection interrupts the optimization will continue to run on the server).
func (m *ManticoreClient) OptimizeTable(tableName string, foreground bool) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	sync := 0
	if foreground {
		sync = 1
	}

	return m.RunCli([]byte(fmt.Sprintf("OPTIMIZE TABLE %s OPTION sync=%d", tableName, sync)))
}
func (m *ManticoreClient) OptimizeTableCustom(tableName string, foreground bool, cutoff int) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	sync := 0
	if foreground {
		sync = 1
	}

	return m.RunCli([]byte(fmt.Sprintf("OPTIMIZE TABLE %s OPTION sync=%d,cutoff=%d", tableName, sync, cutoff)))
}

// FREEZE readies a real-time/plain table for a secure backup.
func (m *ManticoreClient) FreezeTable(tableNames ...string) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCli([]byte(fmt.Sprintf("FREEZE %s", strings.Join(tableNames, ","))))
}

// UNFREEZE reactivates previously blocked operations and resumes the internal compaction service. All operations waiting for a table to unfreeze will also be unfrozen and complete normally.
func (m *ManticoreClient) UnfreezeTable(tableNames ...string) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCli([]byte(fmt.Sprintf("UNFREEZE %s", strings.Join(tableNames, ","))))
}

// The SQL statement EXPLAIN QUERY allows displaying the execution tree of a provided full-text query without running an actual search query on the table.
func (m *ManticoreClient) ExplainQuery(tableName string, query string) (resp *interface{}, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	return m.RunCliRaw([]byte(fmt.Sprintf("EXPLAIN QUERY %s '%s'", tableName, query)))
}

func (m *ManticoreClient) ShowStatus(like string) (resp *MCDocumentMainResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	if like == "" {
		// show all
		return m.RunCli([]byte("SHOW STATUS"))
	}

	return m.RunCli([]byte(fmt.Sprintf("SHOW STATUS LIKE '%s%%'", like)))
}

/*
Endpoint: POST /insert JSON
*/
func (m *ManticoreClient) Insert(item MCDocumentUpsertRequest) (resp *MCDocumentResponse, err error) {
	return m.upsert(MCApiRouteInsert, item)
}

/*
Endpoint: POST /bulk "Content-Type: application/x-ndjson" JSON

- you might want to increase max_packet_size value to allow bigger batches
*/
func (m *ManticoreClient) BulkInsert(items ...MCDocumentUpsertRequest) (resp *MCDocumentBulkResponse, err error) {

	payload := []MCDocumentBulkUpsertRequest{}
	for _, item := range items {
		payload = append(payload, MCDocumentBulkUpsertRequest{
			Insert: item,
		})
	}

	return m.bulkUpsert(MCApiRouteInsert, payload...)
}

/*
Endpoint: POST /update JSON

UPDATE changes row-wise attribute values of existing documents in a specified table with new values. Note that you can't update contents of a fulltext field or a columnar attribute. If there's such a need, use REPLACE.
*/
func (m *ManticoreClient) Update(item MCDocumentUpsertRequest) (resp *MCDocumentResponse, err error) {
	return m.upsert(MCApiRouteUpdate, item)
}

/*
Endpoint: POST /bulk "Content-Type: application/x-ndjson" JSON

- you might want to increase max_packet_size value to allow bigger batches
*/
func (m *ManticoreClient) BulkUpdate(items ...MCDocumentUpsertRequest) (resp *MCDocumentBulkResponse, err error) {

	payload := []MCDocumentBulkUpsertRequest{}
	for _, item := range items {
		payload = append(payload, MCDocumentBulkUpsertRequest{
			Update: item,
		})
	}

	return m.bulkUpsert(MCApiRouteUpdate, payload...)
}

/*
Endpoint: POST /replace JSON

REPLACE works similar to INSERT except that if an old document has the same ID as the new document, the old document is marked as deleted before the new document is inserted. Note that the old document does not get physically deleted from the table. The deletion can only happen when chunks are merged in a table, e.g. as a result of an OPTIMIZE.
For HTTP JSON protocol, two request formats are available: Manticore and Elasticsearch-like. You can find both examples in the provided examples.
*/
func (m *ManticoreClient) Replace(item MCDocumentUpsertRequest) (resp *MCDocumentResponse, err error) {
	return m.upsert(MCApiRouteReplace, item)
}

/*
Endpoint: POST /bulk "Content-Type: application/x-ndjson" JSON

- you might want to increase max_packet_size value to allow bigger batches
*/
func (m *ManticoreClient) BulkReplace(items ...MCDocumentUpsertRequest) (resp *MCDocumentBulkResponse, err error) {

	payload := []MCDocumentBulkUpsertRequest{}
	for _, item := range items {
		payload = append(payload, MCDocumentBulkUpsertRequest{
			Update: item,
		})
	}

	return m.bulkUpsert(MCApiRouteReplace, payload...)
}

// Alias insert,replace and delete method
func (m *ManticoreClient) upsert(action string, v MCDocumentUpsertRequest) (resp *MCDocumentResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	// Request
	code, body, errs := m.client.PostJSON(m.generateUrl([]string{action}), v)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if m.client.debug {
		fmt.Printf("\nBody: %s - Status: %d\n", string(body), code)
	}

	// catch main response json
	mainResp := MCDocumentMainResponse{}
	err = json.Unmarshal(body, &mainResp)
	if err == nil {
		if m.client.debug {
			fmt.Printf("MainResp: %#v\n", mainResp)
		}

		// return main response - something wrong? stupid response from manticore server http api
		return nil, errors.New(mainResp[0].Error)
	}

	// catch document response json
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if m.client.debug {
		fmt.Printf("Resp: %#v\n", resp)
	}

	return resp, nil
}

// Alias insert,replace and delete bulk method
func (m *ManticoreClient) bulkUpsert(action string, v ...MCDocumentBulkUpsertRequest) (resp *MCDocumentBulkResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	// Payload (Newline JSON)
	payload := new(bytes.Buffer)
	enc := json.NewEncoder(payload)
	for _, item := range v {
		err := enc.Encode(item)
		if err != nil {
			return nil, err
		}
	}

	// Request
	code, body, errs := m.client.PostNDJSON(m.generateUrl([]string{MCApiRouteBulk}), payload.Bytes())
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if m.client.debug {
		fmt.Printf("\nBody: %s - Status: %d\n", string(body), code)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Errors {
		// something wrong: what a stupid way find errors?
		msg := ""
		if action == MCApiRouteInsert {
			msg = resp.Items[0].Bulk.Error.Type
		} else if action == MCApiRouteUpdate {
			msg = resp.Items[0].Update.Error.Type
		} else if action == MCApiRouteReplace {
			msg = resp.Items[0].Replace.Error.Type
		}

		return nil, errors.New(msg)
	}

	if m.client.debug {
		fmt.Printf("Resp: %#v\n", resp)
	}

	return resp, nil
}

/*
Endpoint: POST /delete JSON

You can delete existing rows (documents) from an existing table based on ID or conditions.
To delete all documents from a table it's recommended to use instead the table truncation as it's a much faster operation.
*/
func (m *ManticoreClient) Delete(v MCDocumentDeleteRequest) (resp *MCDocumentResponse, err error) {
	if m.IsReadOnly() {
		return nil, errors.New("readonly mode active")
	}

	// Request
	code, body, errs := m.client.PostJSON(m.generateUrl([]string{MCApiRouteDelete}), v)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if m.client.debug {
		fmt.Printf("\nBody: %s - Status: %d\n", string(body), code)
	}

	// catch main response json
	mainResp := MCDocumentMainResponse{}
	err = json.Unmarshal(body, &mainResp)
	if err == nil {
		if m.client.debug {
			fmt.Printf("MainResp: %#v\n", mainResp)
		}

		// return main response - something wrong? stupid response from manticore server http api
		return nil, errors.New(mainResp[0].Error)
	}

	// catch document response json
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if m.client.debug {
		fmt.Printf("Resp: %#v\n", resp)
	}

	return resp, nil
}

/*
Endpoint: POST /search JSON

Via: https://manual.manticoresearch.com/Full-text%20matching
- "match" is a simple query that matches the specified keywords in the specified fields
- "match_phrase" is a query that matches the entire phrase. It is similar to a phrase operator in SQL. Here's an example:

- "hello  world" There always is implicit AND operator, so "hello world" means that both "hello" and "world" must be present in matching document.
- "hello | world" OR operator precedence is higher than AND, so looking for cat | dog | mouse means looking for ( cat | dog | mouse ) and not (looking for cat) | dog | mouse.
- "hello MAYBE world" MAYBE operator works much like operator | but doesn't return documents which match only right subtree expression.
- "hello -world" "hello !world" Queries having only negations are not supported by default in Manticore Search. There's the server option not_terms_only_allowed to enable it.
- "@title hello @body world" Field Search
- "@(title,body) hello world" Multiple-field search
- "@!title hello world" Ignore field search
- "@!(title,body) hello world" Ignore multiple-field search
- "@* hello" All field search

Full Text Operators
- "exact * phrase * * for terms" phrase search
- "hello world"~10 Proximity distance search
- "the world is a wonderful place"/3 Quorum matching search for fuzzy matching
- "aaa << bbb << ccc" "(bag of words) << "exact phrase" << red|green|blue" Strict order search
- raining =cats and =dogs ="exact phrase" Exact form keyword modifier search
- "nation* *nation* *national" wildcard search (Requires min_infix_len for prefix (expansion in trail) and/or sufix (expansion in head). If only prefixing is wanted, min_prefix_len can be used instead.)
- - ? can match any(one) character: t?st will match test, but not teast
- - % can match zero or one character : tes% will match tes or test, but not testing
- "^hello world$" Field-start and field-end keyword modifiers
- "boosted^1.234 boostedfieldend$^1.234" boost modifier increases the word IDF score
- "hello NEAR/3 world NEAR/4 "my test"" near
- "Church NOTNEAR/3 street" not near
- "all SENTENCE words SENTENCE "in one sentence"" "Bill Gates" PARAGRAPH "Steve Jobs" SENTENCE and PARAGRAPH operators matches the document when both its arguments are within the same sentence or the same paragraph of text, respectively
- ZONE:(h3,h4) ZONE limit operator is quite similar to field limit operator
- ZONESPAN:(h2) ZONESPAN limit operator is similar to the ZONE operator, but requires the match to occur in a single contiguous span. In the example above, ZONESPAN:th hello world would not match the document, since "hello" and "world" do not occur within the same span.

Escape Characters: !    "    $    '    (    )    -    /    <    @    \    ^    |    ~

All full-text match clauses can be combined with must, must_not and should operators of an HTTP bool query.
*/
func (m *ManticoreClient) Search(builder McSearchQueryBuilder, profile bool) error {
	return nil
}

/*
Endpoint: PUT /json/pq/{pq_table_name}/doc/{?id}?refresh=1 JSON
*/
func (m *ManticoreClient) InsertPq(v interface{}, profile bool) error {
	return errors.New("not implemented yet")
}
func (m *ManticoreClient) ReplacePq(v interface{}, profile bool) error {
	return errors.New("not implemented yet")
}
func (m *ManticoreClient) UpsertPq(v interface{}, profile bool) error {
	return errors.New("not implemented yet")
}
func (m *ManticoreClient) DeletePq(v MCDocumentDeleteRequest) (resp *MCDocumentResponse, err error) {
	return m.Delete(v)
}

/*
Endpoint: POST /json/pq/{pq_table_name}/search JSON
*/
func (m *ManticoreClient) SearchPq(builder McPercolateQueryBuilder, profile bool) error {
	return errors.New("not implemented yet")
}

// TODO: Autocomplete
// TODO: Suggestion

/*
Backup

You can also back up your data through SQL by running the simple command BACKUP TO /path/to/backup.
-> BACKUP

	[{TABLE | TABLES} a[, b]]
	[{OPTION | OPTIONS}
	  async = {on | off | 1 | 0 | true | false | yes | no}
	  [, compress = {on | off | 1 | 0 | true | false | yes | no}]
	]
	TO path_to_backup

For instance, to back up tables a and b to the /backup directory, run the following command:

-> BACKUP TABLES a, b TO /backup

There are options available to control and adjust the backup process, such as:

async: makes the backup non-blocking, allowing you to receive a response with the query ID immediately and run other queries while the backup is ongoing. The default value is 0.
compress: enables file compression using zstd. The default value is 0. For example, to run a backup of all tables in async mode with compression enabled to the /tmp directory:

-> BACKUP OPTION async = yes, compress = yes TO /tmp
*/
func (m *ManticoreClient) Backup(opt MCBackupRequest) error {
	cmd := []string{"BACKUP"}

	if len(opt.Tables) == 1 {
		cmd = append(cmd, "TABLE")
	} else if len(opt.Tables) > 1 {
		cmd = append(cmd, "TABLES", strings.Join(opt.Tables, ","))
	}

	// Options
	cmd = append(cmd, "OPTIONS", fmt.Sprintf("async=%t, compress=%t", opt.Options.Async, opt.Options.Compress))

	// to path
	if opt.Path == "" {
		opt.Path = "/tmp"
	}
	cmd = append(cmd, "TO", opt.Path)

	resp, err := m.RunCliRaw([]byte(strings.Join(cmd, " ")))
	if err != nil {
		return err
	}

	if m.client.debug {
		fmt.Printf("Resp: %#v\n", resp)
	}

	return err
}

/*
Restore - alias import table

If you decide to migrate from Plain mode to RT mode and in some other cases, real-time and percolate tables built in the Plain mode can be imported to Manticore running in the RT mode using the IMPORT TABLE statement. The general syntax is as follows:

-> IMPORT TABLE table_name FROM 'path'
*/
func (m *ManticoreClient) Restore(tableName, path string) error {
	resp, err := m.RunCliRaw([]byte(fmt.Sprintf("IMPORT TABLE %s FROM '%s'", tableName, path)))
	if err != nil {
		return err
	}

	if m.client.debug {
		fmt.Printf("Resp: %#v\n", resp)
	}

	return err
}

// TODO: reload table
// TODO: rotate table

// TODO: https://manual.manticoresearch.com/Updating_table_schema_and_settings
// TODO: https://manual.manticoresearch.com/Data_creation_and_modification/Adding_data_from_external_storages/Adding_data_to_tables/Attaching_a_plain_table_to_RT_table#Attaching-table---general-syntax

// TODO:
type McPercolateQueryBuilder struct{}
