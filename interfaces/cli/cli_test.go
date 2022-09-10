package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strings"
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
			args:     []string{"book", "fixtures/1_input_balanced_book.csv"},
			wantPath: "fixtures/1_output_balanced_book.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 2_input_shallow_bid",
			args:     []string{"book", "fixtures/2_input_shallow_bid.csv"},
			wantPath: "fixtures/2_output_shallow_bid.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 3_input_shallow_ask",
			args:     []string{"book", "fixtures/3_input_shallow_ask.csv"},
			wantPath: "fixtures/3_output_shallow_ask.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 4_input_balanced_book_limit_below_best",
			args:     []string{"book", "fixtures/4_input_balanced_book_limit_below_best.csv"},
			wantPath: "fixtures/4_output_balanced_book_limit_below_best_bid.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 5_input_balanced_book_limit_above_best_ask",
			args:     []string{"book", "fixtures/5_input_balanced_book_limit_above_best_ask.csv"},
			wantPath: "fixtures/5_output_balanced_book_limit_above_best_ask.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 6_input_spread_through_new_limit",
			args:     []string{"book", "fixtures/6_input_spread_through_new_limit.csv"},
			wantPath: "fixtures/6_output_spread_through_new_limit.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 7_input_balanced_book_limit_sell_partial",
			args:     []string{"book", "fixtures/7_input_balanced_book_limit_sell_partial.csv"},
			wantPath: "fixtures/7_output_balanced_book_limit_sell_partial.csv",
			wantErr:  false,
		},
		{
			name:     "runs: 8_input_balanced_book_limit_buy_partial",
			args:     []string{"book", "fixtures/8_input_balanced_book_limit_buy_partial.csv"},
			wantPath: "fixtures/8_output_balanced_book_limit_buy_partial.csv",
			wantErr:  false,
		},
		// {
		// 	name:     "runs: 9_input_balanced_book_cancel_best_bid_and_offer",
		// 	args:     []string{"book", "fixtures/9_input_balanced_book_cancel_best_bid_and_offer.csv"},
		// 	wantPath: "fixtures/9_output_balanced_book_cancel_best_bid_and_offer.csv",
		// 	wantErr:  false,
		// },
		{
			name:     "runs: 10_input_balanced_book_cancel_behind_best_bid_and_offer",
			args:     []string{"book", "fixtures/10_input_balanced_book_cancel_behind_best_bid_and_offer.csv"},
			wantPath: "fixtures/10_output_balanced_book_cancel_behind_best_bid_and_offer.csv",
			wantErr:  false,
		},
		// {
		// 	name:     "runs: 11_input_balanced_book_cancel_all_bids",
		// 	args:     []string{"book", "fixtures/11_input_balanced_book_cancel_all_bids.csv"},
		// 	wantPath: "fixtures/11_output_balanced_book_cancel_all_bids.csv",
		// 	wantErr:  false,
		// },
		// {
		// 	name:     "runs: 12_input_balanced_book_tob_volume_changes",
		// 	args:     []string{"book", "fixtures/12_input_balanced_book_tob_volume_changes.csv"},
		// 	wantPath: "fixtures/12_output_balanced_book_tob_volume_changes.csv",
		// 	wantErr:  false,
		// },
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
			resultPath := tt.wantPath + "_result.csv"
			tt.args = append(tt.args, resultPath)
			err = Init(tt.args, repos)

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v, fixture: %s", err, tt.wantErr, tt.args[1])
			}

			if !tt.wantErr {
				// check if the output exists
				_, err = os.Stat(resultPath)
				if err == nil {
					// TODO: uncomment after making sure all is as expected
					// defer os.Remove(resultPath)

					// check the output data
					buf, err := ioutil.ReadFile(resultPath)
					if err != nil {
						t.Fatal(err)
					}

					outputBuf, err := ioutil.ReadFile(tt.wantPath)
					if err != nil {
						t.Fatal(err)
					}

					// normalize the data as strings so we can easily check for equality
					bufString := strings.TrimSpace(string(buf))
					bufString = strings.ReplaceAll(bufString, " ", "")
					bufString = strings.TrimSuffix(bufString, "\n")

					outputBufString := strings.TrimSpace(string(outputBuf))
					outputBufString = strings.ReplaceAll(outputBufString, " ", "")
					outputBufString = strings.TrimSuffix(outputBufString, "\n")

					if bufString != outputBufString {
						t.Errorf("Data = %v, want %v", bufString, outputBufString)
						t.Fatal("Data not what expected")
					}
				} else if errors.Is(err, os.ErrNotExist) {
					t.Errorf("File exists = %v, want %v", false, true)
				}
			}
		})
	}
}
