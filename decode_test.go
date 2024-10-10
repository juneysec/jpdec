package jpdec_test

import (
	"fmt"
	"testing"

	"github.com/juneysec/jpdec"
	"github.com/stretchr/testify/assert"
)

// DetectEncoding() のサンプル
func ExampleDetectEncoding() {
	byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
	enc, _ := jpdec.DetectEncoding(byteArray)
	fmt.Println(enc)

	// Output:
	// Shift_JIS
}

// Decode() のサンプル
func ExampleDecode() {
	byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
	str, _ := jpdec.Decode(byteArray)
	fmt.Println(str)

	// Output:
	// こんにちは
}

// EUC-JP のエンコーディング判別と変換のテスト
func TestDetectEncodingEUCJP(t *testing.T) {
	byteArray := []byte{0xA4, 0xB3, 0xA4, 0xF3, 0xA4, 0xCB, 0xA4, 0xC1, 0xA4, 0xCF}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingEUCJP, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// JIS のエンコーディング判別と変換のテスト
func TestDetectEncodingJIS(t *testing.T) {
	byteArray := []byte{0x1B, 0x24, 0x42, 0x24, 0x33, 0x24, 0x73, 0x24, 0x4B, 0x24, 0x41, 0x24, 0x4F, 0x1B, 0x28, 0x42}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingJIS, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// Shift_JIS のエンコーディング判別と変換のテスト
func TestDetectEncodingShiftJIS(t *testing.T) {
	byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingShiftJIS, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// UTF8 のエンコーディング判別と変換のテスト
func TestDetectEncodingUTF8(t *testing.T) {
	byteArray := []byte{0xE3, 0x81, 0x93, 0xE3, 0x82, 0x93, 0xE3, 0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF8, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// UTF8(BOM有り) のエンコーディング判別と変換のテスト
func TestDetectEncodingUTF8BOM(t *testing.T) {
	byteArray := []byte{0xEF, 0xBB, 0xBF, 0xE3, 0x81, 0x93, 0xE3, 0x82, 0x93, 0xE3, 0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF8BOM, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// UTF16(BOM無し) のエンコーディング判別と変換のテスト
func TestDetectEncodingUTF16LE(t *testing.T) {
	byteArray := []byte{0x53, 0x30, 0x93, 0x30, 0x6B, 0x30, 0x61, 0x30, 0x6F, 0x30}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF16LE, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// 一 (U+4E00)を含むケース
func TestDetectEncodingUTF16LE4E00(t *testing.T) {
	byteArray := []byte{0xA4, 0x30, 0xF3, 0x30, 0xB9, 0x30, 0xC8, 0x30, 0xFC, 0x30, 0xEB, 0x30, 0x67, 0x30, 0x4D, 0x30, 0x8B, 0x30, 0x09, 0x67, 0xB9, 0x52, 0x6A, 0x30, 0xC7, 0x30, 0xA3, 0x30, 0xB9, 0x30, 0xC8, 0x30, 0xEA, 0x30, 0xD3, 0x30, 0xE5, 0x30, 0xFC, 0x30, 0xB7, 0x30, 0xE7, 0x30, 0xF3, 0x30, 0x6E, 0x30, 0x00, 0x4E, 0xA7, 0x89, 0x92, 0x30, 0x21, 0x6B, 0x6B, 0x30, 0x3A, 0x79, 0x57, 0x30, 0x7E, 0x30, 0x59, 0x30, 0x02, 0x30, 0x0D, 0x00, 0x0A, 0x00, 0x27, 0x00, 0x77, 0x00, 0x73, 0x00, 0x6C, 0x00, 0x2E, 0x00, 0x65, 0x00, 0x78, 0x00, 0x65, 0x00, 0x20, 0x00, 0x2D, 0x00, 0x2D, 0x00, 0x69, 0x00, 0x6E, 0x00, 0x73, 0x00, 0x74, 0x00, 0x61, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x20, 0x00, 0x3C, 0x00}
	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF16LE, enc)
}

// UTF16(BOM有り) のエンコーディング判別と変換のテスト
func TestDetectEncodingUTF16LEBOM(t *testing.T) {
	byteArray := []byte{0xFF, 0xFE, 0x53, 0x30, 0x93, 0x30, 0x6B, 0x30, 0x61, 0x30, 0x6F, 0x30}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF16LEBOM, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// UTF16(BOM無し) のエンコーディング判別と変換のテスト
func TestDetectEncodingUTF16BE(t *testing.T) {
	byteArray := []byte{0x30, 0x53, 0x30, 0x93, 0x30, 0x6B, 0x30, 0x61, 0x30, 0x6F}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF16BE, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}

// UTF16(BOM有り) のエンコーディング判別と変換のテスト
func TestDetectEncodingUTF16BEBOM(t *testing.T) {
	byteArray := []byte{0xFE, 0xFF, 0x30, 0x53, 0x30, 0x93, 0x30, 0x6B, 0x30, 0x61, 0x30, 0x6F}

	enc, err := jpdec.DetectEncoding(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, jpdec.EncodingUTF16BEBOM, enc)

	str, err := jpdec.Decode(byteArray)
	assert.NoError(t, err)
	assert.Equal(t, "こんにちは", str)
}
