package goparsephone

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	TypeAll          = 0
	TypeMobileOnly   = 1
	TypeLandlineOnly = 2
)

type Phone struct {
	Raw       string
	Formatted string
	Carrier   string
	StartsAt  int
	EndsAt    int
}

func Pattern(findType int) string {
	var pieces []string

	switch findType {
	case TypeAll:
		pieces = append(MobileCarrierNumbers(), LandlineCarrierNumbers()...)
	case TypeMobileOnly:
		pieces = MobileCarrierNumbers()
	case TypeLandlineOnly:
		pieces = LandlineCarrierNumbers()
	}

	return `(\+)?(00)?(84)?(0)?(` + strings.Join(pieces, "|") + `)(\d{7})`
}

func RemoveAllSeparatorsAndSavePositions(text string) (string, []int) {
	var regex = regexp.MustCompile(`[ .\-]+`)

	var replaced = 0
	var replacedText = text

	var matches = regex.FindAllStringIndex(text, -1)
	var originOffsetIndex []int

	for i, match := range matches {
		// Remove separator.
		replacedText = replacedText[:match[0]-replaced] + replacedText[match[0]-replaced+1:]

		// Save separator position.
		originOffsetIndex = append(originOffsetIndex, match[0]-i)

		replaced++
	}

	return replacedText, originOffsetIndex
}

func FindInText(text string, findType int) []Phone {
	var result []Phone

	if text == "" || len(text) < 9 {
		return result
	}

	textWithoutSeparators, positionIndexes := RemoveAllSeparatorsAndSavePositions(text)

	if textWithoutSeparators == "" || len(textWithoutSeparators) < 9 {
		return result
	}

	regex := regexp.MustCompile(Pattern(findType))
	matches := regex.FindAllStringSubmatchIndex(textWithoutSeparators, -1)

	for _, subMatch := range matches {
		match := []int{subMatch[0], subMatch[1]}
		matched := textWithoutSeparators[match[0]:][:match[1]-match[0]]

		// Get position of phone number.
		m := ArrayFilter(positionIndexes, func(i int) bool {
			return i < (match[0] + len(matched))
		})

		start := len(ArrayFilter(m, func(i int) bool {
			return i <= match[0]
		})) + match[0]

		end := match[0] + len(m) + len(matched)

		// Get origin string.
		origin := text[start:end]

		// Filter phone number.
		if start > 0 && (end+1) < len(text) {
			if IsNumeric(text[start-1 : start]) {
				continue
			}

			check := text[end : end+1]

			if IsNumeric(check) && check != "+" {
				continue
			}

			regex = regexp.MustCompile(`[ .\-]`)

			if len(regex.FindAllString(origin, -1)) > 4 {
				continue
			}
		}

		// Get carrier.
		carrierNumber, _ := strconv.Atoi(textWithoutSeparators[subMatch[10]:subMatch[12]])
		carrier := Carriers()[carrierNumber]

		// Get national format of phone number.
		national := "+84" + textWithoutSeparators[subMatch[10]:subMatch[13]]

		// Export to structs.
		result = append(result, Phone{
			Raw:       origin,
			Formatted: national,
			Carrier:   carrier,
			StartsAt:  start,
			EndsAt:    end,
		})
	}

	return result
}
