package jpdec

import (
	"errors"
	"sort"
)

// エンコーディング
const (
	// 不明
	EncodingUnknown = ""

	// Shift_JIS
	EncodingShiftJIS = "Shift_JIS"

	// JIS
	EncodingJIS = "ISO-2022-JP"

	// EUC
	EncodingEUCJP = "EUC-JP"

	// UTF-8
	EncodingUTF8 = "UTF-8"

	// UTF-8(BOM付)
	EncodingUTF8BOM = "UTF-8BOM"

	// UTF-16(BOM無し：リトルエンディアン)
	EncodingUTF16LE = "UTF-16LE"

	// UTF-16(BOM無し：リトルエンディアン)
	EncodingUTF16BE = "UTF-16BE"

	// UTF-16(BOM有り：リトルエンディアン)
	EncodingUTF16LEBOM = "UTF-16LEBOM"

	// UTF-16(BOM有り：ビッグエンディアン)
	EncodingUTF16BEBOM = "UTF-16BEBOM"
)

// エラー定義
var (
	// エンコーディングの検出に失敗
	ErrUnknown = errors.New("unknown encoding")

	// バイナリファイルっぽい
	ErrPossiblyBinary = errors.New("possibly binary")
)

// エンコーディングの検出
//
// エンコーディングを示す文字列である EncodingXxx 定数のいずれかを返します。
// 検出に失敗した場合は EncodingUnknown と ErrUnknown/ErrPossiblyBinary エラーを返します。
func DetectEncoding(bytes []byte) (string, error) {
	if isUTF8BOM(bytes) {
		return EncodingUTF8BOM, nil
	} else if isUnicodeBOM(bytes) {
		return EncodingUTF16LEBOM, nil
	} else if isUnicodeBEBOM(bytes) {
		return EncodingUTF16BEBOM, nil
	}

	mayBeJIS := true
	mayBeAscii := true
	mayBeUTF16 := false
	utf16LEScore := 0
	utf16BEScore := 0

	for i := 0; i < len(bytes); i++ {
		b1 := bytes[i]
		b2 := byte(0)
		if i+1 < len(bytes) {
			b2 = bytes[i+1]
		}

		// ISO-2022-JP (0x00～0x7F) の範囲外の場合はASCIIでもJISでもない
		if b1 >= 0x80 {
			mayBeAscii = false
			mayBeJIS = false
		}

		// JIS 判定
		if mayBeJIS && isJISEsc(bytes, i) {
			return EncodingJIS, nil
		}

		// UNICODE 判定
		if b1 == 0x00 {
			mayBeUTF16 = true
			mayBeAscii = false
			mayBeJIS = false

			// 1 つ後ろの文字または前の文字が 00 ～ 7F の間にある場合は UTF-16
			if b2 <= 0x7F || (i > 0 && bytes[i-1] <= 0x7F) {
				if i%2 == 0 {
					utf16BEScore++
				} else {
					utf16LEScore++
				}
			} else {
				return EncodingUnknown, ErrPossiblyBinary
			}
		} else if b1 == 0x30 {
			// 0x30 XX 0x30 XX が続くような形
			if i+2 < len(bytes) && bytes[i+1] != 0x30 && bytes[i+2] == 0x30 {
				if i%2 == 0 {
					utf16BEScore++
				} else {
					utf16LEScore++
				}
			}
		}
	}

	// 0x00 を吹く場合は UTF16LE, BE を優先
	if mayBeUTF16 {
		if utf16LEScore > utf16BEScore {
			return EncodingUTF16LE, nil
		} else if utf16BEScore > utf16LEScore {
			return EncodingUTF16BE, nil
		}
	}

	// ASCIIファイルの場合は UTF-8 を返す
	if mayBeAscii {
		return EncodingUTF8, nil
	}

	// 各エンコーディングのスコアを計算
	sjisScore := scoreShiftJIS(bytes)
	eucScore := scoreEUC(bytes)
	utf8Score := scoreUTF8(bytes)

	scoreByencoding := map[string]int{
		EncodingShiftJIS: sjisScore,
		EncodingEUCJP:    eucScore,
		EncodingUTF8:     utf8Score,
		EncodingUTF16BE:  utf16BEScore,
		EncodingUTF16LE:  utf16LEScore,
	}

	keys := make([]string, 0, len(scoreByencoding))
	for k := range scoreByencoding {
		keys = append(keys, k)
	}

	// 降順ソート
	sort.Slice(keys, func(i, j int) bool {
		return scoreByencoding[keys[i]] > scoreByencoding[keys[j]]
	})

	if scoreByencoding[keys[0]] > 0 {
		return keys[0], nil
	}

	// いずれにも該当しない場合はエラーを返す
	return EncodingUnknown, ErrUnknown
}

