package addressparser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

const fuzzyScore = 10

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func mapToSlice(mapData map[string]string) []string {
	var output []string
	for key, value := range mapData {
		output = append(output, key, value)
	}
	return output
}

func mapKeysToSlice(mapData map[string]string) []string {
	var output []string

	for k := range mapData {
		output = append(output, k)
	}
	return output
}

func getMatchKey(item string, mapData map[string]string) string {
	for key, value := range mapData {
		if key == item || value == item {
			return key
		}
	}
	return ""
}

func fuzzyMatch(str string, mapList map[string]string) []string {
	var output []string
	if str == "" {
		return output
	}

	list := mapToSlice(mapList)

	// check exact match
	if stringInSlice(str, list) {
		output = append(output, str)
		return output
	}
	// return if no exact macth for single char string
	if len(str) == 1 {
		return output
	}
	// no exact match do fuzzy
	fuzzyResults := fuzzy.RankFind(str, list)
	if len(fuzzyResults) == 0 {
		return output
	}
	// sort the results
	sort.Sort(fuzzyResults)

	if fuzzyResults[0].Distance >= fuzzyScore {
		// if distance is not within limits, check if partial result.
		if !strings.HasSuffix(fuzzyResults[0].Target, fuzzyResults[0].Source) {
			// no results found
			// return empty
			return output
		}
	}
	// return first item of results.
	output = append(output, fuzzyResults[0].Target)
	return output
}

func splitNumberRange(addressPart string) (int, int) {
	var start int
	var end int
	checkFileRegex, err := regexp.Compile("^([0-9]+)-([0-9]+)$")
	if err != nil {
		return start, end
	}

	matched := checkFileRegex.FindStringSubmatch(addressPart)
	if matched == nil {
		return start, end
	}
	start, _ = strconv.Atoi(matched[1])
	end, _ = strconv.Atoi(matched[2])
	return start, end
}

func splitMixedIndex(addressPart string) (int, string) {
	var foundString string
	var foundInt int
	checkFileRegex, err := regexp.Compile("^([0-9]+)([A-Z]{1,2})$")
	if err != nil {
		return foundInt, foundString
	}

	matched := checkFileRegex.FindStringSubmatch(addressPart)
	if matched == nil {
		return foundInt, foundString
	}
	foundInt, _ = strconv.Atoi(matched[1])
	foundString = matched[2]
	return foundInt, foundString
}

// AddressParts - parts fo address found
type AddressParts struct {
	AddressString      string
	AddressStringParts []string
	// parts of the address
	LevelType          string
	LevelNumber        int
	FlatType           string
	FlatNumber         int
	FlatNumberSuffix   string
	StreetNumber       int
	StreetNumberEnd    int
	StreetNumberSuffix string
	StreetName         string
	StreetType         string
	StreetSuffix       string
	Suburb             string
	PostCode           int
	State              string
}

// NewAddress parse an address string into address struct
func NewAddress(address string) (*AddressParts, error) {
	var err error
	addressParts := new(AddressParts)
	err = addressParts.LoadAddressString(address)
	if err != nil {
		return addressParts, err
	}
	addressParts.ProcessAddress()

	return addressParts, err
}

// LoadAddressString load a address string
func (ap *AddressParts) LoadAddressString(addressString string) error {
	var err error
	reg, err := regexp.Compile("[^A-Za-z0-9 -]+")
	if err != nil {
		return errors.Wrap(err, "Failed to compile regex")
	}

	addressString = strings.ToUpper(reg.ReplaceAllString(addressString, ""))
	ap.AddressString = addressString
	ap.AddressStringParts = strings.Split(addressString, " ")
	if len(ap.AddressStringParts) < 3 {
		return errors.New("Address string too short")
	}

	return err
}

