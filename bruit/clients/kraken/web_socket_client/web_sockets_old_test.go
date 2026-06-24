//go:build old

package web_socket

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"strings"
	"testing"
)

// Helper to compute the expected CRC32 checksum of the book state
func calculateChecksum(book *types.BookRespV2UpdateJSON) uint32 {
	crc32q := crc32.MakeTable(crc32.IEEE)
	var priceAsks strings.Builder
	var priceBids strings.Builder

	// Mimics types.verifyLevelTree
	bidsItt := book.Bids.Iterator()
	for bidsItt.Begin(); bidsItt.Next(); {
		price := strings.Replace(string(bidsItt.Key().(types.NumericString)), ".", "", -1)
		qty := strings.Replace(string(bidsItt.Value().(types.NumericString)), ".", "", -1)
		priceBids.WriteString(strings.TrimLeft(price, "0") + strings.TrimLeft(qty, "0"))
	}

	asksItt := book.Asks.Iterator()
	for asksItt.Begin(); asksItt.Next(); {
		price := strings.Replace(string(asksItt.Key().(types.NumericString)), ".", "", -1)
		qty := strings.Replace(string(asksItt.Value().(types.NumericString)), ".", "", -1)
		priceAsks.WriteString(strings.TrimLeft(price, "0") + strings.TrimLeft(qty, "0"))
	}

	priceAsks.WriteString(priceBids.String())
	return crc32.Checksum([]byte(priceAsks.String()), crc32q)
}

func BenchmarkBookProcessing500Old(b *testing.B) {
	// Discard standard logging to avoid terminal writing overhead and spam
	oldOutput := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(oldOutput)

	// 1. Define base snapshot
	snapshotJSON := []byte(`{
		"channel": "book",
		"type": "snapshot",
		"data": [{
			"symbol": "BTC/USD",
			"bids": [{"price": "60000.0", "qty": "1.50000000"}],
			"asks": [{"price": "60001.0", "qty": "2.00000000"}],
			"checksum": 1785233261
		}]
	}`)

	// 2. Pre-generate 500 mock updates with valid checksums
	client := &WebSocketClient{}
	client.InitBook()
	bookCh := make(chan types.BookRespV2UpdateJSON, 1000)
	logger := settings.LoggingSettings{} // console logging disabled

	// Feed snapshot to initialize book
	client.BookJsonDecoder(string(snapshotJSON), logger, bookCh, nil, 10)
	<-bookCh

	updatesJSON := make([][]byte, 500)
	for i := 0; i < 500; i++ {
		// Apply update to client state first to get the correct checksum
		price := "60000.0"
		qty := fmt.Sprintf("%.8f", 1.5+float64(i)*0.01)

		entry := client.orderBooks["BTC/USD"]
		entry.Book.Bids.Put(types.NumericString(price), types.NumericString(qty))
		cs := calculateChecksum(entry.Book)

		// Create the raw JSON representing this update with the correct checksum
		updatesJSON[i] = []byte(fmt.Sprintf(`{
			"channel": "book",
			"type": "update",
			"data": [{
				"symbol": "BTC/USD",
				"bids": [{"price": "%s", "qty": "%s"}],
				"asks": [],
				"checksum": %d
			}]
		}`, price, qty, cs))
	}

	// 3. Run the Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Re-initialize client state for each iteration
		b.StopTimer()
		client = &WebSocketClient{}
		client.InitBook()
		client.BookJsonDecoder(string(snapshotJSON), logger, bookCh, nil, 10)
		<-bookCh
		b.StartTimer()

		// Decode 500 updates
		for j := 0; j < 500; j++ {
			client.BookJsonDecoder(string(updatesJSON[j]), logger, bookCh, nil, 10)
			<-bookCh // Drain channel to simulate system processing
		}
	}
}
