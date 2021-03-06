/**
 * Copyright (c) Ladifire, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

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
	Raw         string
	Formatted   string
	UnFormatted string
	Carrier     string
	StartsAt    int
	EndsAt      int
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

	return `(((0084|\+84)(0?))|0)(` + strings.Join(pieces, "|") + `)(\d{7})`
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

func GetCarrier(text string, carrierNumber string) (string, string, string) {
	var formatted string
	var raw string

	if text != "" && carrierNumber == "" {
		regex := regexp.MustCompile(Pattern(TypeAll))
		matches := regex.FindStringSubmatch(text)

		if matches == nil || len(matches) == 0 {
			return "", "", ""
		}

		carrierNumber = matches[5]
		raw = matches[5] + matches[6]
		formatted = "+84" + raw
	}

	n, _ := strconv.Atoi(carrierNumber)

	return Carriers()[n], formatted, "0" + raw
}

func FindInText(text string, findType int) []Phone {
	var result []Phone

	if text == "" || len(text) < 9 {
		return result
	}

	// Replace all UTF-8 to ASCII.
	u := regexp.MustCompile(`[^\x00-\x7F]`)
	text = u.ReplaceAllString(text, "a")

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

			if len(regex.FindAllString(origin, -1)) > 3 {
				continue
			}
		}

		// Get carrier.
		carrier, _, _ := GetCarrier("", textWithoutSeparators[subMatch[10]:subMatch[12]])

		// Get national format of phone number.
		national := "+84" + textWithoutSeparators[subMatch[10]:subMatch[13]]

		// Export to structs.
		result = append(result, Phone{
			Raw:         origin,
			Formatted:   national,
			UnFormatted: "0" + textWithoutSeparators[subMatch[10]:subMatch[13]],
			Carrier:     carrier,
			StartsAt:    start,
			EndsAt:      end,
		})
	}

	return result
}
