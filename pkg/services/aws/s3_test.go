package services

import "testing"

func TestGenerateFileName(t *testing.T) {
	type args struct {
		file     string
		fileType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				file:     "yooo",
				fileType: ".json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateFileName(tt.args.file, tt.args.fileType); got != tt.want {
				t.Errorf("GenerateFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
