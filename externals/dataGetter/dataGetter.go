package dataGetter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/lat1992/blockchain-data-aggregator/externals"
)

type DataGetter struct {
	recordChannel chan externals.Record
	endChannel    chan bool
}

func NewDataGetter() *DataGetter {
	return &DataGetter{
		recordChannel: make(chan externals.Record),
		endChannel:    make(chan bool),
	}
}

func (g *DataGetter) ReadDataFromFiles(path string) error {
	p, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open path: %w", err)
	}
	files, err := p.Readdir(0)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	for _, f := range files {
		g.readDataFromFile(path + "/" + f.Name())
	}
	g.endChannel <- true
	return nil
}

func (g *DataGetter) readDataFromFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close file", "err", err)
		}
	}()

	parser := csv.NewReader(file)

	header, err := g.getHeader(parser)
	if err != nil {
		return fmt.Errorf("failed to get header: %w", err)
	}

	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}
		g.recordChannel <- externals.Record{
			Timestamp: record[header["ts"]],
			Event:     record[header["event"]],
			ProjectId: record[header["project_id"]],
			Props:     record[header["props"]],
			Nums:      record[header["nums"]],
		}
	}
	return nil
}

func (g *DataGetter) getHeader(parser *csv.Reader) (map[string]int, error) {
	record, err := parser.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	header := make(map[string]int)
	for i, r := range record {
		if r == "ts" || r == "event" || r == "project_id" || r == "props" || r == "nums" {
			header[r] = i
		}
	}
	return header, nil
}

func (g *DataGetter) Channel() chan externals.Record {
	return g.recordChannel
}

func (g *DataGetter) EndChannel() chan bool {
	return g.endChannel
}
