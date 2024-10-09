// Detecting and Decode for Japanese encoding
//
// # Supported Encoding
//
//   - Shift_JIS
//   - EUC-JP
//   - JIS(ISO-2022-JP)
//   - UTF-8
//   - UTF-8(with BOM)
//   - UTF-16
//   - UTF-16(with BOM)
//   - UTF-16(Big Endian)
//   - UTF-16(Big Endian with BOM)
//
// To detect encoding exactly, use follow encoding:
//   - JIS
//   - UTF-8(with BOM)
//   - UTF-16(with BOM)
//   - UTF-16(Big Endian with BOM)
//   - UTF-16(Little Endian containing ASCII char)
//   - UTF-16(Big Endian containing ASCII char)
//
// Other encodings are not accurate, just pick the one that is more likely.
package jpdec

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// エンコーディングの自動判別とデコード
//
// SJIS, EUC-JP などの文字列のバイト配列 byteArray を文字列に変換します。
// 失敗した場合、string(byteArray) と ErrUnknown/ErrPossiblyBinary エラーを返します。
func Decode(byteArray []byte) (string, error) {
	encoding, err := DetectEncoding(byteArray)
	if err != nil {
		return string(byteArray), err
	}

	switch encoding {
	// Shift_JIS
	case EncodingShiftJIS:
		return decode(byteArray, japanese.ShiftJIS.NewDecoder().Transformer)

	// JIS
	case EncodingJIS:
		return decode(byteArray, japanese.ISO2022JP.NewDecoder().Transformer)

	// EUC
	case EncodingEUCJP:
		return decode(byteArray, japanese.EUCJP.NewDecoder().Transformer)

	// UTF-8
	case EncodingUTF8:
		return string(byteArray), nil

	// UTF-8(BOM付)
	case EncodingUTF8BOM:
		return string(byteArray[3:]), nil

	// UTF-16(リトルエンディアン)
	case EncodingUTF16LE:
		return decode(byteArray, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder())

	// UTF-16(ビッグエンディアン)
	case EncodingUTF16BE:
		return decode(byteArray, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder())

	// UTF-16(BOM有り：リトルエンディアン)
	case EncodingUTF16LEBOM:
		return decode(byteArray[2:], unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder())

	// UTF-16(BOM有り：ビッグエンディアン)
	case EncodingUTF16BEBOM:
		return decode(byteArray[2:], unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder())

	}

	return string(byteArray), err
}

func decode(byteArray []byte, decoder transform.Transformer) (string, error) {
	reader := transform.NewReader(bytes.NewReader(byteArray), decoder)
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
