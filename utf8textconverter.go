package utilities

import (
	"bytes"
	"io"
	"log"

	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func ConvertToUTF8(strBytes []byte) ([]byte, error) {
	detector := chardet.NewTextDetector()
	detectorResult, err := detector.DetectBest(strBytes)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(detectorResult.Charset, byteReader)
	strBytes, _ = io.ReadAll(reader)
	return strBytes, nil
}

// BytesToStr converts a []byte representation of a plaintext file to string
func BytesToStr(data []byte) (string, error) {

	utf8Bytes, err := ConvertToUTF8(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(utf8Bytes), nil
}

func ConvertStringToUTF8(content string) (string, error) {
	return BytesToStr([]byte(content))
}
