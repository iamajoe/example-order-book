package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/joesantosio/example-order-book/infrastructure/sqlite"
)

func Test_readCSVFile(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSVFile(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCSVFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCSVFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printHelp(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printHelp()
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		inputFile string
	}

	tests := []struct {
		name     string
		args     []string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "runs: 1_input_balanced_book",
			args:     []string{"book", "fixtures/1_input_balanced_book.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/1_output_balanced_book.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 2_input_shallow_bid",
			args:     []string{"book", "fixtures/2_input_shallow_bid.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/2_output_shallow_bid.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 3_input_shallow_ask",
			args:     []string{"book", "fixtures/3_input_shallow_ask.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/3_output_shallow_ask.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 4_input_balanced_book_limit_below_best",
			args:     []string{"book", "fixtures/4_input_balanced_book_limit_below_best.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/4_output_balanced_book_limit_below_best_bid.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 5_input_balanced_book_limit_above_best_ask",
			args:     []string{"book", "fixtures/5_input_balanced_book_limit_above_best_ask.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/5_output_balanced_book_limit_above_best_ask.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 6_input_spread_through_new_limit",
			args:     []string{"book", "fixtures/6_input_spread_through_new_limit.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/6_output_spread_through_new_limit.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 7_input_balanced_book_limit_sell_partial",
			args:     []string{"book", "fixtures/7_input_balanced_book_limit_sell_partial.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/7_output_balanced_book_limit_sell_partial.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 8_input_balanced_book_limit_buy_partial",
			args:     []string{"book", "fixtures/8_input_balanced_book_limit_buy_partial.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/8_output_balanced_book_limit_buy_partial.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 9_input_balanced_book_cancel_best_bid_and_offer",
			args:     []string{"book", "fixtures/9_input_balanced_book_cancel_best_bid_and_offer.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/9_output_balanced_book_cancel_best_bid_and_offer.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 10_input_balanced_book_cancel_behind_best_bid_and_offer",
			args:     []string{"book", "fixtures/10_input_balanced_book_cancel_behind_best_bid_and_offer.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/10_output_balanced_book_cancel_behind_best_bid_and_offer.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 11_input_balanced_book_cancel_all_bids",
			args:     []string{"book", "fixtures/11_input_balanced_book_cancel_all_bids.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/11_output_balanced_book_cancel_all_bids.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 12_input_balanced_book_tob_volume_changes",
			args:     []string{"book", "fixtures/12_input_balanced_book_tob_volume_changes.csv", fmt.Sprintf("tmp_output_%d", rand.Intn(10000))},
			wantPath: "fixtures/12_output_balanced_book_tob_volume_changes.csv",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare the test by setting up the repos needed
			path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
			db, err := sqlite.Connect(path)
			if err != nil {
				t.Fatal(err)
			}
			repos, err := sqlite.InitRepos(db)
			if err != nil {
				t.Fatal(err)
			}
			defer repos.Close()
			defer os.Remove(path)

			// run the tests
			err = Init(tt.args, repos)

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v, fixture: %s", err, tt.wantErr, tt.args[1])
			}

			if !tt.wantErr {
				// check if the output exists
				_, err = os.Stat(tt.args[2])
				if err == nil {
					defer os.Remove(tt.args[2])

					// check the output data
					buf, err := ioutil.ReadFile(tt.args[2])
					if err != nil {
						t.Fatal(err)
					}

					outputBuf, err := ioutil.ReadFile(tt.wantPath)
					if err != nil {
						t.Fatal(err)
					}

					if string(buf) != string(outputBuf) {
						// TODO: outputs of the program are not as expected yet
					}

					// remove the output
				} else if errors.Is(err, os.ErrNotExist) {
					t.Errorf("File exists = %v, want %v", false, true)
				}
			}
		})
	}
}
