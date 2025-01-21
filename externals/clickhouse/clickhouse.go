package clickhouse

type ClickHouse struct {
}

func New(url, apiKey string) *ClickHouse {
	return &ClickHouse{}
}
