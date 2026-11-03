package resources

import (
	"bytes"
	"encoding/binary"
	"sort"
	"unicode/utf16"
)

// pixelUnits is the size of a font pixel in font units.
const pixelUnits = 128

// TTF converts the bitmap font into a TrueType font. Every pixel becomes a
// square, so the font is rendered exactly like the game when it is used with
// a font size of Height * n pixels (n = 1, 2, ...).
//
// Only the minimal set of tables needed by browsers is written: cmap, glyf,
// head, hhea, hmtx, loca, maxp, name, OS/2 and post.
func (f *BitmapFont) TTF(family string) []byte {
	type glyph struct {
		r       rune
		data    []byte
		advance int
		bbox    [4]int16 // xMin, yMin, xMax, yMax
		points  int
		contour int
	}

	// glyph 0 (.notdef) and 1 (space) are empty
	glyphs := []glyph{
		{advance: f.Width / 2 * pixelUnits},
		{r: ' ', advance: f.Width / 2 * pixelUnits},
	}
	for i := range f.Glyphs {
		r := f.Rune(i)
		width := f.GlyphWidth(i)
		// the letters start in the second column (the first one is the
		// outline), so this leaves one pixel between two letters
		g := glyph{r: r, advance: width * pixelUnits}
		if width == 0 {
			// unused glyph
			continue
		}
		g.data, g.bbox, g.points, g.contour = f.outline(i)
		glyphs = append(glyphs, g)
	}

	// glyf and loca
	var glyf bytes.Buffer
	loca := make([]uint32, 0, len(glyphs)+1)
	maxPoints, maxContours, maxAdvance := 0, 0, 0
	var xMax, yMin, yMax int16
	for _, g := range glyphs {
		loca = append(loca, uint32(glyf.Len()))
		glyf.Write(g.data)
		for glyf.Len()%4 != 0 {
			glyf.WriteByte(0)
		}
		maxPoints = max(maxPoints, g.points)
		maxContours = max(maxContours, g.contour)
		maxAdvance = max(maxAdvance, g.advance)
		xMax, yMin, yMax = max(xMax, g.bbox[2]), min(yMin, g.bbox[1]), max(yMax, g.bbox[3])
	}
	loca = append(loca, uint32(glyf.Len()))

	ascent := int16(f.Baseline * pixelUnits)
	descent := int16((f.Height - f.Baseline) * pixelUnits)
	unitsPerEm := uint16(f.Height * pixelUnits)

	head := pack(
		uint32(0x00010000), uint32(0x00010000), uint32(0), uint32(0x5F0F3CF5),
		uint16(0x000B), unitsPerEm, uint64(0), uint64(0),
		int16(0), yMin, xMax, yMax,
		uint16(0), uint16(8), int16(2), int16(1), int16(0),
	)
	hhea := pack(
		uint32(0x00010000), ascent, -descent, int16(0),
		uint16(maxAdvance), int16(0), int16(0), xMax,
		int16(1), int16(0), int16(0), int16(0), int16(0), int16(0), int16(0),
		int16(0), uint16(len(glyphs)),
	)
	maxp := pack(
		uint32(0x00010000), uint16(len(glyphs)), uint16(maxPoints), uint16(maxContours),
		uint16(0), uint16(0), uint16(2), uint16(0), uint16(0), uint16(0), uint16(0),
		uint16(0), uint16(0), uint16(0), uint16(0),
	)

	var hmtx bytes.Buffer
	for _, g := range glyphs {
		hmtx.Write(pack(uint16(g.advance), g.bbox[0]))
	}

	// cmap format 4 with one segment per character
	type mapping struct {
		r     rune
		glyph int
	}
	var mappings []mapping
	for id, g := range glyphs {
		if id > 0 && g.r < 0xFFFF {
			mappings = append(mappings, mapping{g.r, id})
		}
	}
	sort.Slice(mappings, func(i, j int) bool { return mappings[i].r < mappings[j].r })
	segments := len(mappings) + 1
	var ends, starts, deltas, offsets bytes.Buffer
	for _, m := range mappings {
		ends.Write(pack(uint16(m.r)))
		starts.Write(pack(uint16(m.r)))
		deltas.Write(pack(uint16(m.glyph - int(m.r))))
		offsets.Write(pack(uint16(0)))
	}
	ends.Write(pack(uint16(0xFFFF)))
	starts.Write(pack(uint16(0xFFFF)))
	deltas.Write(pack(uint16(1)))
	offsets.Write(pack(uint16(0)))
	searchRange, entrySelector := searchParams(segments, 2)
	subtable := pack(
		uint16(4), uint16(16+8*segments), uint16(0),
		uint16(segments*2), uint16(searchRange), uint16(entrySelector), uint16(segments*2-searchRange),
	)
	subtable = append(subtable, ends.Bytes()...)
	subtable = append(subtable, pack(uint16(0))...)
	subtable = append(subtable, starts.Bytes()...)
	subtable = append(subtable, deltas.Bytes()...)
	subtable = append(subtable, offsets.Bytes()...)
	cmap := append(pack(uint16(0), uint16(1), uint16(3), uint16(1), uint32(12)), subtable...)

	firstRune, lastRune := mappings[0].r, mappings[len(mappings)-1].r
	os2 := pack(
		uint16(4), int16(f.Width/2*pixelUnits), uint16(400), uint16(5), uint16(0),
		// subscript, superscript and strikeout
		int16(unitsPerEm/2), int16(unitsPerEm/2), int16(0), int16(0),
		int16(unitsPerEm/2), int16(unitsPerEm/2), int16(0), int16(unitsPerEm/2),
		int16(pixelUnits), int16(ascent/2),
		int16(0), // family class
		[10]byte{},
		uint32(0x00000003), uint32(0), uint32(0), uint32(0), // Basic Latin, Latin-1
		[4]byte{'X', 'C', 'O', 'M'},
		uint16(0x0040), uint16(firstRune), uint16(min(lastRune, 0xFFFF)),
		ascent, -descent, int16(0), uint16(ascent), uint16(descent),
		uint32(1), uint32(0), // code page Latin 1
		int16(ascent/2), ascent, uint16(0), uint16(' '), uint16(1),
	)

	name := nameTable(family)

	post := pack(
		uint32(0x00030000), uint32(0), int16(-pixelUnits), int16(pixelUnits),
		uint32(0), uint32(0), uint32(0), uint32(0), uint32(0),
	)

	var locaBuf bytes.Buffer
	for _, o := range loca {
		locaBuf.Write(pack(o))
	}

	return buildFont(map[string][]byte{
		"cmap": cmap,
		"glyf": glyf.Bytes(),
		"head": head,
		"hhea": hhea,
		"hmtx": hmtx.Bytes(),
		"loca": locaBuf.Bytes(),
		"maxp": maxp,
		"name": name,
		"OS/2": os2,
		"post": post,
	})
}

