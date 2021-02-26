package goparsephone

import (
	"reflect"
	"testing"
)

var testString = "0989999888 09899998880989999888 09 89 99 98 88 0989 999 888      0989-999-888 0989.999.888  0084 868 606 701 +840868584147"

func TestRemoveAllSeparatorsAndSavePositions(t *testing.T) {
	wantString := "09899998880989999888098999988809899998880989999888     09899998880989999888 0084868606701+840868584147"
	wantPositions := []int{10, 30, 32, 34, 36, 38, 40, 44, 47, 50, 59, 62, 65, 69, 72, 75, 80, 83, 86, 89}

	result, positions := RemoveAllSeparatorsAndSavePositions(testString)

	if wantString != result || !reflect.DeepEqual(wantPositions, positions) {
		t.Fatalf(`RemoveAllSeparatorsAndSavePositions("%s") = %s, %v, want match for %s, %v`, testString, result, positions, wantString, wantPositions)
	}
}

func TestFindInText(t *testing.T) {
	result := FindInText(testString, TypeAll)

	wantResult := []Phone{
		{
			Raw:       "0989999888",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  0,
			EndsAt:    10,
		},
		{
			Raw:       "09 89 99 98 88",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  32,
			EndsAt:    46,
		},
		{
			Raw:       "0989 999 888",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  47,
			EndsAt:    59,
		},
		{
			Raw:       "0989-999-888",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  65,
			EndsAt:    77,
		},
		{
			Raw:       "0989.999.888",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  78,
			EndsAt:    90,
		},
		{
			Raw:       "0084 868 606 701",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  92,
			EndsAt:    108,
		},
		{
			Raw:       "+840868584147",
			Formatted: "+84989999888",
			Carrier:   "Viettel",
			StartsAt:  109,
			EndsAt:    122,
		},
	}

	if len(result) != len(wantResult) {
		t.Fatalf(`TestFindInText("%s") = %v, want match for %v`, testString, result, wantResult)
	}
}

func TestGetCarrier(t *testing.T) {
	text := "0978123456"
	want := Carriers()[97]
	wantE164 := "+84978123456"

	carrier, e164 := GetCarrier(text, "")

	if carrier != want || e164 != wantE164 {
		t.Fatalf(`TestGetCarrier("%s") = %s, %s, want match for %s, %s`, text, carrier, e164, want, wantE164)
	}
}
