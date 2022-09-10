package sqlite

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func Test_repositoryStock_Create(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	repo, err := createRepositoryStock(db)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		userID int
		symbol string
		qty    int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "runs",
			args: args{
				userID: rand.Intn(10000),
				symbol: fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				qty:    rand.Int() * 10000,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Create(tt.args.userID, tt.args.symbol, tt.args.qty)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryStock.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got >= 0) != tt.want {
				t.Errorf("repositoryStock.Create() = %v, want %v", got >= 0, tt.want)
			}
		})
	}
}

func Test_repositoryStock_Update(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	repo, err := createRepositoryStock(db)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		userID int
		symbol string
		qty    int
	}
	tests := []struct {
		name         string
		args         args
		createBefore bool
		want         bool
		wantErr      bool
	}{
		{
			name: "runs",
			args: args{
				userID: rand.Intn(10000),
				symbol: fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				qty:    rand.Int() * 10000,
			},
			createBefore: true,
			want:         true,
			wantErr:      false,
		},
		{
			name: "runs and creates if it doesnt exist",
			args: args{
				userID: rand.Intn(10000),
				symbol: fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				qty:    rand.Int() * 10000,
			},
			createBefore: false,
			want:         true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createBefore {
				_, err := repo.Create(tt.args.userID, tt.args.symbol, rand.Intn(10000))
				if err != nil {
					t.Fatal(err)
				}
			}

			got, err := repo.Update(tt.args.userID, tt.args.symbol, tt.args.qty)
			if err != nil {
				t.Fatal(err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryStock.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("repositoryStock.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repositoryStock_RemoveTable(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{"runs", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := createRepositoryStock(db)
			if err != nil {
				t.Fatal(err)
			}

			got, err := repo.removeTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryStock.RemoveTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repositoryStock.RemoveTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createRepositoryStock(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{"runs", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createRepositoryStock(db)
			if (err != nil) != tt.wantErr {
				t.Errorf("createRepositoryStock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("createRepositoryStock() = %v, want %v", (got != nil), tt.want)
			}
		})
	}
}
