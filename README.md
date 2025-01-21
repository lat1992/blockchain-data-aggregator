# blockchain-data-aggregator

A Go-based service that aggregates blockchain transaction data, fetches cryptocurrency prices from CoinGecko, and stores market statistics in ClickHouse.

## Overview

This project processes blockchain transaction data from CSV files, enriches it with historical cryptocurrency prices from CoinGecko API, and stores aggregated market statistics in a ClickHouse database. The system uses a pipeline architecture with concurrent processing for improved performance.

## Features

- CSV file processing with concurrent data reading
- Integration with CoinGecko API for historical cryptocurrency prices
- ClickHouse database for storing market statistics
- Configurable number of concurrent processors
- Docker support for easy deployment
- Comprehensive test coverage

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Make

## Installation

1. Clone the repository:
```bash
git clone https://github.com/lat1992/blockchain-data-aggregator.git
cd blockchain-data-aggregator
```

2. Build and run with Docker:
```bash
make install
```

## Build locally

1. Install and start clickhouse with this document: `https://clickhouse.com/docs/en/getting-started/quick-start`.

Or you can use the docker-compose file: `docker-compose -f docker-compose-clickhouse.yml up -d`

2. Build and run the aggregator locally:
```bash
make build
./build/blockchain-data-aggregator
```

## Configuration

Create a `.env` file from `.env_sample` file:

```env
COINGECKO_URL=https://api.coingecko.com/api/v3
COINGECKO_API_KEY=your-api-key
DATA_PATH=datas
CLICKHOUSE_HOSTNAME=clickhouse:9000
CLICKHOUSE_DATABASE=default
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=
GOROUTINE_NUM=2
```

## Project Structure

```
├── cmd/
│   └── aggregator/         # Main application entry point
├── config/                 # Configuration management
├── externals/              # External service integrations
│   ├── clickhouse/         # ClickHouse database client
│   ├── coingecko/          # CoinGecko API client
│   └── dataGetter/         # CSV data processing
├── internal/               # Internal packages
│   └── services/           # Core business logic
├── mocks/                  # Test mocks
├── scripts/
│   └── clickhouse/         # Database initialization scripts
└── docker-compose.yml      # Docker composition file
└── Dockerfile              # Docker build file
└── Makefile                # Makefile for build and run
```

## Result
```
ch_postgres :) SELECT * FROM market_stats;

SELECT *
FROM market_stats

Query id: 16564976-cb13-4fa5-9e9b-efabaf30c773

   ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
1. │ 2024-04-16 │          0 │               62 │    31.03050136904114 │
2. │ 2024-04-16 │       1609 │                7 │  5528392794798174000 │
3. │ 2024-04-16 │       4974 │               40 │ 21225546994184160000 │
   └────────────┴────────────┴──────────────────┴──────────────────────┘
   ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
4. │ 2024-04-02 │          0 │              104 │ 38908772594862440000 │
5. │ 2024-04-02 │       1609 │                9 │ 21136940686386532000 │
6. │ 2024-04-02 │       4974 │               97 │  3686094245829074000 │
   └────────────┴────────────┴──────────────────┴──────────────────────┘
    ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
 7. │ 2024-04-15 │          0 │              292 │   312.01391314510664 │
 8. │ 2024-04-15 │         43 │                4 │ 16690913441800393000 │
 9. │ 2024-04-15 │       1609 │               20 │ 87114442288911270000 │
10. │ 2024-04-15 │       1660 │                1 │    5.175627913055766 │
11. │ 2024-04-15 │       4974 │              149 │   171.34546636037103 │
    └────────────┴────────────┴──────────────────┴──────────────────────┘
    ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
12. │ 2024-04-01 │          0 │              102 │ 54202152847868090000 │
13. │ 2024-04-01 │       1609 │               13 │ 24059758475721257000 │
14. │ 2024-04-01 │       4974 │              100 │  6179917954455790000 │
    └────────────┴────────────┴──────────────────┴──────────────────────┘

14 rows in set. Elapsed: 0.011 sec.
```

## Future

- Separate the pipeline with two microservices: reader and indexer
