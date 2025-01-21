package clickhouse

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/lat1992/blockchain-data-aggregator/internal"
)

type ClickHouse struct {
	conn driver.Conn
}

func New(host, database, user, password string) (*ClickHouse, error) {
	conn, err := connect(host, database, user, password)
	if err != nil {
		return nil, fmt.Errorf("error connecting to clickhouse: %w", err)
	}
	return &ClickHouse{
		conn: conn,
	}, nil
}

func connect(host, database, user, password string) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{host},
			Auth: clickhouse.Auth{
				Database: database,
				Username: user,
				Password: password,
			},
			DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "tcp", addr)
			},
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
			DialTimeout:          time.Second * 30,
			MaxOpenConns:         5,
			MaxIdleConns:         5,
			ConnMaxLifetime:      time.Duration(10) * time.Minute,
			ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
			BlockBufferSize:      10,
			MaxCompressionBuffer: 10240,
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "aggregator-client", Version: "0.1"},
				},
			},
		})
	)

	if err != nil {
		return nil, fmt.Errorf("error connecting to clickhouse: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			slog.Error("Exception", "code", exception.Code, "msg", exception.Message, "trace", exception.StackTrace)
		}
		return nil, fmt.Errorf("error pinging clickhouse: %w", err)
	}
	return conn, nil
}

func (c *ClickHouse) InsertMarket(stats []internal.MarketStat) error {
	ctx := context.Background()
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO market_stats")
	if err != nil {
		return err
	}
	for _, stat := range stats {
		err := batch.Append(stat.Date, stat.ProjectId, stat.NumTx, stat.TotalVolume)
		if err != nil {
			return fmt.Errorf("error appending to batch: %w", err)
		}
	}
	return batch.Send()
}
