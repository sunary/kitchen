package j

import (
	"strings"
)

func MarshalJson(body map[string]string) string {
	sb := strings.Builder{}
	sb.WriteString("{")

	i := 0
	for k, v := range body {
		sb.WriteString("\"")
		sb.WriteString(k)
		sb.WriteString("\":\"")
		v = strings.Replace(v, "\"", "\\\"", -1)
		sb.WriteString(v)
		sb.WriteString("\"")

		if i != len(body)-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String()
}

func UnmarshalJson(message string) map[string]string {
	m := map[string]string{}
	for i := 0; i < len(message)-2; {
		i += readChar(message, '"', i+1)
		key := readString(message, i+1)
		i += len(key) - 1
		i += 1 // "
		i += readChar(message, ':', i+1)
		i += readChar(message, '"', i+1)
		value := readString(message, i+1)
		i += len(value) - 1
		i += 1 // "
		m[key] = value
		i += 1 // ,
	}
	return m
}

func readChar(message string, char rune, start int) int {
	i := 0
	for ; ; i++ {
		if message[start+i] == uint8(char) && (char != '"' || message[start+i-1] != uint8('\\')) {
			break
		}
	}
	return i + 1
}

func readString(message string, start int) string {
	sb := strings.Builder{}
	for message[start] != uint8('"') || message[start-1] == uint8('\\') {
		sb.WriteByte(message[start])
		if start == len(message)-1 {
			break
		}
		start++
	}
	return sb.String()
}
