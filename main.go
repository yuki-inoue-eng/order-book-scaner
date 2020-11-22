package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yuki-inoue-eng/order-book-searcher/lib"
	"github.com/yuki-inoue-eng/order-book-searcher/lib/oanda"
)

var (
	fileNamePrefix       = flag.String("fname", "ob-search", "")
	periodStr            = flag.String("period", "none", "specify the aggregation period.")
	instrumentStr        = flag.String("instrument", "USD_JPY", "specify a instrument.")
	stopOrderStr         = flag.String("stop-order", "none", "")
	limitOrderStr        = flag.String("limit-order", "none", "")
	losingPositionStr    = flag.String("losingPosition", "none", "")
	profitingPositionStr = flag.String("profiting-position", "none", "")
	//netAmount            = flag.Bool("net-amount", false, "") 純額は後ほど
)

func main() {
	flag.Parse()

	// validate period
	if len(*periodStr) == 0 {
		log.Fatal("period is required")
		return
	}
	period := strings.Split(*periodStr, "-")
	if len(period) != 2 {
		log.Fatalf("invalid period: %s (ex: 2020/10/01-2020/11/01)", *periodStr)
		return
	}
	layout := "2006/01/02"
	since, err := time.Parse(layout, period[0])
	if err != nil {
		log.Fatalf("failed to parse period to date time: %v", period[0])
		return
	}
	until, err := time.Parse(layout, period[1])
	if err != nil {
		log.Fatalf("failed to parse period to date time: %v", period[1])
		return
	}

	// validate instrument
	if len(*instrumentStr) == 0 {
		log.Fatal("instrument is required")
		return
	}
	instrument := oanda.ToInstrument(*instrumentStr)
	if instrument == oanda.InstrumentUNKNOWN {
		log.Fatalf("invalid instrument: %s", *instrumentStr)
		return
	}

	// validate stop-order
	var stopOrderLowerLimits []float64
	if len(*stopOrderStr) == 0 {
		for _, s := range strings.Split(*stopOrderStr, "-") {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatalf("invalid stop-order: %s (ex: 0.8-1.0): %v", *stopOrderStr, err)
				return
			}
			stopOrderLowerLimits = append(stopOrderLowerLimits, v)
		}
	}

	// validate limit-order
	var limitOrderLowerLimits []float64
	if len(*limitOrderStr) == 0 {
		for _, l := range strings.Split(*limitOrderStr, "-") {
			v, err := strconv.ParseFloat(l, 64)
			if err != nil {
				log.Fatalf("invalid limit-order: %s (ex: 0.8-1.0): %v", *limitOrderStr, err)
				return
			}
			limitOrderLowerLimits = append(limitOrderLowerLimits, v)
		}
	}

	// validate losing-position
	var losingPositionLowerLimits []float64
	if len(*losingPositionStr) == 0 {
		for _, l := range strings.Split(*losingPositionStr, "-") {
			v, err := strconv.ParseFloat(l, 64)
			if err != nil {
				log.Fatalf("invalid losing-position: %s (ex: 0.8-1.0): %v", *losingPositionStr, err)
				return
			}
			losingPositionLowerLimits = append(losingPositionLowerLimits, v)
		}
	}

	// validate profiting-position
	var profitingPositionLowerLimits []float64
	if len(*profitingPositionStr) == 0 {
		for _, p := range strings.Split(*profitingPositionStr, "-") {
			v, err := strconv.ParseFloat(p, 64)
			if err != nil {
				log.Fatalf("invalid profiting-position: %s (ex: 0.8-1.0): %v", *profitingPositionStr, err)
				return
			}
			profitingPositionLowerLimits = append(profitingPositionLowerLimits, v)
		}
	}

	if len(stopOrderLowerLimits) == 0 && len(limitOrderLowerLimits) == 0 &&
		len(profitingPositionLowerLimits) == 0 && len(losingPositionLowerLimits) == 0 {
		log.Fatalf("at least one of stop-order, limit-order, profiting-position, losing-position is required")
		return
	}

	var allRecords []record
	const twentyMinutes = 1200
	for iTime := since.Unix(); iTime < until.Unix(); iTime += twentyMinutes {
		t := time.Unix(iTime, 0)
		client := oanda.NewClient(lib.ParamTradeForceOandaAPIKey.FetchValue(), "Practice")
		orderBook, err := client.FetchOrderBook(instrument, &t)
		if err != nil {
			log.Printf("failed to fetch order book (at %s): %v", t.String(), err)
			continue
		}
		client = oanda.NewClient(lib.ParamTradeForceOandaAPIKey.FetchValue(), "Practice")
		positionBook, err := client.FetchPositionBook(instrument, &t)
		if err != nil {
			log.Printf("failed to fetch position book (at %s): %v ", t.String(), err)
			continue
		}
		price := orderBook.Price
		dateTime := orderBook.Time
		instrument := orderBook.Instrument
		const targetRange = 20
		oShort, oLong, err := orderBook.ExtractBucketVicinityOfPrice(price, targetRange)
		if err != nil {
			log.Printf("failed to extract order book buckets: %v", err)
			continue
		}
		pShort, pLong, err := positionBook.ExtractBucketVicinityOfPrice(price, targetRange)
		if err != nil {
			log.Printf("failed to extract position book buckets: %v", err)
			continue
		}

		// search applicable stop loss order
		var stopOrderRecords []record
		if len(stopOrderLowerLimits) > 0 {
			var shortBuckets []bucket
			var longBuckets []bucket
			for i := 0; i < targetRange-len(stopOrderLowerLimits); i++ {
				for j := 0; j < len(stopOrderLowerLimits); j++ {
					if oShort[i+j].ShortCountPercent >= stopOrderLowerLimits[j] {
						b := bucket{
							priceRange:    oShort[i+j].Price,
							shortOrder:    oShort[i+j].ShortCountPercent,
							longOrder:     oShort[i+j].LongCountPercent,
							shortPosition: pShort[i+j].ShortCountPercent,
							longPosition:  pShort[i+j].LongCountPercent,
						}
						shortBuckets = append(shortBuckets, b)
					}
					if oLong[i+j].LongCountPercent >= stopOrderLowerLimits[j] {
						b := bucket{
							priceRange:    oLong[i+j].Price,
							shortOrder:    oLong[i+j].ShortCountPercent,
							longOrder:     oLong[i+j].LongCountPercent,
							shortPosition: pLong[i+j].ShortCountPercent,
							longPosition:  pLong[i+j].LongCountPercent,
						}
						longBuckets = append(longBuckets, b)
					}
				}
			}
			if len(shortBuckets) != 0 {
				stopOrderRecords = append(stopOrderRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    shortBuckets,
				})
			}
			if len(longBuckets) != 0 {
				stopOrderRecords = append(stopOrderRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    longBuckets,
				})
			}
		}

		// search applicable limit order
		var limitOrderRecords []record
		if len(limitOrderLowerLimits) > 0 {
			var shortBuckets []bucket
			var longBuckets []bucket
			for i := 0; i < targetRange-len(limitOrderLowerLimits); i++ {
				for j := 0; j < len(limitOrderLowerLimits); j++ {
					if oShort[i+j].LongCountPercent >= limitOrderLowerLimits[j] {
						b := bucket{
							priceRange:    oShort[i+j].Price,
							shortOrder:    oShort[i+j].ShortCountPercent,
							longOrder:     oShort[i+j].LongCountPercent,
							shortPosition: pShort[i+j].ShortCountPercent,
							longPosition:  pShort[i+j].LongCountPercent,
						}
						shortBuckets = append(shortBuckets, b)
					}
					if oLong[i+j].ShortCountPercent >= limitOrderLowerLimits[j] {
						b := bucket{
							priceRange:    oLong[i+j].Price,
							shortOrder:    oLong[i+j].ShortCountPercent,
							longOrder:     oLong[i+j].LongCountPercent,
							shortPosition: pLong[i+j].ShortCountPercent,
							longPosition:  pLong[i+j].LongCountPercent,
						}
						longBuckets = append(longBuckets, b)
					}
				}
			}
			if len(shortBuckets) != 0 {
				limitOrderRecords = append(limitOrderRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    shortBuckets,
				})
			}
			if len(longBuckets) != 0 {
				limitOrderRecords = append(limitOrderRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    longBuckets,
				})
			}
		}

		// search applicable losing position
		var losingPositionRecords []record
		if len(losingPositionLowerLimits) > 0 {
			var shortBuckets []bucket
			var longBuckets []bucket
			for i := 0; i < targetRange-len(losingPositionLowerLimits); i++ {
				for j := 0; j < len(losingPositionLowerLimits); j++ {
					if oShort[i+j].ShortCountPercent >= losingPositionLowerLimits[j] {
						b := bucket{
							priceRange:    oShort[i+j].Price,
							shortOrder:    oShort[i+j].ShortCountPercent,
							longOrder:     oShort[i+j].LongCountPercent,
							shortPosition: pShort[i+j].ShortCountPercent,
							longPosition:  pShort[i+j].LongCountPercent,
						}
						shortBuckets = append(shortBuckets, b)
					}
					if oLong[i+j].LongCountPercent >= losingPositionLowerLimits[j] {
						b := bucket{
							priceRange:    oLong[i+j].Price,
							shortOrder:    oLong[i+j].ShortCountPercent,
							longOrder:     oLong[i+j].LongCountPercent,
							shortPosition: pLong[i+j].ShortCountPercent,
							longPosition:  pLong[i+j].LongCountPercent,
						}
						longBuckets = append(longBuckets, b)
					}
				}
			}
			if len(shortBuckets) != 0 {
				losingPositionRecords = append(losingPositionRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    shortBuckets,
				})
			}
			if len(longBuckets) != 0 {
				losingPositionRecords = append(losingPositionRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    longBuckets,
				})
			}
		}

		// search applicable profiting position
		var profitingPositionRecords []record
		if len(profitingPositionLowerLimits) > 0 {
			var shortBuckets []bucket
			var longBuckets []bucket
			for i := 0; i < targetRange-len(profitingPositionLowerLimits); i++ {
				for j := 0; j < len(profitingPositionLowerLimits); j++ {
					if oShort[i+j].LongCountPercent >= profitingPositionLowerLimits[j] {
						b := bucket{
							priceRange:    oShort[i+j].Price,
							shortOrder:    oShort[i+j].ShortCountPercent,
							longOrder:     oShort[i+j].LongCountPercent,
							shortPosition: pShort[i+j].ShortCountPercent,
							longPosition:  pShort[i+j].LongCountPercent,
						}
						shortBuckets = append(shortBuckets, b)
					}
					if oLong[i+j].ShortCountPercent >= profitingPositionLowerLimits[j] {
						b := bucket{
							priceRange:    oLong[i+j].Price,
							shortOrder:    oLong[i+j].ShortCountPercent,
							longOrder:     oLong[i+j].LongCountPercent,
							shortPosition: pLong[i+j].ShortCountPercent,
							longPosition:  pLong[i+j].LongCountPercent,
						}
						longBuckets = append(longBuckets, b)
					}
				}
			}
			if len(shortBuckets) != 0 {
				profitingPositionRecords = append(profitingPositionRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    shortBuckets,
				})
			}
			if len(longBuckets) != 0 {
				profitingPositionRecords = append(profitingPositionRecords, record{
					dateTime:   dateTime,
					price:      price,
					instrument: instrument,
					buckets:    longBuckets,
				})
			}
		}

		allRecords = append(allRecords, stopOrderRecords...)
		allRecords = append(allRecords, limitOrderRecords...)
		allRecords = append(allRecords, losingPositionRecords...)
		allRecords = append(allRecords, profitingPositionRecords...)
	}

	// open file
	f, err := os.Create(buildFileName(*fileNamePrefix, *instrumentStr, *periodStr))
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
		return
	}

	// write csv
	baseHeader := []string{"date-time", "price"}
	bucketHeader := []string{"price-range", "short-order", "long-order", "short-position", "long-position"}
	bucketHeaderMaxSize := len(stopOrderLowerLimits)
	if bucketHeaderMaxSize < len(limitOrderLowerLimits) {
		bucketHeaderMaxSize = len(limitOrderLowerLimits)
	}
	if bucketHeaderMaxSize < len(losingPositionLowerLimits) {
		bucketHeaderMaxSize = len(losingPositionLowerLimits)
	}
	if bucketHeaderMaxSize < len(profitingPositionLowerLimits) {
		bucketHeaderMaxSize = len(profitingPositionLowerLimits)
	}
	if err := writeCSV(f, baseHeader, bucketHeader, bucketHeaderMaxSize, allRecords); err != nil {
		log.Fatalf("failed to write csv: %v", err)
	}
	return
}

