package pdf

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
)

func ExtractText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ING-001: read PDF file: %w", err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return "", fmt.Errorf("ING-003: invalid PDF header")
	}
	text := extractLiteralText(data)
	if text == "" {
		return "", fmt.Errorf("ING-005: PDF text layer is empty or unsupported")
	}
	return text, nil
}

func extractLiteralText(data []byte) string {
	var out []rune
	for i := 0; i < len(data); i++ {
		if data[i] != '(' {
			continue
		}
		i++
		var runes []rune
		escaped := false
		for ; i < len(data); i++ {
			ch := data[i]
			if escaped {
				switch ch {
				case 'n':
					runes = append(runes, '\n')
				case 'r':
					runes = append(runes, '\r')
				case 't':
					runes = append(runes, '\t')
				default:
					runes = append(runes, rune(ch))
				}
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == ')' {
				break
			}
			runes = append(runes, rune(ch))
		}
		if meaningfulText(runes) {
			out = append(out, runes...)
			out = append(out, '\n')
		}
	}
	return string(out)
}

func meaningfulText(runes []rune) bool {
	letters := 0
	for _, r := range runes {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			letters++
		}
	}
	return letters >= 2
}
