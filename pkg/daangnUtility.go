package pkg

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ParsePriceString(letter string) int {
	ret := ""
	for _, val := range letter {
		if unicode.IsDigit(val) {
			ret += string(val)
		}
	}
	retInt, _ := strconv.Atoi(ret)
	return retInt
}

// "관심 8 ∙채팅 2∙조회 158"
func ParseViewCountsString(letter string) (int, int, int) {
	likeCount, viewCount, chatCount := 0, 0, 0
	letter = strings.Replace(letter, " ", "", -1)
	letter = strings.Replace(letter, "\t", "", -1)
	letter = strings.Replace(letter, "\n", "", -1)
	tokens := strings.Split(letter, "∙")
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	for _, token := range tokens {
		if strings.Contains(token, "관심") {
			likeCount, _ = strconv.Atoi(re.FindString(token))
		} else if strings.Contains(token, "조회") {
			viewCount, _ = strconv.Atoi(re.FindString(token))
		} else if strings.Contains(token, "채팅") {
			chatCount, _ = strconv.Atoi(re.FindString(token))
		}
	}
	return likeCount, viewCount, chatCount
}

func ExtractIntFromString(input string) (int, error) {
	var result int
	var sb strings.Builder
	for _, char := range input {
		if char >= '0' && char <= '9' {
			sb.WriteRune(char)
		}
	}
	parsedInt, err := strconv.Atoi(sb.String())
	if err == nil {
		result = parsedInt
	}
	return result, err
}
