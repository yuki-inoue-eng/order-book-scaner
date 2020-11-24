package oanda

import (
	"reflect"
	"testing"
)

func TestBook_ExtractBucketVicinityOfPrice(t *testing.T) {
	type fields struct {
		Instrument Instrument
		Price      Price
		Buckets    []BookBucket
	}
	type args struct {
		price Price
		n     int
	}
	tests := []struct {
		fields    fields
		args      args
		wantShort []BookBucket
		wantLong  []BookBucket
		wantErr   bool
	}{
		{
			fields: fields{
				Instrument: InstrumentUSDJPY,
				Buckets: []BookBucket{
					{Price: 99.85},
					{Price: 99.90},
					{Price: 99.95},
					{Price: 100.00},
					{Price: 100.05},
					{Price: 100.10},
					{Price: 100.15},
					{Price: 100.20},
				},
			},
			args: args{price: 100.001, n: 3},
			wantShort: []BookBucket{
				{Price: 100.00},
				{Price: 99.95},
				{Price: 99.90},
			},
			wantLong: []BookBucket{
				{Price: 100.05},
				{Price: 100.10},
				{Price: 100.15},
			},
			wantErr: false,
		},
		{
			fields: fields{
				Instrument: InstrumentEURGBP,
				Buckets: []BookBucket{
					{Price: 0.90500},
					{Price: 0.90550},
					{Price: 0.90600},
					{Price: 0.90650},
					{Price: 0.90700},
					{Price: 0.90750},
					{Price: 0.90800},
				},
			},
			args: args{price: 0.90670, n: 3},
			wantShort: []BookBucket{
				{Price: 0.90650},
				{Price: 0.90600},
				{Price: 0.90550},
			},
			wantLong: []BookBucket{
				{Price: 0.90700},
				{Price: 0.90750},
				{Price: 0.90800},
			},
			wantErr: false,
		},
	}
	for i, tt := range tests {
		o := &Book{
			Instrument: tt.fields.Instrument,
			Buckets:    tt.fields.Buckets,
		}
		gotShort, gotLong, err := o.ExtractBucketVicinityOfPrice(tt.args.price, tt.args.n)
		if (err != nil) != tt.wantErr {
			t.Errorf("#%d ExtractBucketVicinityOfPrice() error = %v, wantErr %v", i, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(gotShort, tt.wantShort) {
			t.Errorf("#%d ExtractBucketVicinityOfPrice() gotShort = %v, want %v", i, gotShort, tt.wantShort)
		}
		if !reflect.DeepEqual(gotLong, tt.wantLong) {
			t.Errorf("#%d ExtractBucketVicinityOfPrice() gotLong = %v, want %v", i, gotLong, tt.wantLong)
		}
	}
}
