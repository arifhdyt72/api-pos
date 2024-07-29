package tools

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseUint(val string) uint {
	eval, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return uint(0)
	}

	return uint(eval)
}

func ParseInt(val string) int {
	eval, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0
	}

	return int(eval)
}

func ParseUintP(val string) *uint {
	eval, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return nil
	}
	p := uint(eval)
	return &p
}

func ParseIntP(val string) *int {
	eval, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return nil
	}
	p := int(eval)
	return &p
}

func ToString(v interface{}) string {
	switch value := v.(type) {
	case string:
		return value
	case int:
		i := strconv.Itoa(value)
		return i
	case int64:
		i := strconv.FormatInt(value, 10)
		return i
	case uint:
		i := strconv.FormatUint(uint64(value), 10)
		return i
	case *uint:
		i := strconv.FormatUint(uint64(*value), 10)
		return i
	case *int:
		i := strconv.FormatInt(int64(*value), 10)
		return i
	case *int64:
		i := strconv.FormatInt(*value, 10)
		return i
	default:
		return ""
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func SaveImageToDisk(path string, filename string) (string, error) {
	idx := strings.Index(path, ";base64,")
	if idx < 0 {
		return "error", errors.New("invalid base64")
	}
	ImageType := path[11:idx]
	unbased, err := base64.StdEncoding.DecodeString(path[idx+8:])
	if err != nil {
		return "error", err
	}
	r := bytes.NewReader(unbased)

	fullname := fmt.Sprintf("images/icon/%s-%s.%s", time.Now().Format("02-01-2006"), filename, ImageType)

	switch ImageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			return "error", err
		}

		f, err := os.OpenFile(fullname, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "error", err
		}

		png.Encode(f, im)
	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			return "error", err
		}

		f, err := os.OpenFile(fullname, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "error", err
		}

		jpeg.Encode(f, im, nil)
	}

	return fullname, nil
}
