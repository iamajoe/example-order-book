package config

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		getenv func(string) string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "runs",
			args: args{
				getenv: func(s string) string {
					switch s {
					case "DB_PATH":
						return "dbpath"
					}

					return ""
				},
			},
			want:    Config{"dbpath"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.getenv)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