// b が from 以上 to 以下の場合に true を返す。
func between(b, from, to byte) bool {
	return from <= b && b <= to
}

// bytes の先頭3バイトが UTF-8 の BOM の場合に true を返す。
func isUTF8BOM(bytes []byte) bool {
	return len(bytes) > 3 && bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF
}

// bytes の先頭2バイトが UTF-16 の BOM の場合に true を返す。
func isUnicodeBOM(bytes []byte) bool {
	return len(bytes) > 2 && bytes[0] == 0xFF && bytes[1] == 0xFE
}

// bytes の先頭2バイトが UTF-16 の BOM(BE) の場合に true を返す。
func isUnicodeBEBOM(bytes []byte) bool {
	return len(bytes) > 2 && bytes[0] == 0xFE && bytes[1] == 0xFF
}

// bytes の index 番目から始まる文字列が JIS のエスケープ文字の場合に true を返す。
func isJISEsc(bytes []byte, index int) bool {
	if bytes[index] == 0x1B { // ESC
		if index+6 < len(bytes) {
			if bytes[index+1] == '&' && bytes[index+2] == '@' && bytes[index+3] == 0x1B &&
				bytes[index+4] == '$' && bytes[index+5] == 'B' {
				return true
			}
		}

		if index+3 < len(bytes) {
			if bytes[index+1] == '(' {
				return bytes[index+2] == 'B' || bytes[index+2] == 'J' || bytes[index+2] == 'I'
			} else if bytes[index+1] == '$' {
				if bytes[index+2] == '@' || bytes[index+2] == 'B' {
					return true
				}
				return index+4 < len(bytes) && bytes[index+2] == '(' && bytes[index+3] == 'D'
			}
		}
	}

	return false
}

// bytes バイト列に対して EUC の度合を示すスコアを計算する
func scoreEUC(bytes []byte) int {
	score := 0

	for i := 0; i < len(bytes); i++ {
		b1 := bytes[i]
		var b2, b3 byte
		if i+1 < len(bytes) {
			b2 = bytes[i+1]
		}
		if i+2 < len(bytes) {
			b3 = bytes[i+2]
		}

		if b1 >= 0x80 {
			// 漢字: 第1バイト・第2バイトとも0xA1～0xFE
			if between(b1, 0xA1, 0xFE) && between(b2, 0xA1, 0xFE) {
				score += 2
				i++
			} else if b1 == 0x8E && between(b2, 0xA1, 0xDF) {
				// 半角カタカナ: 第1バイト 0x8E, 第2バイト 0xA1～0xDF
				score += 2
				i++
			} else if b1 == 0x8F && between(b2, 0xA1, 0xFE) && between(b3, 0xA1, 0xFE) {
				// 補助漢字: 第1バイト 0x8F, 第2バイト・第3バイトとも 0xA1～0xFE
				score += 3
				i += 2
			} else {
				return 0
			}
		}
	}

	return score
}

// bytes バイト列に対して Shift_JIS の度合を示すスコアを計算する
func scoreShiftJIS(bytes []byte) int {
	score := 0

	for i := 0; i < len(bytes); i++ {
		b1 := bytes[i]
		var b2 byte
		if i+1 < len(bytes) {
			b2 = bytes[i+1]
		}

		if b1 >= 0x80 {
			// 全角文字
			if between(b1, 0x81, 0x9F) || between(b1, 0xE0, 0xFC) {
				if between(b2, 0x40, 0x7E) || between(b2, 0x80, 0xFC) {
					score += 2
					i++
				} else {
					// Shift_JIS ではない
					return 0
				}
			} else if !between(b1, 0xA1, 0xDF) {
				// 半角カタカナ(0xA1～0xDF)でも無い場合
				return 0
			}
		}
	}

	return score
}

// bytes バイト列に対して UTF-8 の度合を示すスコアを計算する
func scoreUTF8(bytes []byte) int {
	score := 0

	for i := 0; i < len(bytes); i++ {
		b1 := bytes[i]

		if b1 >= 0x80 {
			if b1 < 0xC0 {
				// UTF-8 の第1バイトは 0xC0 ～ 始まる
				return 0
			}

			// 頭 4bit の 1 の数が1文字のバイト数を表す
			var length int
			switch b1 >> 4 {
			case 0xF:
				length = 4
			case 0xE:
				length = 3
			case 0xC, 0xD:
				length = 2
			default:
				return 0
			}

			// 2バイト目以降の範囲は 0x80 から 0xBF の範囲
			for j := 1; j < length; j++ {
				if i+j >= len(bytes) || !between(bytes[i+j], 0x80, 0xBF) {
					return 0
				}
			}

			score += length
			i += length - 1
		}
	}

	return score
}
