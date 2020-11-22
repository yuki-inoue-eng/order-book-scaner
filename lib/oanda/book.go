package oanda

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type book struct {
	Instrument  string    `json:"instrument"`
	Time        time.Time `json:"time"`
	Price       string    `json:"price"`
	BucketWidth string    `json:"bucketWidth"`
	Buckets     []bucket  `json:"buckets"`
}

type retrievedOrderBook struct {
	Book book `json:"orderBook"`
}

type retrievedPositionBook struct {
	Book book `json:"positionBook"`
}

type bucket struct {
	Price             string `json:"price"`
	LongCountPercent  string `json:"longCountPercent"`
	ShortCountPercent string `json:"shortCountPercent"`
}
type Book struct {
	Instrument Instrument
	Time       time.Time
	Price      Price
	Buckets    []BookBucket
}
type BookBucket struct {
	Price             Price   `json:"price"`
	LongCountPercent  float64 `json:"longCountPercent"`
	ShortCountPercent float64 `json:"shortCountPercent"`
}

func (b *book) toBook() (*Book, error) {
	price, err := strconv.ParseFloat(b.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price to float64: %v", err)
	}
	var buckets []BookBucket
	for _, bu := range b.Buckets {
		p, err := strconv.ParseFloat(bu.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bucket price to type of float64: %v", err)
		}
		l, err := strconv.ParseFloat(bu.LongCountPercent, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse long count percent to float64: %v", err)
		}
		s, err := strconv.ParseFloat(bu.ShortCountPercent, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse short count percent to float64: %v", err)
		}
		buckets = append(buckets, BookBucket{
			Price(p),
			l,
			s,
		})
	}
	return &Book{
		Instrument(b.Instrument),
		b.Time,
		Price(price),
		buckets,
	}, nil
}

func (o *Book) ExtractBucket(maxPrice, minPrice float64) {
	var buckets []BookBucket
	for _, b := range o.Buckets {
		if Price(maxPrice) >= b.Price && b.Price >= Price(minPrice) {
			buckets = append(buckets, b)
		}
	}
	o.Buckets = buckets
}

func (o *Book) ExtractBucketVicinityOfPrice(price Price, n int) (short, long []BookBucket, err error) {
	var lowerBuckets []BookBucket
	var higherBuckets []BookBucket
	for i, b := range o.Buckets {
		if b.Price > price {
			lowerBuckets = o.Buckets[:i-1]
			higherBuckets = o.Buckets[i-1:]
			break
		}
	}
	for i, j := 0, len(lowerBuckets)-1; i < j; i, j = i+1, j-1 {
		lowerBuckets[i], lowerBuckets[j] = lowerBuckets[j], lowerBuckets[i]
	}
	if len(lowerBuckets[:n]) < n {
		return nil, nil, fmt.Errorf("price is too low: lowerBuckets[%d] is not exist", n-1)
	}
	if len(higherBuckets[:n]) < n {
		return nil, nil, fmt.Errorf("price is too high: higherBuckets[%d] is not exist", n-1)
	}
	return lowerBuckets[:n], higherBuckets[:n], nil
}

func (c *Client) FetchOrderBook(instrument Instrument, dateTime *time.Time) (*Book, error) {
	body, err := c.fetchOrderBook(instrument, dateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order book: %v", err)
	}
	var rb retrievedOrderBook
	if err := json.Unmarshal(body, &rb); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal: %v", err)
	}
	ob, err := rb.Book.toBook()
	if err != nil {
		return nil, fmt.Errorf("failed to convert book to order book: %v", err)
	}
	return ob, nil
}

func (c *Client) FetchPositionBook(instrument Instrument, dateTime *time.Time) (*Book, error) {
	body, err := c.fetchPositionBook(instrument, dateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch position book: %v", err)
	}
	var rb retrievedPositionBook
	if err := json.Unmarshal(body, &rb); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal: %v", err)
	}
	ob, err := rb.Book.toBook()
	if err != nil {
		return nil, fmt.Errorf("failed to convert book to position book: %v", err)
	}
	return ob, nil
}

func (c *Client) FetchOrderBookJSON(instrument Instrument, dateTime *time.Time) ([]byte, error) {
	return c.fetchOrderBook(instrument, dateTime)
}
