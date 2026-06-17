CREATE KEYSPACE IF NOT EXISTS stock_keyspace
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

USE stock_keyspace;

CREATE TABLE IF NOT EXISTS candles (
    symbol text,
    datetime timestamp,
    high double,
    low double,
    close double,
    volume bigint,
    PRIMARY KEY ((symbol), datetime)
) WITH CLUSTERING ORDER BY (datetime ASC);
