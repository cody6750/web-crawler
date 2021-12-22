package webcrawler

import "testing"

func Test_writeURLResponseToFile(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				url: "https://www.newegg.com/p/pl?d=rtx+3080&LeftPriceRange=1000+",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeURLResponseToFile(tt.args.url)
		})
	}
}
