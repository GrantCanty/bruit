//go:build !old

package web_socket

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/buger/jsonparser"
)

// parseSide parses a single side (bids or asks) within its array boundaries
func parseSide(payload []byte, side []byte, book *FastOrderBook, isBid bool, depth int) {
	sideIdx := bytes.Index(payload, side)
	if sideIdx == -1 {
		return
	}

	endIdx := bytes.IndexByte(payload[sideIdx:], ']')
	if endIdx == -1 {
		return
	}
	endIdx += sideIdx

	cursor := sideIdx + len(side)
	for cursor < endIdx {
		priceIdx := bytes.Index(payload[cursor:endIdx], []byte(`"price"`))
		if priceIdx == -1 {
			break
		}
		cursor += priceIdx + 7

		pStart := bytes.IndexByte(payload[cursor:endIdx], '"') + 1 + cursor
		pEnd := bytes.IndexByte(payload[pStart:endIdx], '"') + pStart
		priceStr := unsafeString(payload[pStart:pEnd])

		cursor = pEnd + 1
		qStart := bytes.Index(payload[cursor:endIdx], []byte(`"qty"`))
		if qStart == -1 {
			break
		}
		cursor += qStart + 5
		qValStart := bytes.IndexByte(payload[cursor:endIdx], '"') + 1 + cursor
		qValEnd := bytes.IndexByte(payload[qValStart:endIdx], '"') + qValStart
		qtyStr := unsafeString(payload[qValStart:qValEnd])

		price, _ := strconv.ParseFloat(priceStr, 64)
		qty, _ := strconv.ParseFloat(qtyStr, 64)
		if isBid {
			book.UpdateBid(price, qty, depth)
		} else {
			book.UpdateAsk(price, qty, depth)
		}

		cursor = qValEnd + 1
	}
}

// parseUpdateUltraFast is a production-safe, zero-allocation, specialized byte scanner.
func parseUpdateUltraFast(payload []byte, book *FastOrderBook, depth int) {
	parseSide(payload, []byte(`"bids"`), book, true, depth)
	parseSide(payload, []byte(`"asks"`), book, false, depth)
}

// parseUpdateSafe parses multiple bids and asks safely with zero heap allocations using jsonparser
func parseUpdateSafe(payload []byte, book *FastOrderBook, depth int) {
	// 1. Process bids
	jsonparser.ArrayEach(payload, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		priceBuf, _, _, _ := jsonparser.Get(value, "price")
		qtyBuf, _, _, _ := jsonparser.Get(value, "qty")
		price, _ := strconv.ParseFloat(unsafeString(priceBuf), 64)
		qty, _ := strconv.ParseFloat(unsafeString(qtyBuf), 64)
		book.UpdateBid(price, qty, depth)
	}, "data", "[0]", "bids")

	// 2. Process asks
	jsonparser.ArrayEach(payload, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		priceBuf, _, _, _ := jsonparser.Get(value, "price")
		qtyBuf, _, _, _ := jsonparser.Get(value, "qty")
		price, _ := strconv.ParseFloat(unsafeString(priceBuf), 64)
		qty, _ := strconv.ParseFloat(unsafeString(qtyBuf), 64)
		book.UpdateAsk(price, qty, depth)
	}, "data", "[0]", "asks")
}