// ProcessAddress process the address parts
func (ap *AddressParts) ProcessAddress() {
	var foundIndex int

	// first part is string, most level or flat type
	if ap.isPartString(0) {

		// look for flat types
		foundIndex = ap.findIndex(1, ap.isPartAnyNumber, ap.isPartString)
		if foundIndex != 0 {
			matchedFlatType, matchedIndex := ap.matchAddressPart(foundIndex-1, flatTypes)
			if matchedFlatType != "" {
				ap.FlatType = matchedFlatType
				if !ap.addressPartNoNumber(matchedFlatType) {
					ap.FlatNumber, ap.FlatNumberSuffix, _ = ap.getAddressNumber(foundIndex)
					ap.removeParts(foundIndex)
				}
				ap.removeParts(matchedIndex...)
			}
		}

		// look for level types
		foundIndex = ap.findIndex(1, ap.isPartAnyNumber, ap.isPartString)
		if foundIndex != 0 {
			matchedLevelType, matchedIndex := ap.matchAddressPart(foundIndex-1, levelTypes)
			if matchedLevelType != "" {
				ap.LevelType = matchedLevelType
				if !ap.addressPartNoNumber(matchedLevelType) {
					ap.LevelNumber, _, _ = ap.getAddressNumber(foundIndex)
					ap.removeParts(foundIndex)
				}
				ap.removeParts(matchedIndex...)
			}
		}

	}

	// if last part is a number it is prossibly a postcode
	lastIndex := len(ap.AddressStringParts) - 1
	if ap.isPartNumber(lastIndex) {
		postCode := ap.matchPostCode(lastIndex)
		if postCode != 0 {
			ap.PostCode = postCode
			ap.removeParts(lastIndex)
		}
	}

	// look for the state - usually the last string on the address
	foundIndex = ap.findIndexReverse(lastIndex, ap.isPartString, ap.isPartAny)
	if foundIndex != 0 {
		matchedState, matchedIndex := ap.matchAddressPart(foundIndex, australianStates)

		if matchedState != "" {
			ap.State = matchedState
			ap.removeParts(matchedIndex...)
		}
	}

	// look for the state - usually the last string on the address followed by a string
	foundIndex = ap.findIndex(1, ap.isPartString, ap.isPartAnyNumber)
	if foundIndex != 0 {
		ap.StreetNumber,
			ap.StreetNumberSuffix,
			ap.StreetNumberEnd = ap.getAddressNumber(foundIndex - 1)
		ap.removeParts(foundIndex - 1)
	}

	// find street type
	foundIndex = ap.findIndex(1, ap.isPartStreetType, ap.isPartString)
	if foundIndex != 0 {
		if (foundIndex-1 >= 0) && ap.AddressStringParts[foundIndex-1] == "THE" {
			ap.StreetName = fmt.Sprintf(
				"%s %s",
				ap.AddressStringParts[foundIndex-1],
				ap.AddressStringParts[foundIndex],
			)
			ap.StreetType = "-"
			ap.removeParts(foundIndex-1, foundIndex)
		} else {
			ap.StreetType = getMatchKey(ap.AddressStringParts[foundIndex], streetTypes)
			ap.removeParts(foundIndex)
		}
	}

	// before street type should be street Name
	//
	if foundIndex > 0 && ap.StreetName == "" {
		var matchedIndex []int
		ap.StreetName, matchedIndex = ap.getStringBefore(foundIndex)
		ap.removeParts(matchedIndex...)
	} else if ap.StreetType == "" && ap.StreetName == "" && ap.isPartString(lastIndex) {
		// street name/type not found, last string might be the street name
		ap.StreetName = ap.AddressStringParts[lastIndex]
	}

	// after street type should be suburb
	nextIndex := foundIndex + 1
	if foundIndex > 0 && ap.hasPartIndex(nextIndex) {
		// look for street suffix
		streetSuffixesKeys := mapKeysToSlice(streetSuffixes)
		if stringInSlice(ap.AddressStringParts[nextIndex], streetSuffixesKeys) {
			ap.StreetSuffix = ap.AddressStringParts[nextIndex]
			ap.removeParts(nextIndex)
			foundIndex = nextIndex
		}
		var matchedIndex []int
		ap.Suburb, matchedIndex = ap.getStringAfter(foundIndex)
		ap.removeParts(matchedIndex...)
	}

}

func (ap *AddressParts) hasPartIndex(index int) bool {
	lastIndex := len(ap.AddressStringParts) - 1
	if index >= 0 && index <= lastIndex {
		return true
	}
	return false
}

func (ap *AddressParts) removeParts(indexes ...int) {
	for _, index := range indexes {
		ap.AddressStringParts[index] = ""
	}
}

func (ap *AddressParts) addressPartNoNumber(addressPart string) bool {
	addressTypesNoNumberSlice := mapToSlice(addressTypesNoNumber)
	return stringInSlice(addressPart, addressTypesNoNumberSlice)
}

func (ap *AddressParts) matchAddressPart(index int, addressTypes map[string]string) (string, []int) {

	var matchedPart string
	var matchedIndex []int
	var matchResult []string

	for i := index; i >= 0; i-- {
		var currentPart string
		if matchedPart == "" {
			currentPart = ap.AddressStringParts[i]
		} else {
			currentPart = fmt.Sprintf("%s %s", ap.AddressStringParts[i], matchedPart)
		}

		result := fuzzyMatch(currentPart, addressTypes)

		if result == nil {
			// no match - quit loop
			break
		}
		matchResult = result
		matchedIndex = append(matchedIndex, i)

		matchedPart = currentPart
	}

	if len(matchResult) == 1 {
		matchedPart = matchResult[0]
	}

	// switch the matched part to the key (code) value
	matchedPart = getMatchKey(matchedPart, addressTypes)

	return matchedPart, matchedIndex
}

