package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

var crlf = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Header is not in correct form")
	}
	name := parts[0]
	val := bytes.TrimSpace(parts[1])

	//fieldname cannot have whitespace before semicolon
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("fieldname has whitespace")
	}
	strname := string(name)
	strname = strings.ToLower(strname)
	for _, c := range strname {
		if unicode.IsDigit(c) || unicode.IsLetter(c) {
			continue
		}
		switch c {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return "", "", fmt.Errorf("invalid character")

		}
	}
	return strname, string(val), nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	read := 0
	done = false
	for {
		idx := bytes.Index(data[read:], crlf)
		if idx == -1 {
			break
		}
		//empty header
		if idx == 0 {
			done = true
			read += len(crlf)
			return read, done, nil
		}
		name, val, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(crlf)
		value, ok := h[name]
		if ok {
			v := value + ", " + val
			h[name] = v
		} else {
			h[name] = val
		}
	}
	return read, done, nil

}
