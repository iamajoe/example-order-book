package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func Test_initRepos(t *testing.T) {
	type args struct {
		dbPath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "runs",
			args:    args{fmt.Sprintf("tmp_test_%d.db", rand.Intn(10000))},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initRepos(tt.args.dbPath)
			defer got.Close()
			defer os.Remove(tt.args.dbPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("initRepos() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want && got.GetOrder() == nil {
				t.Errorf("GetOrder() = %v, want %v", got.GetOrder() != nil, got.GetOrder() == nil)
			}
		})
	}
}
