//go:build !old

package web_socket

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"
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

func BenchmarkBookProcessing500(b *testing.B) {
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
	// We dynamically calculate the checksums so VerifyChecksumUpdate passes successfully.
	client := &WebSocketClient{}
	client.InitBook()
	bookCh := make(chan types.BookRespV2UpdateJSON, 1000)
	logger := settings.LoggingSettings{} // console logging disabled

	// Feed snapshot to initialize book
	client.BookJsonDecoder(snapshotJSON, logger, bookCh, nil, 10)
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
		client.BookJsonDecoder(snapshotJSON, logger, bookCh, nil, 10)
		<-bookCh
		b.StartTimer()

		// Decode 500 updates
		for j := 0; j < 500; j++ {
			client.BookJsonDecoder(updatesJSON[j], logger, bookCh, nil, 10)
			<-bookCh // Drain channel to simulate system processing
		}
	}
}

func BenchmarkBookProcessing500WithCasting(b *testing.B) {
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
	client.BookJsonDecoder(snapshotJSON, logger, bookCh, nil, 10)
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

	// 3. Run the Benchmark with simulated string castings
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		client = &WebSocketClient{}
		client.InitBook()
		
		// Simulate []byte -> string -> []byte conversion
		snapshotStr := string(snapshotJSON)
		snapshotBytes := []byte(snapshotStr)
		client.BookJsonDecoder(snapshotBytes, logger, bookCh, nil, 10)
		<-bookCh
		b.StartTimer()

		// Decode 500 updates
		for j := 0; j < 500; j++ {
			// Simulate []byte -> string -> []byte conversion
			updateStr := string(updatesJSON[j])
			updateBytes := []byte(updateStr)
			client.BookJsonDecoder(updateBytes, logger, bookCh, nil, 10)
			<-bookCh // Drain channel to simulate system processing
		}
	}
}

func (client *WebSocketClient) BookJsonDecoderOptimized(byteResponse []byte, logger settings.LoggingSettings, Bookch chan types.BookRespV2UpdateJSON, bookDepth int) {
	// Optimization 1: Bypass Double Unmarshaling using bytes.Contains
	isBookUpdate := bytes.Contains(byteResponse, []byte(`"channel"`)) &&
		bytes.Contains(byteResponse, []byte(`"book"`)) &&
		bytes.Contains(byteResponse, []byte(`"type"`)) &&
		bytes.Contains(byteResponse, []byte(`"update"`))

	if isBookUpdate {
		// Optimization 2: Use json.Unmarshal directly (avoiding bytes.NewReader + json.NewDecoder allocation)
		var resp types.UpdateBookRespV2WS
		if err := json.Unmarshal(byteResponse, &resp); err == nil {
			symbol := resp.Data[0].Symbol

			client.orderBooksMutex.RLock()
			entry := client.orderBooks[symbol]
			if entry == nil {
				client.orderBooksMutex.RUnlock()
				return
			}
			entry.Mutex.Lock()
			client.orderBooksMutex.RUnlock()
			book := entry.Book
			if !resp.Data[0].Timestamp.IsZero() {
				book.Timestamp = resp.Data[0].Timestamp
			}

			for _, bid := range resp.Data[0].Bids {
				if val, err := bid.Quantity.Float64(); err == nil {
					if val == 0 {
						if _, ok := book.Bids.Get(bid.Price); ok {
							book.Bids.Remove(bid.Price)
						} else {
							client.orderBooksMutex.Lock()
							delete(client.orderBooks, symbol)
							client.orderBooksMutex.Unlock()
							entry.Mutex.Unlock()
							return
						}
					} else {
						book.Bids.Put(bid.Price, bid.Quantity)
					}
				}
			}

			for _, ask := range resp.Data[0].Asks {
				if val, err := ask.Quantity.Float64(); err == nil {
					if val == 0 {
						if _, ok := book.Asks.Get(ask.Price); ok {
							book.Asks.Remove(ask.Price)
						} else {
							client.orderBooksMutex.Lock()
							delete(client.orderBooks, symbol)
							client.orderBooksMutex.Unlock()
							entry.Mutex.Unlock()
							return
						}
					} else {
						book.Asks.Put(ask.Price, ask.Quantity)
					}
				}
			}

			if book.Bids.Size() > bookDepth {
				keys := book.Bids.Keys()
				for i := bookDepth; i < len(keys); i++ {
					book.Bids.Remove(keys[i])
				}
			}

			if book.Asks.Size() > bookDepth {
				keys := book.Asks.Keys()
				for i := bookDepth; i < len(keys); i++ {
					book.Asks.Remove(keys[i])
				}
			}

			if ok := types.VerifyChecksumUpdate(*book, resp); !ok {
				client.orderBooksMutex.Lock()
				delete(client.orderBooks, symbol)
				client.orderBooksMutex.Unlock()
				entry.Mutex.Unlock()
				return
			}

			bookCopy := types.DeepCopyOrderBook(*book)
			entry.Mutex.Unlock()
			Bookch <- bookCopy
			return
		}
	}

	// Fallback for snapshots, subscriptions, etc.
	client.BookJsonDecoder(byteResponse, logger, Bookch, nil, bookDepth)
}

