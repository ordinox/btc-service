// Copyright 2017 Equim. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package leb128 implements LEB128 codec.
//
// The encode part is a fork of https://golang.org/src/cmd/internal/dwarf/dwarf.go
package runes // import "ekyu.moe/leb128"

var sevenbits = [...]byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
	0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
}

// AppendUleb128 appends v to b using unsigned LEB128 encoding.
func AppendUleb128(b []byte, v uint64) []byte {
	// If it's less than or equal to 7-bit
	if v < 0x80 {
		return append(b, sevenbits[v])
	}

	for {
		c := uint8(v & 0x7f)
		v >>= 7

		if v != 0 {
			c |= 0x80
		}

		b = append(b, c)

		if c&0x80 == 0 {
			break
		}
	}

	return b
}

// AppendSleb128 appends v to b using signed LEB128 encoding.
func AppendSleb128(b []byte, v int64) []byte {
	// If it's less than or equal to 7-bit
	if v >= 0 && v <= 0x3f {
		return append(b, sevenbits[v])
	} else if v < 0 && v >= ^0x3f {
		return append(b, sevenbits[0x80+v])
	}

	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7

		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}

		b = append(b, c)

		if c&0x80 == 0 {
			break
		}
	}

	return b
}

// DecodeUleb128 decodes b to u with unsigned LEB128 encoding and returns the
// number of bytes read. On error (bad encoded b), n will be 0 and therefore u
// must not be trusted.
func DecodeUleb128(b []byte) (u uint64, n uint8) {
	l := uint8(len(b) & 0xff)
	// The longest LEB128 encoded sequence is 10 byte long (9 0xff's and 1 0x7f)
	// so make sure we won't overflow.
	if l > 10 {
		l = 10
	}

	var i uint8
	for i = 0; i < l; i++ {
		u |= uint64(b[i]&0x7f) << (7 * i)
		if b[i]&0x80 == 0 {
			n = uint8(i + 1)
			return
		}
	}

	return
}

// DecodeSleb128 decodes b to s with signed LEB128 encoding and returns the
// number of bytes read. On error (bad encoded b), n will be 0 and therefore s
// must not be trusted.
func DecodeSleb128(b []byte) (s int64, n uint8) {
	l := uint8(len(b) & 0xff)
	if l > 10 {
		l = 10
	}

	var i uint8
	for i = 0; i < l; i++ {
		s |= int64(b[i]&0x7f) << (7 * i)
		if b[i]&0x80 == 0 {
			// If it's signed
			if b[i]&0x40 != 0 {
				s |= ^0 << (7 * (i + 1))
			}
			n = uint8(i + 1)
			return
		}
	}

	return
}
