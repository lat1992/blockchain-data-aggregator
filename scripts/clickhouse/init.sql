CREATE TABLE IF NOT EXISTS market_stats (
    date Date,
    project_id UInt64,
    num_transactions UInt64,
    total_volume_usd Float64,
    INDEX project_id_index (project_id) TYPE
    SET
        (100) GRANULARITY 4,
) ENGINE = MergeTree ()
PARTITION BY
    date
ORDER BY
    (project_id, date) SETTINGS index_granularity = 8192;