func BenchmarkBookProcessing500Optimized(b *testing.B) {
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

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		client = &WebSocketClient{}
		client.InitBook()
		client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
		<-bookCh
		b.StartTimer()

		for j := 0; j < 500; j++ {
			client.BookJsonDecoderOptimized(updatesJSON[j], logger, bookCh, 10)
			<-bookCh
		}
	}
}

func TestLatencyDistribution(t *testing.T) {
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

	// Helper to run latency tests
	runTest := func(optimized bool) []time.Duration {
		client := &WebSocketClient{}
		client.InitBook()
		bookCh := make(chan types.BookRespV2UpdateJSON, 1000)
		logger := settings.LoggingSettings{}

		// Feed snapshot
		if optimized {
			client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
		} else {
			client.BookJsonDecoder(snapshotJSON, logger, bookCh, nil, 10)
		}
		<-bookCh

		// Pre-generate 500 updates
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

		runs := 2000
		durations := make([]time.Duration, runs)

		for r := 0; r < runs; r++ {
			// Re-init client to keep book size stable
			client = &WebSocketClient{}
			client.InitBook()
			if optimized {
				client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
			} else {
				client.BookJsonDecoder(snapshotJSON, logger, bookCh, nil, 10)
			}
			<-bookCh

			start := time.Now()
			for j := 0; j < 500; j++ {
				if optimized {
					client.BookJsonDecoderOptimized(updatesJSON[j], logger, bookCh, 10)
				} else {
					client.BookJsonDecoder(updatesJSON[j], logger, bookCh, nil, 10)
				}
				<-bookCh
			}
			durations[r] = time.Since(start)
		}
		return durations
	}

	fmt.Println("\nRunning latency test over 2,000 iterations (each iteration processes 500 updates)...")
	originalDurations := runTest(false)
	optimizedDurations := runTest(true)

	// Sort durations to compute percentiles
	sort.Slice(originalDurations, func(i, j int) bool { return originalDurations[i] < originalDurations[j] })
	sort.Slice(optimizedDurations, func(i, j int) bool { return optimizedDurations[i] < optimizedDurations[j] })

	printStats := func(label string, d []time.Duration) {
		n := len(d)
		minVal := d[0]
		p50 := d[n*50/100]
		p90 := d[n*90/100]
		p99 := d[n*99/100]
		maxVal := d[n-1]

		fmt.Printf("%-12s | Min: %8v | Median (p50): %8v | p90: %8v | p99: %8v | Max: %8v\n",
			label, minVal, p50, p90, p99, maxVal)
	}

	fmt.Println("------------------------------------------------------------------------------------------------------")
	printStats("Original", originalDurations)
	printStats("Optimized", optimizedDurations)
	fmt.Println("------------------------------------------------------------------------------------------------------")
}