func writeCSV(f io.Writer, baseHeader, bucketHeader []string, bucketHeaderMaxSize int, records []record) error {

	// build header
	header := baseHeader
	for i := 0; i < bucketHeaderMaxSize; i++ {
		for j, s := range bucketHeader {
			bucketHeader[j] = fmt.Sprintf("%s-%d", s, i)
		}
		header = append(header, bucketHeader...)
	}

	// build records
	var csvRecords [][]string
	csvRecords = append(csvRecords, header)
	for _, r := range records {
		var bucketRecord []string
		for _, b := range r.buckets {
			bucketRecord = append(bucketRecord, b.priceRange.PriceStr(r.instrument))
			bucketRecord = append(bucketRecord, strconv.FormatFloat(b.shortOrder, 'f', 2, 64))
			bucketRecord = append(bucketRecord, strconv.FormatFloat(b.longOrder, 'f', 2, 64))
			bucketRecord = append(bucketRecord, strconv.FormatFloat(b.shortPosition, 'f', 2, 64))
			bucketRecord = append(bucketRecord, strconv.FormatFloat(b.longPosition, 'f', 2, 64))
		}
		csvRecords = append(csvRecords, append([]string{r.dateTime.String(), r.price.PriceStr(r.instrument)}, bucketRecord...))
	}

	// write csv
	w := csv.NewWriter(f)
	if err := w.WriteAll(csvRecords); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func buildFileName(prefix, instrument, period string) string {
	return fmt.Sprintf("%s_%s_%s.csv", prefix, instrument, period)
}

type record struct {
	dateTime   time.Time
	price      oanda.Price
	instrument oanda.Instrument
	buckets    []bucket
}

type bucket struct {
	priceRange    oanda.Price
	shortOrder    float64
	longOrder     float64
	shortPosition float64
	longPosition  float64
}
