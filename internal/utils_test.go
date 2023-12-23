package internal_test

import (
	"testing"

	"github.com/redtoad/xcom-editor/internal"
)

func TestLoadDATFile(t *testing.T) {

	tests := []struct {
		name    string
		path    string
		obj     interface{}
		wantErr bool
	}{
		{
			"empty path",
			"", nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := internal.LoadDATFile(tt.path, &tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("LoadDATFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
