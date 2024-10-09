package jpdec

import (
	"io"
	"os"
)

// 指定されたファイルを読み込んで、文字列として返します。
func ReadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// ファイルの内容をすべて読み込む
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return Decode(bytes)
}
