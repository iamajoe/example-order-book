package cli

import (
	"testing"
)

func Test_readCSVFile(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "should error",
			args:    args{"README.md"},
			want:    -1,
			wantErr: true,
		},
		{
			name:    "runs",
			args:    args{"fixtures/1_output_balanced_book.csv"},
			want:    13,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSVFile(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCSVFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("readCSVFile() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
