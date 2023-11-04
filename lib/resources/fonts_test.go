package resources_test

import (
	"image"
	"reflect"
	"testing"

	"github.com/redtoad/xcom-editor/lib/resources"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font"
)

func TestFont_ImplementsFaceInterface(t *testing.T) {
	assert.Implements(t, (*font.Face)(nil), new(resources.Font))
}

func TestFindBoundingBox(t *testing.T) {
	tests := []struct {
		name   string
		pixels []byte
		width  int
		want   image.Rectangle
	}{
		{
			"one pixel in every corner",
			[]byte{
				0x01, 0x00, 0x00, 0x02,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x05,
			},
			4,
			image.Rect(0, 0, 3, 3),
		},
		{
			"one pixel",
			[]byte{
				0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			4,
			image.Rect(1, 0, 1, 0),
		},
		{
			"block in middle",
			[]byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x03, 0x00,
				0x00, 0x04, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			4,
			image.Rect(1, 1, 2, 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBb := resources.FindBoundingBox(tt.pixels, tt.width); !reflect.DeepEqual(gotBb, tt.want) {
				t.Errorf("FindBoundingBox() = %v, want %v", gotBb, tt.want)
			}
		})
	}
}
