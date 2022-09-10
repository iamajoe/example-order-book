package sqlite

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func Test_repositoryOrder_Create(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	repo, err := createRepositoryOrder(db)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		userOrderID int
		userID      int
		symbol      string
		side        string
		price       int
		size        int
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
				userOrderID: rand.Intn(10000),
				userID:      rand.Intn(10000),
				symbol:      fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				side:        fmt.Sprintf("tmp_side_%d", rand.Intn(10000)), // TODO: we should enum this
				price:       rand.Int() * 10000,
				size:        rand.Int() * 10000,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Create(tt.args.userOrderID, tt.args.userID, tt.args.symbol, tt.args.side, tt.args.price, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryOrder.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got >= 0) != tt.want {
				t.Errorf("repositoryOrder.Create() = %v, want %v", got >= 0, tt.want)
			}
		})
	}
}

func Test_repositoryOrder_Cancel(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	repo, err := createRepositoryOrder(db)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		userOrderID int
		userID      int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "runs",
			args:    args{rand.Intn(10000), rand.Intn(10000)},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare the test
			_, err := repo.Create(
				rand.Intn(10000),
				tt.args.userID,
				fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				fmt.Sprintf("tmp_side_%d", rand.Intn(10000)),
				rand.Intn(10000),
				rand.Intn(10000),
			)
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Create(
				tt.args.userOrderID,
				tt.args.userID,
				fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				fmt.Sprintf("tmp_side_%d", rand.Intn(10000)),
				rand.Intn(10000),
				rand.Intn(10000),
			)
			if err != nil {
				t.Fatal(err)
			}

			// run the test
			got, err := repo.Cancel(tt.args.userOrderID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryOrder.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repositoryOrder.Cancel() = %v, want %v", got, tt.want)
			}

			// check if anything on the db
			ids := []int{}

			rows, err := repo.db.db.Query("SELECT id FROM "+repo.tableName+" WHERE userid=$1", tt.args.userID)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				err = rows.Scan(&id)
				if err != nil {
					t.Fatal(err)
				}

				if id == tt.args.userOrderID {
					t.Fatal(errors.New("id should have been removed"))
				}

				ids = append(ids, id)
			}

			if len(ids) > 1 || ids[0] == tt.args.userOrderID {
				t.Errorf("rows length = %v, want %v", len(ids), 1)
			}
		})
	}
}

func Test_repositoryOrder_Empty(t *testing.T) {
	path := fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))
	db, err := Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(path)

	repo, err := createRepositoryOrder(db)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "runs",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare the test
			_, err = repo.Create(
				rand.Intn(10000),
				rand.Intn(10000),
				fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				fmt.Sprintf("tmp_side_%d", rand.Intn(10000)),
				rand.Intn(10000),
				rand.Intn(10000),
			)
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Create(
				rand.Intn(10000),
				rand.Intn(10000),
				fmt.Sprintf("tmp_symbol_%d", rand.Intn(10000)),
				fmt.Sprintf("tmp_side_%d", rand.Intn(10000)),
				rand.Intn(10000),
				rand.Intn(10000),
			)
			if err != nil {
				t.Fatal(err)
			}

			// run the test
			got, err := repo.Empty()
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryOrder.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repositoryOrder.Cancel() = %v, want %v", got, tt.want)
			}

			// check if anything on the db
			rows, err := repo.db.db.Query("SELECT id FROM " + repo.tableName)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				err = rows.Scan(&id)
				if err != nil {
					t.Fatal(err)
				}

				if id != -1 {
					t.Fatal(errors.New("table should be empty"))
				}
			}
		})
	}
}

func Test_repositoryOrder_RemoveTable(t *testing.T) {
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
			repo, err := createRepositoryOrder(db)
			if err != nil {
				t.Fatal(err)
			}

			got, err := repo.removeTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryOrder.RemoveTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repositoryOrder.RemoveTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createRepositoryOrder(t *testing.T) {
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
			got, err := createRepositoryOrder(db)
			if (err != nil) != tt.wantErr {
				t.Errorf("createRepositoryOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("createRepositoryOrder() = %v, want %v", (got != nil), tt.want)
			}
		})
	}
}