// outline returns the TrueType glyph data of glyph i. Set pixels are
// combined into rectangles (horizontal runs, merged with identical runs in
// the rows below).
func (f *BitmapFont) outline(i int) (data []byte, bbox [4]int16, points, contours int) {
	type rect struct{ x0, x1, y0, y1 int } // in pixels, y from the top
	var rects []rect
	open := map[[2]int]int{} // run (x0, x1) -> index of rect ending in the previous row
	for y := 0; y < f.Height; y++ {
		next := map[[2]int]int{}
		for x := 0; x < f.Width; {
			if !f.Set(i, x, y) {
				x++
				continue
			}
			start := x
			for x < f.Width && f.Set(i, x, y) {
				x++
			}
			run := [2]int{start, x}
			if idx, ok := open[run]; ok {
				rects[idx].y1 = y + 1
				next[run] = idx
			} else {
				rects = append(rects, rect{start, x, y, y + 1})
				next[run] = len(rects) - 1
			}
		}
		open = next
	}

	bbox = [4]int16{32767, 32767, -32768, -32768}
	var xs, ys []int16
	var ends []uint16
	for _, r := range rects {
		// font units, y up from the baseline
		x0, x1 := int16(r.x0*pixelUnits), int16(r.x1*pixelUnits)
		y0, y1 := int16((f.Baseline-r.y1)*pixelUnits), int16((f.Baseline-r.y0)*pixelUnits)
		// clockwise: up the left side, along the top, down the right side
		xs = append(xs, x0, x0, x1, x1)
		ys = append(ys, y0, y1, y1, y0)
		ends = append(ends, uint16(len(xs)-1))
		bbox = [4]int16{min(bbox[0], x0), min(bbox[1], y0), max(bbox[2], x1), max(bbox[3], y1)}
	}

	var buf bytes.Buffer
	buf.Write(pack(int16(len(rects)), bbox[0], bbox[1], bbox[2], bbox[3]))
	for _, e := range ends {
		buf.Write(pack(e))
	}
	buf.Write(pack(uint16(0))) // no instructions
	for range xs {
		buf.WriteByte(0x01) // on curve, 16-bit coordinates
	}
	var prev int16
	for _, x := range xs {
		buf.Write(pack(x - prev))
		prev = x
	}
	prev = 0
	for _, y := range ys {
		buf.Write(pack(y - prev))
		prev = y
	}
	return buf.Bytes(), bbox, len(xs), len(rects)
}

