package internal_test

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal"
)

// testStruct is a simple binary-serialisable struct used across tests.
type testStruct struct {
	A uint16 `struct:"uint16"`
	B uint16 `struct:"uint16"`
}

func TestLoadDATFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		obj     interface{}
		wantErr bool
	}{
		{"empty path", "", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := internal.LoadDATFile(tt.path, &tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("LoadDATFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadDATFile_RoundTrip(t *testing.T) {
	want := testStruct{A: 0x0102, B: 0x0304}
	data, err := restruct.Pack(binary.LittleEndian, &want)
	if err != nil {
		t.Fatalf("could not pack test data: %v", err)
	}

	f, err := os.CreateTemp(t.TempDir(), "*.dat")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	if _, err = f.Write(data); err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	f.Close()

	var got testStruct
	if err = internal.LoadDATFile(f.Name(), &got); err != nil {
		t.Fatalf("LoadDATFile() unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("LoadDATFile() = %+v, want %+v", got, want)
	}
}

func TestMarshall(t *testing.T) {
	obj := testStruct{A: 1, B: 2}
	data, err := internal.Marshall(&obj)
	if err != nil {
		t.Fatalf("Marshall() unexpected error: %v", err)
	}
	if len(data) != 4 {
		t.Errorf("Marshall() len = %d, want 4", len(data))
	}
}

func TestUnmarshall(t *testing.T) {
	obj := testStruct{A: 0xABCD, B: 0x1234}
	data, err := restruct.Pack(binary.LittleEndian, &obj)
	if err != nil {
		t.Fatalf("could not pack test data: %v", err)
	}

	var got testStruct
	if err = internal.Unmarshall(data, &got); err != nil {
		t.Fatalf("Unmarshall() unexpected error: %v", err)
	}
	if got != obj {
		t.Errorf("Unmarshall() = %+v, want %+v", got, obj)
	}
}

func TestSaveDATFile(t *testing.T) {
	obj := testStruct{A: 10, B: 20}
	data, err := restruct.Pack(binary.LittleEndian, &obj)
	if err != nil {
		t.Fatalf("could not pack test data: %v", err)
	}

	t.Run("creates file with correct content", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.dat")
		// pre-populate so SaveDATFile can read the original
		if err = os.WriteFile(path, data, 0644); err != nil {
			t.Fatalf("could not write initial file: %v", err)
		}
		if err = internal.SaveDATFile(path, &obj); err != nil {
			t.Fatalf("SaveDATFile() unexpected error: %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("could not read back file: %v", err)
		}
		if string(got) != string(data) {
			t.Errorf("SaveDATFile() file content mismatch")
		}
	})

	t.Run("skips write when content unchanged", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "unchanged.dat")
		if err = os.WriteFile(path, data, 0644); err != nil {
			t.Fatalf("could not write initial file: %v", err)
		}
		before, _ := os.Stat(path)
		if err = internal.SaveDATFile(path, &obj); err != nil {
			t.Fatalf("SaveDATFile() unexpected error: %v", err)
		}
		after, _ := os.Stat(path)
		if !before.ModTime().Equal(after.ModTime()) {
			t.Error("SaveDATFile() wrote file despite unchanged content")
		}
	})

	t.Run("returns error for missing file", func(t *testing.T) {
		if err = internal.SaveDATFile("/nonexistent/path.dat", &obj); err == nil {
			t.Error("SaveDATFile() expected error for missing file, got nil")
		}
	})
}
