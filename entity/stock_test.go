package entity

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestNewStock(t *testing.T) {
	type args struct {
		id     int
		userID int
		symbol string
		qty    int
	}
	type testStruct struct {
		name string
		args args
		want Stock
	}

	tests := []testStruct{
		func() testStruct {
			stock := Stock{
				ID:     rand.Intn(10000),
				UserID: rand.Intn(10000),
				Symbol: fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				Qty:    rand.Int() * 10000,
			}

			return testStruct{
				name: "runs",
				args: args{
					stock.ID,
					stock.UserID,
					stock.Symbol,
					stock.Qty,
				},
				want: stock,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStock(
				tt.args.id,
				tt.args.userID,
				tt.args.symbol,
				tt.args.qty,
			)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStock() = %v, want %v", got, tt.want)
			}
		})
	}
}
