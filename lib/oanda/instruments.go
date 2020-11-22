package oanda

type Instrument string

const (
	InstrumentUSDJPY  = Instrument("USD_JPY")
	InstrumentEURJPY  = Instrument("EUR_JPY")
	InstrumentAUDJPY  = Instrument("AUD_JPY")
	InstrumentGBPJPY  = Instrument("GBP_JPY")
	InstrumentEURUSD  = Instrument("EUR_USD")
	InstrumentGBPUSD  = Instrument("GBP_USD")
	InstrumentAUDUSD  = Instrument("AUD_USD")
	InstrumentNZDUSD  = Instrument("NZD_USD")
	InstrumentEURGBP  = Instrument("EUR_GBP")
	InstrumentUNKNOWN = Instrument("UNKNOWN")
)

var instrumentsMap = map[string]Instrument{
	"USD_JPY": InstrumentUSDJPY,
	"EUR_JPY": InstrumentEURJPY,
	"AUD_JPY": InstrumentAUDJPY,
	"GBP_JPY": InstrumentGBPJPY,
	"EUR_USD": InstrumentEURUSD,
	"GBP_USD": InstrumentGBPUSD,
	"AUD_USD": InstrumentAUDUSD,
	"NZD_USD": InstrumentNZDUSD,
	"EUR_GBP": InstrumentEURGBP,
	"UNKNOWN": InstrumentUNKNOWN,
}

func ToInstrument(str string) Instrument {
	for k, v := range instrumentsMap {
		if k == str {
			return v
		}
	}
	return InstrumentUNKNOWN
}