func (ap *AddressParts) getStringBefore(startIndex int) (string, []int) {
	var matchedString string
	var matchedIndex []int

	for i := (startIndex - 1); i >= 0; i-- {
		addressPart := ap.AddressStringParts[i]

		if !ap.isPartString(i) {
			return matchedString, matchedIndex
		}
		if matchedString == "" {
			matchedString = addressPart
		} else {
			matchedString = fmt.Sprintf("%s %s", addressPart, matchedString)
		}
		matchedIndex = append(matchedIndex, i)
	}

	return matchedString, matchedIndex
}

func (ap *AddressParts) getStringAfter(startIndex int) (string, []int) {
	var matchedString string
	var matchedIndex []int

	startIndex = startIndex + 1
	for i := startIndex; i < len(ap.AddressStringParts); i++ {
		addressPart := ap.AddressStringParts[i]

		if !ap.isPartString(i) {
			return matchedString, matchedIndex
		}
		if matchedString == "" {
			matchedString = addressPart
		} else {
			matchedString = fmt.Sprintf("%s %s", matchedString, addressPart)
		}
		matchedIndex = append(matchedIndex, i)
	}

	return matchedString, matchedIndex
}

func (ap *AddressParts) getAddressNumber(index int) (int, string, int) {
	var addressNumber int
	var addressNumberSuffix string
	var addressNumberEnd int

	if ap.isPartMixed(index) {
		// set address type number and suffix
		addressNumber, addressNumberSuffix = splitMixedIndex(ap.AddressStringParts[index])
	} else if ap.isPartNumberRange(index) {
		addressNumber, addressNumberEnd = splitNumberRange(ap.AddressStringParts[index])
	} else {
		// set address type number
		addressNumber, _ = strconv.Atoi(ap.AddressStringParts[index])
	}
	// clear address type number
	return addressNumber, addressNumberSuffix, addressNumberEnd
}

func (ap *AddressParts) isPartString(index int) bool {
	if index < 0 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z]+$", ap.AddressStringParts[index])
	return matched
}

func (ap *AddressParts) isPartNumber(index int) bool {
	matched, _ := regexp.MatchString("^[0-9]+$", ap.AddressStringParts[index])
	return matched
}

func (ap *AddressParts) isPartNumberRange(index int) bool {
	matched, _ := regexp.MatchString("^([0-9]+)-([0-9]+)$", ap.AddressStringParts[index])
	return matched
}

func (ap *AddressParts) isPartMixed(index int) bool {
	matched, _ := regexp.MatchString("^([0-9]+)([A-Z]{1,2})$", ap.AddressStringParts[index])
	return matched
}

func (ap *AddressParts) isPartNumberOrMixed(index int) bool {
	return (ap.isPartMixed(index) || ap.isPartNumber(index))
}

func (ap *AddressParts) isPartAnyNumber(index int) bool {
	return (ap.isPartMixed(index) || ap.isPartNumber(index) || ap.isPartNumberRange(index))
}

func (ap *AddressParts) isPartAny(_ int) bool {
	return true
}

func (ap *AddressParts) isPartStreetType(index int) bool {
	if index < 0 {
		return false
	}
	result := fuzzyMatch(ap.AddressStringParts[index], streetTypes)
	return (len(result) > 0)
}

// check if index mathces postcode
func (ap *AddressParts) matchPostCode(index int) int {
	var postCode int

	if matched, _ := regexp.MatchString("^[0-9]{4}$", ap.AddressStringParts[index]); matched {
		// found potential postcode
		postCode, _ = strconv.Atoi(ap.AddressStringParts[index])
		return postCode
	}
	// not found
	return postCode
}

func (ap *AddressParts) findIndex(
	startIndex int, current func(int) bool, last func(int) bool,
) int {
	for index := startIndex; index < len(ap.AddressStringParts); index++ {
		if current(index) && last(index-1) {
			return index
		}
	}
	return 0
}

func (ap *AddressParts) findIndexReverse(
	startIndex int, current func(int) bool, last func(int) bool,
) int {
	for index := startIndex; index >= 0; index-- {
		if current(index) && last(index+1) {
			return index
		}
	}
	return 0
}
