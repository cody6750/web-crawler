package webcrawler

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type args struct {
		v interface{}
		c *FilterConfiguration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "String test",
			args: args{
				v: "string",
				c: &FilterConfiguration{
					Contains:              "st",
					IsEqualTo:             "string",
					IsNotEqualTo:          "strin",
					ConvertStringToNumber: "false",
				},
			},
			want: true,
		},
		{
			name: "int test",
			args: args{
				v: "6",
				c: &FilterConfiguration{
					Contains:              "st",
					IsEqualTo:             6,
					IsNotEqualTo:          3,
					ConvertStringToNumber: "true",
				},
			},
			want: true,
		},
		{
			name: "float test",
			args: args{
				v: "6.0",
				c: &FilterConfiguration{
					Contains:              "st",
					IsEqualTo:             6.0,
					IsNotEqualTo:          6.1,
					ConvertStringToNumber: "true",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.args.v, tt.args.c); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringToNunber(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "float",
			args: args{
				s: "float $1,062.2",
			},
		},
		{
			name: "int",
			args: args{
				"int 83",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConvertStringToNunber(tt.args.s)
		})
	}
}

func Test_isURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: " ",
			want: false,
			args: args{
				url: "google.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isURL(tt.args.url); got != tt.want {
				t.Errorf("isURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