func TestLatencyComparisonOptimizedVsZeroAlloc(t *testing.T) {
	// Discard standard logging
	oldOutput := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(oldOutput)

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

	// Pre-generate 500 mock updates
	client := &WebSocketClient{}
	client.InitBook()
	bookCh := make(chan types.BookRespV2UpdateJSON, 1000)
	logger := settings.LoggingSettings{}

	client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
	<-bookCh

	updatesJSON := make([][]byte, 500)
	for i := 0; i < 500; i++ {
		price := "60000.0"
		qty := fmt.Sprintf("%.8f", 1.5+float64(i)*0.01)
		entry := client.orderBooks["BTC/USD"]
		entry.Book.Bids.Put(types.NumericString(price), types.NumericString(qty))
		cs := calculateChecksum(entry.Book)
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

	// Helper to print nice percentiles
	printPercentiles := func(label string, d []time.Duration) {
		n := len(d)
		minVal := d[0]
		p50 := d[n*50/100]
		p90 := d[n*90/100]
		p95 := d[n*95/100]
		p99 := d[n*99/100]
		p999 := d[n*999/1000]
		p9999 := d[n*9999/10000]
		maxVal := d[n-1]

		fmt.Printf("%-18s | Min: %8v | p50: %8v | p90: %8v | p95: %8v | p99: %8v | p99.9: %8v | p99.99: %8v | Max: %8v\n",
			label, minVal, p50, p90, p95, p99, p999, p9999, maxVal)
	}

	// =========================================================================
	// TEST 1: Latency per batch of 500 updates (2,000 runs)
	// =========================================================================
	fmt.Println("\n=========================================================================================================================")
	fmt.Println("TEST 1: Batch Latency (Time taken to process 500 updates at once, over 2,000 iterations)")
	fmt.Println("=========================================================================================================================")

	runs := 2000
	optimizedBatchDurations := make([]time.Duration, runs)
	zeroAllocUnsafeBatchDurations := make([]time.Duration, runs)
	zeroAllocSafeBatchDurations := make([]time.Duration, runs)
	zeroAllocUltraFastBatchDurations := make([]time.Duration, runs)

	// Run Optimized Tree batch
	for r := 0; r < runs; r++ {
		client = &WebSocketClient{}
		client.InitBook()
		client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
		<-bookCh

		start := time.Now()
		for j := 0; j < 500; j++ {
			client.BookJsonDecoderOptimized(updatesJSON[j], logger, bookCh, 10)
			<-bookCh
		}
		optimizedBatchDurations[r] = time.Since(start)
	}

	// Run Zero Alloc batch (Custom Index Parser)
	for r := 0; r < runs; r++ {
		book := FastOrderBook{
			Symbol: "BTC/USD",
			Bids:   make([]FastLevel, 0, 20),
			Asks:   make([]FastLevel, 0, 20),
		}
		book.Bids = append(book.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
		book.Asks = append(book.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

		start := time.Now()
		for j := 0; j < 500; j++ {
			price, qty := parseMockUpdate(updatesJSON[j])
			book.UpdateBid(price, qty, 10)
		}
		zeroAllocUnsafeBatchDurations[r] = time.Since(start)
	}

	// Run Zero Alloc Safe batch (jsonparser)
	for r := 0; r < runs; r++ {
		book := FastOrderBook{
			Symbol: "BTC/USD",
			Bids:   make([]FastLevel, 0, 20),
			Asks:   make([]FastLevel, 0, 20),
		}
		book.Bids = append(book.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
		book.Asks = append(book.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

		start := time.Now()
		for j := 0; j < 500; j++ {
			parseUpdateSafe(updatesJSON[j], &book, 10)
		}
		zeroAllocSafeBatchDurations[r] = time.Since(start)
	}

	// Run Zero Alloc UltraFast batch (Specialized Scanner)
	for r := 0; r < runs; r++ {
		book := FastOrderBook{
			Symbol: "BTC/USD",
			Bids:   make([]FastLevel, 0, 20),
			Asks:   make([]FastLevel, 0, 20),
		}
		book.Bids = append(book.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
		book.Asks = append(book.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

		start := time.Now()
		for j := 0; j < 500; j++ {
			parseUpdateUltraFast(updatesJSON[j], &book, 10)
		}
		zeroAllocUltraFastBatchDurations[r] = time.Since(start)
	}

	sort.Slice(optimizedBatchDurations, func(i, j int) bool { return optimizedBatchDurations[i] < optimizedBatchDurations[j] })
	sort.Slice(zeroAllocUnsafeBatchDurations, func(i, j int) bool { return zeroAllocUnsafeBatchDurations[i] < zeroAllocUnsafeBatchDurations[j] })
	sort.Slice(zeroAllocSafeBatchDurations, func(i, j int) bool { return zeroAllocSafeBatchDurations[i] < zeroAllocSafeBatchDurations[j] })
	sort.Slice(zeroAllocUltraFastBatchDurations, func(i, j int) bool { return zeroAllocUltraFastBatchDurations[i] < zeroAllocUltraFastBatchDurations[j] })

	printPercentiles("Optimized Tree", optimizedBatchDurations)
	printPercentiles("Zero-Alloc Unsafe", zeroAllocUnsafeBatchDurations)
	printPercentiles("Zero-Alloc Safe", zeroAllocSafeBatchDurations)
	printPercentiles("Zero-Alloc UltraFast", zeroAllocUltraFastBatchDurations)

	// =========================================================================
	// TEST 2: Single-packet Latency (Latency of a single update, over 100,000 updates)
	// =========================================================================
	fmt.Println("\n=========================================================================================================================")
	fmt.Println("TEST 2: Single-Packet Latency (Latency of processing a single websocket update, over 100,000 updates)")
	fmt.Println("=========================================================================================================================")

	totalUpdates := 100000
	optimizedSingleDurations := make([]time.Duration, totalUpdates)
	zeroAllocUnsafeSingleDurations := make([]time.Duration, totalUpdates)
	zeroAllocSafeSingleDurations := make([]time.Duration, totalUpdates)
	zeroAllocUltraFastSingleDurations := make([]time.Duration, totalUpdates)

	// Pre-generate 100,000 updates to keep benchmark pure
	largeUpdatesJSON := make([][]byte, totalUpdates)
	client = &WebSocketClient{}
	client.InitBook()
	client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
	<-bookCh

	for i := 0; i < totalUpdates; i++ {
		price := "60000.0"
		qty := fmt.Sprintf("%.8f", 1.5+float64(i%1000)*0.01)
		entry := client.orderBooks["BTC/USD"]
		entry.Book.Bids.Put(types.NumericString(price), types.NumericString(qty))
		cs := calculateChecksum(entry.Book)
		largeUpdatesJSON[i] = []byte(fmt.Sprintf(`{
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

	// Run Optimized Tree single update test
	client = &WebSocketClient{}
	client.InitBook()
	client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
	<-bookCh
	go func() {
		for range bookCh {}
	}()

	for i := 0; i < totalUpdates; i++ {
		start := time.Now()
		client.BookJsonDecoderOptimized(largeUpdatesJSON[i], logger, bookCh, 10)
		optimizedSingleDurations[i] = time.Since(start)
	}

	// Run Zero Alloc single update test (Custom Index Parser)
	book := FastOrderBook{
		Symbol: "BTC/USD",
		Bids:   make([]FastLevel, 0, 20),
		Asks:   make([]FastLevel, 0, 20),
	}
	book.Bids = append(book.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
	book.Asks = append(book.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

	for i := 0; i < totalUpdates; i++ {
		start := time.Now()
		price, qty := parseMockUpdate(largeUpdatesJSON[i])
		book.UpdateBid(price, qty, 10)
		zeroAllocUnsafeSingleDurations[i] = time.Since(start)
	}

	// Run Zero Alloc Safe single update test (jsonparser)
	bookSafe := FastOrderBook{
		Symbol: "BTC/USD",
		Bids:   make([]FastLevel, 0, 20),
		Asks:   make([]FastLevel, 0, 20),
	}
	bookSafe.Bids = append(bookSafe.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
	bookSafe.Asks = append(bookSafe.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

	for i := 0; i < totalUpdates; i++ {
		start := time.Now()
		parseUpdateSafe(largeUpdatesJSON[i], &bookSafe, 10)
		zeroAllocSafeSingleDurations[i] = time.Since(start)
	}

	// Run Zero Alloc UltraFast single update test (Specialized Scanner)
	bookUltraFast := FastOrderBook{
		Symbol: "BTC/USD",
		Bids:   make([]FastLevel, 0, 20),
		Asks:   make([]FastLevel, 0, 20),
	}
	bookUltraFast.Bids = append(bookUltraFast.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
	bookUltraFast.Asks = append(bookUltraFast.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

	for i := 0; i < totalUpdates; i++ {
		start := time.Now()
		parseUpdateUltraFast(largeUpdatesJSON[i], &bookUltraFast, 10)
		zeroAllocUltraFastSingleDurations[i] = time.Since(start)
	}

	sort.Slice(optimizedSingleDurations, func(i, j int) bool { return optimizedSingleDurations[i] < optimizedSingleDurations[j] })
	sort.Slice(zeroAllocUnsafeSingleDurations, func(i, j int) bool { return zeroAllocUnsafeSingleDurations[i] < zeroAllocUnsafeSingleDurations[j] })
	sort.Slice(zeroAllocSafeSingleDurations, func(i, j int) bool { return zeroAllocSafeSingleDurations[i] < zeroAllocSafeSingleDurations[j] })
	sort.Slice(zeroAllocUltraFastSingleDurations, func(i, j int) bool { return zeroAllocUltraFastSingleDurations[i] < zeroAllocUltraFastSingleDurations[j] })

	printPercentiles("Optimized Tree", optimizedSingleDurations)
	printPercentiles("Zero-Alloc Unsafe", zeroAllocUnsafeSingleDurations)
	printPercentiles("Zero-Alloc Safe", zeroAllocSafeSingleDurations)
	printPercentiles("Zero-Alloc UltraFast", zeroAllocUltraFastSingleDurations)
	fmt.Println("=========================================================================================================================")
	fmt.Println()

	// =========================================================================
	// TEST 3: Latency over 1 Million sequential updates (no batches)
	// =========================================================================
	fmt.Println("\n=========================================================================================================================")
	fmt.Println("TEST 3: Sustained 1 Million Sequential Updates (Latency of single packet processing)")
	fmt.Println("=========================================================================================================================")

	millionUpdates := 1000000
	optimizedMillionDurations := make([]time.Duration, millionUpdates)
	zeroAllocUnsafeMillionDurations := make([]time.Duration, millionUpdates)
	zeroAllocSafeMillionDurations := make([]time.Duration, millionUpdates)
	zeroAllocUltraFastMillionDurations := make([]time.Duration, millionUpdates)

	// Run Optimized Tree 1M updates test
	client = &WebSocketClient{}
	client.InitBook()
	bookCh = make(chan types.BookRespV2UpdateJSON, 1000)
	client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
	<-bookCh

	go func() {
		for range bookCh {}
	}()

	for i := 0; i < millionUpdates; i++ {
		start := time.Now()
		client.BookJsonDecoderOptimized(updatesJSON[i%500], logger, bookCh, 10)
		optimizedMillionDurations[i] = time.Since(start)
	}

	// Run Zero Alloc 1M updates test (Custom Index Parser)
	book = FastOrderBook{
		Symbol: "BTC/USD",
		Bids:   make([]FastLevel, 0, 20),
		Asks:   make([]FastLevel, 0, 20),
	}
	book.Bids = append(book.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
	book.Asks = append(book.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

	for i := 0; i < millionUpdates; i++ {
		start := time.Now()
		price, qty := parseMockUpdate(updatesJSON[i%500])
		book.UpdateBid(price, qty, 10)
		zeroAllocUnsafeMillionDurations[i] = time.Since(start)
	}

	// Run Zero Alloc Safe 1M updates test (jsonparser)
	bookSafeM := FastOrderBook{
		Symbol: "BTC/USD",
		Bids:   make([]FastLevel, 0, 20),
		Asks:   make([]FastLevel, 0, 20),
	}
	bookSafeM.Bids = append(bookSafeM.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
	bookSafeM.Asks = append(bookSafeM.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

	for i := 0; i < millionUpdates; i++ {
		start := time.Now()
		parseUpdateSafe(updatesJSON[i%500], &bookSafeM, 10)
		zeroAllocSafeMillionDurations[i] = time.Since(start)
	}

	// Run Zero Alloc UltraFast 1M updates test (Specialized Scanner)
	bookUltraFastM := FastOrderBook{
		Symbol: "BTC/USD",
		Bids:   make([]FastLevel, 0, 20),
		Asks:   make([]FastLevel, 0, 20),
	}
	bookUltraFastM.Bids = append(bookUltraFastM.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
	bookUltraFastM.Asks = append(bookUltraFastM.Asks, FastLevel{Price: 60001.0, Qty: 2.0})

	for i := 0; i < millionUpdates; i++ {
		start := time.Now()
		parseUpdateUltraFast(updatesJSON[i%500], &bookUltraFastM, 10)
		zeroAllocUltraFastMillionDurations[i] = time.Since(start)
	}

	sort.Slice(optimizedMillionDurations, func(i, j int) bool { return optimizedMillionDurations[i] < optimizedMillionDurations[j] })
	sort.Slice(zeroAllocUnsafeMillionDurations, func(i, j int) bool { return zeroAllocUnsafeMillionDurations[i] < zeroAllocUnsafeMillionDurations[j] })
	sort.Slice(zeroAllocSafeMillionDurations, func(i, j int) bool { return zeroAllocSafeMillionDurations[i] < zeroAllocSafeMillionDurations[j] })
	sort.Slice(zeroAllocUltraFastMillionDurations, func(i, j int) bool { return zeroAllocUltraFastMillionDurations[i] < zeroAllocUltraFastMillionDurations[j] })

	printPercentiles("Optimized Tree", optimizedMillionDurations)
	printPercentiles("Zero-Alloc Unsafe", zeroAllocUnsafeMillionDurations)
	printPercentiles("Zero-Alloc Safe", zeroAllocSafeMillionDurations)
	printPercentiles("Zero-Alloc UltraFast", zeroAllocUltraFastMillionDurations)
	fmt.Println("=========================================================================================================================")
	fmt.Println()
}