func TestSustainedLoad100K(t *testing.T) {
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

	runSustainedTest := func(label string, optimized bool) {
		client := &WebSocketClient{}
		client.InitBook()
		bookCh := make(chan types.BookRespV2UpdateJSON, 1000)
		logger := settings.LoggingSettings{}

		// Feed snapshot
		if optimized {
			client.BookJsonDecoderOptimized(snapshotJSON, logger, bookCh, 10)
		} else {
			client.BookJsonDecoder(snapshotJSON, logger, bookCh, nil, 10)
		}
		<-bookCh

		// We pre-calculate all 100,000 updates to keep the measurement pure of loop generation overhead
		updates := make([][]byte, 100000)
		for i := 0; i < 100000; i++ {
			price := "60000.0"
			qty := fmt.Sprintf("%.8f", 1.5+float64(i%1000)*0.01)
			entry := client.orderBooks["BTC/USD"]
			entry.Book.Bids.Put(types.NumericString(price), types.NumericString(qty))
			cs := calculateChecksum(entry.Book)
			updates[i] = []byte(fmt.Sprintf(`{
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

		// Drain channel in a separate goroutine
		go func() {
			for range bookCh {}
		}()

		durations := make([]time.Duration, 100000)

		// Start measurement
		for i := 0; i < 100000; i++ {
			start := time.Now()
			if optimized {
				client.BookJsonDecoderOptimized(updates[i], logger, bookCh, 10)
			} else {
				client.BookJsonDecoder(updates[i], logger, bookCh, nil, 10)
			}
			durations[i] = time.Since(start)
		}

		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

		n := len(durations)
		minVal := durations[0]
		p50 := durations[n*50/100]
		p90 := durations[n*90/100]
		p99 := durations[n*99/100]
		p999 := durations[n*999/1000]
		p9999 := durations[n*9999/10000]
		maxVal := durations[n-1]

		fmt.Printf("%-10s | Min: %6v | p50: %6v | p90: %6v | p99: %6v | p99.9: %6v | p99.99: %6v | Max: %6v\n", 
			label, minVal, p50, p90, p99, p999, p9999, maxVal)
	}

	fmt.Println("\nRunning sustained load test (100,000 updates individual latency percentiles)...")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------")
	runSustainedTest("Original", false)
	runSustainedTest("Optimized", true)
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------")
}

type FastLevel struct {
	Price float64
	Qty   float64
}

type FastOrderBook struct {
	Symbol string
	Bids   []FastLevel
	Asks   []FastLevel
}

func (b *FastOrderBook) UpdateBid(price, qty float64, depth int) {
	idx := sort.Search(len(b.Bids), func(i int) bool {
		return b.Bids[i].Price <= price
	})
	if idx < len(b.Bids) && b.Bids[idx].Price == price {
		if qty == 0 {
			b.Bids = append(b.Bids[:idx], b.Bids[idx+1:]...)
		} else {
			b.Bids[idx].Qty = qty
		}
	} else if qty != 0 {
		b.Bids = append(b.Bids, FastLevel{})
		copy(b.Bids[idx+1:], b.Bids[idx:])
		b.Bids[idx] = FastLevel{Price: price, Qty: qty}
	}
	if len(b.Bids) > depth {
		b.Bids = b.Bids[:depth]
	}
}

func (b *FastOrderBook) UpdateAsk(price, qty float64, depth int) {
	idx := sort.Search(len(b.Asks), func(i int) bool {
		return b.Asks[i].Price >= price
	})
	if idx < len(b.Asks) && b.Asks[idx].Price == price {
		if qty == 0 {
			b.Asks = append(b.Asks[:idx], b.Asks[idx+1:]...)
		} else {
			b.Asks[idx].Qty = qty
		}
	} else if qty != 0 {
		b.Asks = append(b.Asks, FastLevel{})
		copy(b.Asks[idx+1:], b.Asks[idx:])
		b.Asks[idx] = FastLevel{Price: price, Qty: qty}
	}
	if len(b.Asks) > depth {
		b.Asks = b.Asks[:depth]
	}
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func parseMockUpdate(payload []byte) (price, qty float64) {
	pIdx := bytes.Index(payload, []byte(`"price": "`))
	if pIdx == -1 {
		return 0, 0
	}
	pStart := pIdx + len(`"price": "`)
	pEnd := bytes.Index(payload[pStart:], []byte(`"`)) + pStart
	priceStr := unsafeString(payload[pStart:pEnd])

	qIdx := bytes.Index(payload, []byte(`"qty": "`))
	if qIdx == -1 {
		return 0, 0
	}
	qStart := qIdx + len(`"qty": "`)
	qEnd := bytes.Index(payload[qStart:], []byte(`"`)) + qStart
	qtyStr := unsafeString(payload[qStart:qEnd])

	p, _ := strconv.ParseFloat(priceStr, 64)
	q, _ := strconv.ParseFloat(qtyStr, 64)
	return p, q
}

func BenchmarkBookProcessing500ZeroAlloc(b *testing.B) {
	// Pre-generate 500 mock updates as raw byte buffers
	client := &WebSocketClient{}
	client.InitBook()
	bookCh := make(chan types.BookRespV2UpdateJSON, 1000)
	logger := settings.LoggingSettings{}

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

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// Setup the zero-allocation order book structure
		book := FastOrderBook{
			Symbol: "BTC/USD",
			Bids:   make([]FastLevel, 0, 20),
			Asks:   make([]FastLevel, 0, 20),
		}
		// Populate snapshot values
		book.Bids = append(book.Bids, FastLevel{Price: 60000.0, Qty: 1.5})
		book.Asks = append(book.Asks, FastLevel{Price: 60001.0, Qty: 2.0})
		b.StartTimer()

		// Decode and process 500 updates
		for j := 0; j < 500; j++ {
			// Zero-allocation JSON parsing
			price, qty := parseMockUpdate(updatesJSON[j])
			// Zero-allocation sorted map update using binary search on flat slice
			book.UpdateBid(price, qty, 10)
		}
	}
}
