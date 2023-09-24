package ingrid

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

func Map(mapping Mapfn, scanner *bufio.Scanner) error {
	var lineno int
	var currentSection []byte
	for scanner.Scan() {
		lineno++
		buf := scanner.Bytes()
		buf = bytes.TrimSpace(buf)

		var section, key, value, comment []byte
		var lbrack, rbrack, equal, hash = -1, -1, -1, -1

		for i, b := range buf {
			switch b {
			case '[':
				lbrack = i

			case ']':
				rbrack = i

			case '=':
				equal = i

			case '#':
				hash = i
				break
			}
		}

		if rbrack > lbrack && lbrack >= 0 && rbrack >= 0 {
			section = buf[lbrack+1 : rbrack]
			section = bytes.TrimSpace(section)
			currentSection = section
		}
		if equal >= 0 {
			key = bytes.TrimSpace(buf[:equal])
			if bytes.ContainsAny(key, " ") {
				return fmt.Errorf("syntax error: line %v: %s", lineno, string(buf))
			}

			value = buf[equal+1:]
			value = bytes.TrimSpace(value)
			if len(value) > 0 && bytes.ContainsAny(value[:1], "\"'`") {
				valstr, err := strconv.Unquote(string(value))
				if err != nil {
					return fmt.Errorf("syntax error: line %v: %s", lineno, string(buf))
				}
				value = []byte(valstr)
			}
		}
		if hash == 0 {
			comment = buf[hash:]
		}

		switch {
		case len(section)+len(key)+len(value)+len(comment) > 0:
			mapping(
				string(currentSection),
				string(key),
				string(value),
				string(comment),
			)

		case len(buf) > 0:
			return fmt.Errorf("syntax error: line %v: %s", lineno, string(buf))
		}

	}
	return nil
}

type Mapfn func(section, key, value, comment string) error