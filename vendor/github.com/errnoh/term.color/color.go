// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

// Package color handles conversion from 256 color terminal colors into other color models.
package color

import (
	"image/color"
)

// Term256 satisfies color.Color interface, complete with color.ModelFuncs for converting between color models.
type Term256 struct {
	Val uint8
}

func (c Term256) RGBA() (r, g, b, a uint32) {
	r, g, b, a = toRGBA(c.Val)
	r |= r << 8
	g |= g << 8
	b |= b << 8
	return
}

func term256Model(c color.Color) color.Color {
	if _, ok := c.(Term256); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return Term256{fromRGBA(r, g, b, a)}
}

func term256GreyscaleModel(c color.Color) color.Color {
	if val, ok := c.(Term256); ok {
		if val.Val >= 232 {
			return c
		}
	}
	r, g, b, a := c.RGBA()
	return Term256{greyscaleFromRGBA(r, g, b, a)}
}

var (
	Term256Model          color.Model = color.ModelFunc(term256Model)
	Term256GreyscaleModel color.Model = color.ModelFunc(term256GreyscaleModel)
)

// NOTE: 0 is default, 16 is actual black.
func toRGBA(val byte) (r, g, b, a uint32) {
	var tmp uint32

	switch {
	case val < 8:
		if val&1 != 0 {
			r = 187
		}
		if val&2 != 0 {
			g = 187
		}
		if val&4 != 0 {
			b = 187
		}
	case val < 16:
		r, g, b = 85, 85, 85
		if val&1 != 0 {
			r = 255
		}
		if val&2 != 0 {
			g = 255
		}
		if val&4 != 0 {
			b = 255
		}
	case val < 232:
		tmp = uint32(val) - 16
		r = tmp / 36     // "z"
		g = tmp % 36 / 6 // "y"
		b = tmp % 6      // "x"

		if r > 0 {
			r = 95 + 40*(r-1)
		}
		if g > 0 {
			g = 95 + 40*(g-1)
		}
		if b > 0 {
			b = 95 + 40*(b-1)
		}
	default:
		tmp = uint32(val) - 232
		tmp = 8 + 10*tmp
		r, g, b = tmp, tmp, tmp
	}

	return
}

func fromRGBA(r, g, b, a uint32) (val byte) {
	// 0-74 == 0, 75-114 = 1, etc
	r = r &^ 0xffffff00
	g = g &^ 0xffffff00
	b = b &^ 0xffffff00

	if r >= 35 {
		r -= 35
	}
	if g >= 35 {
		g -= 35
	}
	if b >= 35 {
		b -= 35
	}
	r, g, b = r/40, g/40, b/40

	val = 16 + uint8(r)*36
	val = val + uint8(g)*6
	val = val + uint8(b)
	return
}

// Intensity
func greyscaleFromRGBA(r, g, b, a uint32) (val byte) {
	r = r &^ 0xffffff00
	g = g &^ 0xffffff00
	b = b &^ 0xffffff00

	tmp := (r + g + b) / 3
	if tmp >= 3 {
		tmp -= 3
	} else {
		return 16
	}
	tmp = tmp / 10
	val = 232 + byte(tmp)

	if byte(tmp) > 23 {
		return 231
	}
	return
}
