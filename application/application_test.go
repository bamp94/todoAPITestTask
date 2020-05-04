package application

import (
	"fmt"
	"testing"
)

func TestValidateIPPort(t *testing.T) {
	type args struct {
		proxy string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "result-1", args: args{"192.168.1.1:80"}, want: true},
		{name: "result-2", args: args{"168.1.1:80"}, want: false},
		{name: "result-3", args: args{"1.1:80"}, want: false},
		{name: "result-4", args: args{"1:80"}, want: false},
		{name: "result-5", args: args{""}, want: false},
		{name: "result-6", args: args{"google.com"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.want == true)
			if got := validateIPPort(tt.args.proxy); got != tt.want {
				t.Errorf("EmailToAscii() = %v, want %v", got, tt.want)
			}
		})
	}
}
