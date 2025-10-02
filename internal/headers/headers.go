package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"
const validChars = "^[A-Za-z0-9!#$%&'*+\\-\\.\\^_`|~]+$"

var re = regexp.MustCompile(validChars)

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)

	if !strings.Contains(str, crlf) {
		return 0, false, nil
	}

	if strings.HasPrefix(str, crlf) {
		return 2, true, nil
	}

	i := strings.Index(str, crlf)
	line := str[:i]
	c := strings.IndexByte(line, ':')

	if c == -1 {
		return 0, false, errors.New("invalid colon in the field line")
	}

	if c > 0 && line[c-1] == ' ' {
		return 0, false, errors.New("invalid space before colon and after the field-name")
	}

	key := strings.TrimSpace(line[:c])
	value := strings.TrimSpace(line[c+1:])

	if key == "" {
		return 0, false, errors.New("invalid key")
	}

	if !re.MatchString(key) {
		return 0, false, errors.New("invalid character/s in the key")
	}

	lowerKey := strings.ToLower(key)

	h.Set(lowerKey, value)

	return i + 2, false, nil
}

func (h Headers) Set(key, value string) {
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}
