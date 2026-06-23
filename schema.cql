CREATE KEYSPACE IF NOT EXISTS stock_keyspace
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

CREATE TABLE IF NOT EXISTS stock_keyspace.candles (
	symbol text,
	datetime timestamp,
	open double,
	high double,
	low double,
	close double,
	volume bigint,
	PRIMARY KEY ((symbol), datetime)
) WITH CLUSTERING ORDER BY (datetime ASC)
