package cli

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/joesantosio/example-order-book/entity"
)

func Test_getRejectResponse(t *testing.T) {
	type args struct {
		userOrderID int
		userID      int
	}
	type testStruct struct {
		name string
		args args
		want []string
	}

	tests := []testStruct{
		func() testStruct {
			userOrderID := rand.Intn(10000)
			userID := rand.Intn(10000)

			return testStruct{
				name: "runs",
				args: args{userOrderID: userOrderID, userID: userID},
				want: []string{"R", strconv.Itoa(userID), strconv.Itoa(userOrderID)},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRejectResponse(tt.args.userOrderID, tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRejectResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAcknowledgeResponse(t *testing.T) {
	type args struct {
		userOrderID int
		userID      int
	}
	type testStruct struct {
		name string
		args args
		want []string
	}

	tests := []testStruct{
		func() testStruct {
			userOrderID := rand.Intn(10000)
			userID := rand.Intn(10000)

			return testStruct{
				name: "runs",
				args: args{userOrderID: userOrderID, userID: userID},
				want: []string{"A", strconv.Itoa(userID), strconv.Itoa(userOrderID)},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAcknowledgeResponse(tt.args.userOrderID, tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAcknowledgeResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTopOrderResponse(t *testing.T) {
	type args struct {
		order entity.Order
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "runs side ask",
			args: args{
				entity.Order{
					ID:    0,
					Side:  "ask",
					Price: 1,
					Size:  2,
				},
			},
			want: []string{"B", "S", "1", "2"},
		},
		{
			name: "runs side ask with no top",
			args: args{
				entity.Order{
					ID:    -1,
					Side:  "ask",
					Price: -1,
					Size:  -1,
				},
			},
			want: []string{"B", "S", "-", "-"},
		},
		{
			name: "runs side bid",
			args: args{
				entity.Order{
					ID:    0,
					Side:  "bid",
					Price: 1,
					Size:  2,
				},
			},
			want: []string{"B", "B", "1", "2"},
		},
		{
			name: "runs side bid with no top",
			args: args{
				entity.Order{
					ID:    -1,
					Side:  "bid",
					Price: -1,
					Size:  -1,
				},
			},
			want: []string{"B", "B", "-", "-"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTopOrderResponse(tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTopOrderResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeCSVFile(t *testing.T) {
	type args struct {
		p       string
		records [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "runs",
			args: args{
				p:       fmt.Sprintf("tmp_test_%d.csv", rand.Intn(10000)),
				records: [][]string{{"A"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeCSVFile(tt.args.p, tt.args.records); (err != nil) != tt.wantErr {
				t.Errorf("writeCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				os.Remove(tt.args.p)
			}
		})
	}
}
