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

Query id: 087ded2a-c70f-4b75-9cf1-cc92d61d6f99

   ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
1. │ 2024-04-15 │          0 │               79 │    55.99393790945316 │
2. │ 2024-04-15 │         43 │                4 │ 16690913441800393000 │
3. │ 2024-04-15 │       1609 │                7 │ 81064699453157850000 │
4. │ 2024-04-15 │       1660 │                1 │    5.175627913055766 │
5. │ 2024-04-15 │       4974 │               68 │   120.80818410577972 │
   └────────────┴────────────┴──────────────────┴──────────────────────┘
   ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
6. │ 2024-04-16 │          0 │               44 │   23.506236448723254 │
7. │ 2024-04-16 │       1609 │                7 │  5528392794798174000 │
8. │ 2024-04-16 │       4974 │               14 │ 21225546994184160000 │
   └────────────┴────────────┴──────────────────┴──────────────────────┘
    ┌───────date─┬─project_id─┬─num_transactions─┬────total_volume_usd─┐
 9. │ 2024-04-02 │          0 │               16 │  26.608654884226453 │
10. │ 2024-04-02 │       1609 │                1 │  542913105974824770 │
11. │ 2024-04-02 │       4974 │               65 │ 3686094245829073400 │
    └────────────┴────────────┴──────────────────┴─────────────────────┘
    ┌───────date─┬─project_id─┬─num_transactions─┬─────total_volume_usd─┐
12. │ 2024-04-01 │          0 │               85 │ 54202152847868090000 │
13. │ 2024-04-01 │       1609 │               11 │ 20459319420788400000 │
14. │ 2024-04-01 │       4974 │                7 │   702035470321170600 │
    └────────────┴────────────┴──────────────────┴──────────────────────┘

14 rows in set. Elapsed: 0.006 sec.
```

## Future

- Separate the pipeline with two microservices: reader and indexer