// nameTable returns a name table with the family name.
func nameTable(family string) []byte {
	names := []struct {
		id    uint16
		value string
	}{
		{1, family}, {2, "Regular"}, {3, family + " Regular"}, {4, family}, {6, family},
	}
	var records, strings bytes.Buffer
	for _, n := range names {
		var s bytes.Buffer
		for _, c := range utf16.Encode([]rune(n.value)) {
			s.Write(pack(c))
		}
		records.Write(pack(uint16(3), uint16(1), uint16(0x409), n.id, uint16(s.Len()), uint16(strings.Len())))
		strings.Write(s.Bytes())
	}
	header := pack(uint16(0), uint16(len(names)), uint16(6+records.Len()))
	return append(append(header, records.Bytes()...), strings.Bytes()...)
}

// buildFont writes the table directory and the tables (sorted by tag).
func buildFont(tables map[string][]byte) []byte {
	tags := make([]string, 0, len(tables))
	for tag := range tables {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	searchRange, entrySelector := searchParams(len(tags), 16)
	var out bytes.Buffer
	out.Write(pack(uint32(0x00010000), uint16(len(tags)), uint16(searchRange),
		uint16(entrySelector), uint16(len(tags)*16-searchRange)))

	offset := 12 + 16*len(tags)
	var data bytes.Buffer
	headOffset := 0
	for _, tag := range tags {
		table := tables[tag]
		if tag == "head" {
			headOffset = offset
		}
		out.WriteString(tag)
		out.Write(pack(checksum(table), uint32(offset), uint32(len(table))))
		data.Write(table)
		for data.Len()%4 != 0 {
			data.WriteByte(0)
		}
		offset = 12 + 16*len(tags) + data.Len()
	}
	out.Write(data.Bytes())

	font := out.Bytes()
	// head.checkSumAdjustment
	binary.BigEndian.PutUint32(font[headOffset+8:], 0xB1B0AFBA-checksum(font))
	return font
}

// checksum is the TrueType table checksum.
func checksum(data []byte) uint32 {
	var sum uint32
	for i := 0; i < len(data); i += 4 {
		var word [4]byte
		copy(word[:], data[i:])
		sum += binary.BigEndian.Uint32(word[:])
	}
	return sum
}

// searchParams returns searchRange and entrySelector for n entries of the
// given size.
func searchParams(n, size int) (searchRange, entrySelector int) {
	power := 1
	for power*2 <= n {
		power *= 2
		entrySelector++
	}
	return power * size, entrySelector
}

// pack encodes the values big-endian.
func pack(values ...any) []byte {
	var buf bytes.Buffer
	for _, v := range values {
		_ = binary.Write(&buf, binary.BigEndian, v)
	}
	return buf.Bytes()
}
