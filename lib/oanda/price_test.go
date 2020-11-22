package oanda

import "testing"

func TestPrice_Round(t *testing.T) {
	type inputs struct {
		price      Price
		instrument Instrument
	}
	tests := []struct {
		input    inputs
		expected Price
	}{
		{
			input:    inputs{Price(125.2225), InstrumentUSDJPY},
			expected: Price(125.223),
		},
		{
			input:    inputs{Price(125.2224), InstrumentUSDJPY},
			expected: Price(125.222),
		},
		{
			input:    inputs{Price(1.184776), InstrumentEURUSD},
			expected: Price(1.18478),
		},
		{
			input:    inputs{Price(1.184774), InstrumentEURUSD},
			expected: Price(1.18477),
		},
	}
	for i, test := range tests {
		if actual := test.input.price.Round(test.input.instrument); actual != test.expected {
			t.Errorf("#%d Round() = %v, expected: %v", i, actual, test.expected)
		}
	}
}

func TestPrice_RoundFivePips(t *testing.T) {
	type inputs struct {
		price      Price
		instrument Instrument
	}
	tests := []struct {
		input    inputs
		expected Price
	}{
		{
			input:    inputs{Price(125.125), InstrumentUSDJPY},
			expected: Price(125.150),
		},
		{
			input:    inputs{Price(125.124), InstrumentUSDJPY},
			expected: Price(125.100),
		},
		{
			input:    inputs{Price(1.18475), InstrumentEURUSD},
			expected: Price(1.18500),
		},
		{
			input:    inputs{Price(1.18474), InstrumentEURUSD},
			expected: Price(1.18450),
		},
	}
	for i, test := range tests {
		if actual := test.input.price.RoundFivePips(test.input.instrument); actual != test.expected {
			t.Errorf("#%d RoundFivePips() = %v, expected: %v", i, actual, test.expected)
		}
	}
}
