package hosts

import (
	"reflect"
	"testing"
)

func init() {
	SetLogLevel(6)
}

func TestNewHost(t *testing.T) {
	type args struct {
		entry string
	}
	tests := []struct {
		name    string
		args    args
		want    *Host
		wantErr bool
	}{
		{
			"valid entry",
			args{entry: "127.0.0.1	localhost"},
			&Host{Address: "127.0.0.1", Hostnames: "localhost"},
			false,
		},
		{
			"invalid entry",
			args{entry: "# comment"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHost(tt.args.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
