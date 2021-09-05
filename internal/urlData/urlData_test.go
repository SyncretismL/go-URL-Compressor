package urlData

import "testing"

func TestURLCompressing(t *testing.T) {
	tests := []struct {
		name string
		arg  URLData
		want string
	}{
		{
			name: "success",
			arg: URLData{
				ID:  1,
				URL: "http://hhh.com/data/phonenumber1",
			},
			want: "http://hhh.com/data/aaaaaaaaab",
		},
		{
			name: "success large id",
			arg: URLData{
				ID:  45555545445455,
				URL: "http://hhh.com/data/phonenumber2",
			},
			want: "http://hhh.com/data/aalJMT5CzM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arg.URLCompressing(); tt.arg.URLCompressed != tt.want {
				t.Errorf("URLCompressing() = %v, want %v", tt.arg.URLCompressed, tt.want)
			}
		})
	}
}

func TestReplaceQuery(t *testing.T) {
	tests := []struct {
		name string
		arg  URLData
		want string
	}{
		{
			name: "success",
			arg: URLData{
				URL:           "http://hhh.com/data/phonenumber1",
				URLCompressed: "aaaaaaaaab",
			},
			want: "http://hhh.com/data/aaaaaaaaab",
		},
		{
			name: "success large id",
			arg: URLData{
				URL:           "http://hhh.()%)(%JALFMVAA:SFFcom/data///////phonenumber2",
				URLCompressed: "aalJMT5CzM",
			},
			want: "http://hhh.()%)(%JALFMVAA:SFFcom/data///////aalJMT5CzM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arg.ReplaceQuery(); tt.arg.URLCompressed != tt.want {
				t.Errorf("URLCompressing() = %v, want %v", tt.arg.URLCompressed, tt.want)
			}
		})
	}
}
