English | [日本語](./README.jp.md)

# Overview

Go 言語用の日本語エンコーディングの自動判別とデコードを行うパッケージ

## 対応エンコーディング

- Shift_JIS
- EUC-JP
- JIS(ISO-2022-JP)
- UTF-8
- UTF-8(with BOM)
- UTF-16
- UTF-16(with BOM)
- UTF-16(Big Endian)
- UTF-16(Big Endian with BOM)

To detect encoding exactly, use follow encoding:

- JIS
- UTF-8(with BOM)
- UTF-16(with BOM)
- UTF-16(Big Endian with BOM)
- UTF-16(Little Endian containing ASCII char)
- UTF-16(Big Endian containing ASCII char)

Other encodings are not accurate, just pick the one that is more likely.

## サンプル

デコード:

```go
byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
str, _ := jpdec.Decode(byteArray)
fmt.Println(str)
```

output:

```
こんにちは
```

エンコーディングの判別:

```
byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
enc, _ := jpdec.DetectEncoding(byteArray)
fmt.Println(enc)
```

output:

```
Shift_JIS
```
