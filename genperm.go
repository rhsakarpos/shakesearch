package main

import "strings"

func generateCasePerms(S string) []string {
	res := make([]string, 0, 1<<uint(len(S))+2)
	S = strings.ToLower(S)
	for k, v := range S {
		if isLetter(byte(v)) {
			switch len(res) {
			case 0:
				res = append(res, S, toUpper(S, k))
			default:
				for _, s := range res {
					res = append(res, toUpper(s, k))
				}
			}
		}
	}
	if len(res) == 0 {
		res = append(res, S)
	}
	res = append(res, strings.ToLower(S))
	res = append(res, strings.ToUpper(S))

	return res
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func toUpper(s string, i int) string {
	b := []byte(s)
	b[i] -= 'a' - 'A'
	return string(b)
}
