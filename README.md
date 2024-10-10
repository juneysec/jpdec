English | [日本語](./README.jp.md)

# Overview

Detecting and Decode for Japanese encoding for Go lang.

## Supported Encoding

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

Other encodings are not exact, and the one that is more likely to be chosen.

## Samples

### Decode:

```go
byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
str, _ := jpdec.Decode(byteArray)
fmt.Println(str)
```

output:

```
こんにちは
```

### Detect encoding:

```
byteArray := []byte{0x82, 0xB1, 0x82, 0xF1, 0x82, 0xC9, 0x82, 0xBF, 0x82, 0xCD}
enc, _ := jpdec.DetectEncoding(byteArray)
fmt.Println(enc)
```

output:

```
Shift_JIS
```
