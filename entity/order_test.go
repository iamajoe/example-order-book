package entity

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestNewOrder(t *testing.T) {
	type args struct {
		id        int
		userID    int
		symbol    string
		side      string
		price     int
		size      int
		isOpen    bool
		createdAt int
		updatedAt int
	}
	type testStruct struct {
		name string
		args args
		want Order
	}

	tests := []testStruct{
		func() testStruct {
			order := Order{
				ID:        rand.Intn(10000),
				UserID:    rand.Intn(10000),
				Symbol:    fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				Side:      fmt.Sprintf("tmp_side_%d", rand.Intn(10000)), // TODO: we should enum this
				Price:     rand.Int() * 10000,
				Size:      rand.Int() * 10000,
				IsOpen:    false,
				CreatedAt: rand.Int() * 10000,
				UpdatedAt: rand.Int() * 10000,
			}

			return testStruct{
				name: "runs",
				args: args{
					order.ID,
					order.UserID,
					order.Symbol,
					order.Side,
					order.Price,
					order.Size,
					order.IsOpen,
					order.CreatedAt,
					order.UpdatedAt,
				},
				want: order,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOrder(
				tt.args.id,
				tt.args.userID,
				tt.args.symbol,
				tt.args.side,
				tt.args.price,
				tt.args.size,
				tt.args.isOpen,
				tt.args.createdAt,
				tt.args.updatedAt,
			)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
