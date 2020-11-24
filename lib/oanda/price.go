package oanda

import (
	"math"
	"strconv"
)

type Price float64

type Pips float64 // valid up to the first minority

func (p Pips) PipsToPrice(instrument Instrument) Price {
	if instrument == InstrumentUSDJPY {
		return Price(math.Ceil(float64(p)) / 100)
	}
	if instrument == InstrumentEURJPY {
		return Price(math.Ceil(float64(p)) / 100)
	}
	if instrument == InstrumentAUDJPY {
		return Price(math.Ceil(float64(p)) / 100)
	}
	if instrument == InstrumentGBPJPY {
		return Price(math.Ceil(float64(p)) / 100)
	}
	if instrument == InstrumentEURUSD {
		return Price(math.Ceil(float64(p)) / 10000)
	}
	if instrument == InstrumentGBPUSD {
		return Price(math.Ceil(float64(p)) / 10000)
	}
	if instrument == InstrumentAUDUSD {
		return Price(math.Ceil(float64(p)) / 10000)
	}
	if instrument == InstrumentNZDUSD {
		return Price(math.Ceil(float64(p)) / 10000)
	}
	if instrument == InstrumentEURGBP {
		return Price(math.Ceil(float64(p)) / 10000)
	}
	return 0
}

func (p Price) Round(instrument Instrument) Price {
	r := 10 / float64(Pips(1).PipsToPrice(instrument))
	return Price(math.Round(float64(p)*r) / r)
}

func (p Price) RoundFivePips(instrument Instrument) Price {
	r := float64(Pips(1).PipsToPrice(instrument) * 10)
	return Price(math.Round(float64(p)*2/r) * r / 2).Round(instrument)
}

// PriceStr converts string price. (0.1 pips units)
func (p Price) PriceStr(instrument Instrument) string {
	r := math.Floor(1 / float64(Pips(1).PipsToPrice(instrument)))
	rDigits := int(math.Log10(r)) + 1
	return strconv.FormatFloat(float64(p), 'f', rDigits, 64)
}
