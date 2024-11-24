-- +goose Up
-- +goose StatementBegin
CREATE TABLE erc20_transfers (
	-- same as block timestamp, use for timescale partitions
	timestamp TIMESTAMPTZ NOT NULL,
	-- from it's Ethereum Address, e.g. 0xA
	sender TEXT NOT NULL,
	-- to it's Ethereum Address, e.g. 0xA
	recipient TEXT NOT NULL,
	-- trasfer amount after decimals
	amount NUMERIC NOT NULL,
	-- transaction hash, basically 32 bytes
	tx_hash TEXT NOT NULL,
	-- token contract address, where transfer event was emitted
	token_address TEXT NOT NULL
);

SELECT create_hypertable('erc20_transfers', by_range('timestamp'));

CREATE MATERIALIZED VIEW erc20_transfers_sum_5_min WITH ( timescaledb.continuous ) AS
SELECT
	time_bucket ('5 minutes', timestamp) AS time, token_address, sum(amount) as total_amount
FROM
	erc20_transfers
GROUP BY
	time, token_address
WITH NO DATA;

-- create timscale refresh policies
SELECT add_continuous_aggregate_policy('erc20_transfers_sum_5_min',
	start_offset => INTERVAL '5 minutes',
	end_offset => NULL,
	schedule_interval => INTERVAL '1 minute');
	
-- enable real-time aggregation
ALTER MATERIALIZED VIEW erc20_transfers_sum_5_min SET (timescaledb.materialized_only = false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
