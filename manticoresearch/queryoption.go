package manticoresearch

// Options
type McOptions struct {
	AccurateAggregation int    `json:"accurate_aggregation,omitempty" redis:"accurate_aggregation"`
	AgentQueryTimeout   int    `json:"agent_query_timeout,omitempty" redis:"agent_query_timeout"`
	AgentRetryCount     int    `json:"agent_retry_count,omitempty" redis:"agent_retry_count"`
	MirrorRetryCount    int    `json:"mirror_retry_count,omitempty" redis:"mirror_retry_count"`
	AgentRetryDelay     int    `json:"agent_retry_delay,omitempty" redis:"agent_retry_delay"`
	ClientTimeout       int    `json:"client_timeout,omitempty" redis:"client_timeout"`
	HostnameLookup      string `json:"hostname_lookup,omitempty" redis:"hostname_lookup"`
	ListenTfo           string `json:"listen_tfo,omitempty" redis:"listen_tfo"`
	PersistentConnectionsLimit           int `json:"persistent_connections_limit,omitempty" redis:"persistent_connections_limit"`
}

/*
POST /search
{

	    "index" : "index_name",
	    "options":
	    {
	        "optionname": "value",
	        "optionname2": <value2>
	    }
	}
*/
func NewMcOptions() McOptions {
	return McOptions{}
}

// Integer. Enables or disables guaranteed aggregate accuracy when running groupby queries in multiple threads. Default is 0.
func (qb McOptions) AddAccurateAggregation(enabled bool) McOptions {
	qb.AccurateAggregation = 0
	if enabled {
		qb.AccurateAggregation = 1
	}

	return qb
}

// Integer. Max time in milliseconds to wait for remote queries to complete, (agent_query_timeout = 10000 # our query can be long, allow up to 10 sec)
func (qb McOptions) AddAgentQueryTimeout(timeout int) McOptions {
	qb.AgentQueryTimeout = timeout

	return qb
}

// The agent_retry_count is an integer that specifies how many times Manticore will attempt to connect and query remote agents
// in a distributed table before reporting a fatal query error. It works similarly to agent_retry_count defined in the "searchd"
// section of the configuration file but applies specifically to the table.
func (qb McOptions) AddAgentRetryCount(retry int) McOptions {
	qb.AgentRetryCount = retry

	return qb
}

// mirror_retry_count serves the same purpose as agent_retry_count. If both values are provided, mirror_retry_count will take precedence, and a warning will be raised.
// For example, if you have 10 mirrors and set agent_retry_count=5, he server will attempt up to 50 retries
// (assuming an average of 5 tries per every 10 mirrors). In case of the option ha_strategy = roundrobin,
// it will actually be exactly 5 tries per mirror.
func (qb McOptions) AddMirrorRetryCount(retry int) McOptions {
	qb.MirrorRetryCount = retry

	return qb
}

// The agent_retry_delay is an integer value that determines the amount of time, in milliseconds,
// that Manticore Search will wait before retrying to query a remote agent in case of a failure.
// This value can be specified either globally in the searchd configuration or on a per-query basis
// using the OPTION retry_delay=XXX clause. If both options are provided, the per-query option will
// take precedence over the global one. The default value is 500 milliseconds (0.5 seconds).
// This option is only relevant if agent_retry_count or the per-query OPTION retry_count are non-zero.
func (qb McOptions) AddAgentRetryDelay(delay int) McOptions {
	qb.AgentRetryDelay = delay

	return qb
}

// The client_timeout option sets the maximum waiting time between requests when using persistent connections. This value is expressed in seconds or with a time suffix. The default value is 5 minutes.
// Example: client_timeout = 1h
func (qb McOptions) AddClientTimeout(timeout int) McOptions {
	qb.ClientTimeout = timeout

	return qb
}

// The hostname_lookup option defines the strategy for renewing hostnames.
// By default, the IP addresses of agent host names are cached at server start to avoid excessive access to DNS.
// However, in some cases, the IP can change dynamically (e.g. cloud hosting) and it may be desirable to not cache the IPs.
// Setting this option to 'request' disables the caching and queries the DNS for each query.
// The IP addresses can also be manually renewed using the FLUSH HOSTNAMES command.
func (qb McOptions) AddHostnameLookup(hostname string) McOptions {
	qb.HostnameLookup = hostname

	return qb
}

// The listen_tfo option allows for the use of the TCP_FASTOPEN flag for all listeners.
// By default, it is managed by the system, but it can be explicitly turned off by setting it to '0'.
// For Linux systems, the server checks the variable /proc/sys/net/ipv4/tcp_fastopen and behaves accordingly.
// Bit 0 manages the client side, while bit 1 rules the listeners. By default, the system has this parameter set to 1, i.e.,
// clients are enabled and listeners are disabled.
func (qb McOptions) AddListenTfo(tfo string) McOptions {
	qb.ListenTfo = tfo

	return qb
}

// persistent_connections_limit = 29 # assume that each host of agents has max_connections = 30 (or 29).
// It is recommended to set this value equal to or less than the max_connections option in the agent's configuration.
func (qb McOptions) AddPersistentConnectionsLimit(limit int) McOptions {
	qb.PersistentConnectionsLimit = limit

	return qb
}
