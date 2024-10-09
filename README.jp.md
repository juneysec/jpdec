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

エンコーディングを正確に判別するには、下記のエンコーディングを使用してください：

- JIS
- UTF-8(BOM 有り)
- UTF-16(BOM 有り)
- UTF-16(ASCII 文字を含む)
- UTF-16(ビッグエンディアン BOM 有り)
- UTF-16(ビッグエンディアン ASCII 文字を含む)

その他のエンコーディング(Shift_JIS や EUC-JP, UTF-8 など)は、可能性の高いものが選択されます。

## サンプル

### デコード:

```go
byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
str, _ := jpdec.Decode(byteArray)
fmt.Println(str)
```

output:

```
こんにちは
```

### エンコーディングの判別:

```
byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
enc, _ := jpdec.DetectEncoding(byteArray)
fmt.Println(enc)
```

output:

```
Shift_JIS
```
