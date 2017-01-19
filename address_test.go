package addressparser

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var testAddresses = []AddressParts{
	{
		AddressString:      "123 WESTLING HWY",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       123,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "WESTLING",
		StreetType:         "HIGHWAY",
		Suburb:             "",
		PostCode:           0,
		State:              "",
	},
	{
		AddressString:      "1910 LONG LONG HWY MT NOWHERE VIC 3216",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       1910,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "LONG LONG",
		StreetType:         "HIGHWAY",
		Suburb:             "MT NOWHERE",
		PostCode:           3216,
		State:              "VIC",
	},
	{
		AddressString:      "Level 11 500 Collins St Melbourne 3000",
		LevelType:          "L",
		LevelNumber:        11,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       500,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "COLLINS",
		StreetType:         "STREET",
		Suburb:             "MELBOURNE",
		PostCode:           3000,
		State:              "",
	},
	{
		AddressString:      "UNIT 3A 123 Butcher St St Fakeburb QLD 4568",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "UNIT",
		FlatNumber:         3,
		FlatNumberSuffix:   "A",
		StreetNumber:       123,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "BUTCHER",
		StreetType:         "STREET",
		Suburb:             "ST FAKEBURB",
		PostCode:           4568,
		State:              "QLD",
	},
	{
		AddressString:      "123 The Boulevarde, Flat Oak NSW 2529",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       123,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "THE BOULEVARDE",
		StreetType:         "-",
		Suburb:             "FLAT OAK",
		PostCode:           2529,
		State:              "NSW",
	},
	{
		AddressString:      "545 FLAT ROCK RD, GYMEA BAY NSW 2227",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       545,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "FLAT ROCK",
		StreetType:         "ROAD",
		Suburb:             "GYMEA BAY",
		PostCode:           2227,
		State:              "NSW",
	},
	{
		AddressString:      "UNIT 3 5-11 FLATHEAD RD, ETTALONG BEACH NSW 2257",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "UNIT",
		FlatNumber:         3,
		FlatNumberSuffix:   "",
		StreetNumber:       5,
		StreetNumberEnd:    11,
		StreetNumberSuffix: "",
		StreetName:         "FLATHEAD",
		StreetType:         "ROAD",
		Suburb:             "ETTALONG BEACH",
		PostCode:           2257,
		State:              "NSW",
	},
	{
		AddressString:      "LOT 374 WEST ST, ASCOT PARK SA 5043",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "LOT",
		FlatNumber:         374,
		FlatNumberSuffix:   "",
		StreetNumber:       0,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "WEST",
		StreetType:         "STREET",
		Suburb:             "ASCOT PARK",
		PostCode:           5043,
		State:              "SA",
	},
	{
		AddressString:      "123 HIGHWAY ST, PARK AVENUE QLD 4701",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       123,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "HIGHWAY",
		StreetType:         "STREET",
		Suburb:             "PARK AVENUE",
		PostCode:           4701,
		State:              "QLD",
	},
	{
		AddressString:      "583 ST AGNES CT, ST AGNES SA 5097",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       583,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "ST AGNES",
		StreetType:         "COURT",
		Suburb:             "ST AGNES",
		PostCode:           5097,
		State:              "SA",
	},
	{
		AddressString:      "168A MADEUP RD, MELBOURNE VIC 3000",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       168,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "A",
		StreetName:         "MADEUP",
		StreetType:         "ROAD",
		Suburb:             "MELBOURNE",
		PostCode:           3000,
		State:              "VIC",
	},
	{
		AddressString:      "2123-2127 Carlingford Rd, Carlingford NSW 2118",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       2123,
		StreetNumberEnd:    2127,
		StreetNumberSuffix: "",
		StreetName:         "CARLINGFORD",
		StreetType:         "ROAD",
		Suburb:             "CARLINGFORD",
		PostCode:           2118,
		State:              "NSW",
	},
	{
		AddressString:      "Lot 11, 75 Scotts Head Road, Way Way, NSW 2447",
		LevelType:          "",
		LevelNumber:        0,
		FlatType:           "LOT",
		FlatNumber:         11,
		FlatNumberSuffix:   "",
		StreetNumber:       75,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "SCOTTS HEAD",
		StreetType:         "ROAD",
		Suburb:             "WAY WAY",
		PostCode:           2447,
		State:              "NSW",
	},
	{
		AddressString:      "Lower Ground Floor 15 Phillip street Sydney NSW 2000",
		LevelType:          "LG",
		LevelNumber:        0,
		FlatType:           "",
		FlatNumber:         0,
		FlatNumberSuffix:   "",
		StreetNumber:       15,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "PHILLIP",
		StreetType:         "STREET",
		Suburb:             "SYDNEY",
		PostCode:           2000,
		State:              "NSW",
	},
	{
		AddressString:      "Suite 1 Level 14, 200 Queen St, Melbourne VIC 3000",
		LevelType:          "L",
		LevelNumber:        14,
		FlatType:           "SE",
		FlatNumber:         1,
		FlatNumberSuffix:   "",
		StreetNumber:       200,
		StreetNumberEnd:    0,
		StreetNumberSuffix: "",
		StreetName:         "QUEEN",
		StreetType:         "STREET",
		Suburb:             "MELBOURNE",
		PostCode:           3000,
		State:              "VIC",
	},
}

func errorExpectedString(t *testing.T, expected string, actual string) {
	t.Errorf("expected %s, actual %s", expected, actual)
}

func errorExpectedInt(t *testing.T, expected int, actual int) {
	t.Errorf("expected %d, actual %d", expected, actual)
}

func TestAddress(t *testing.T) {
	addressParts, err := NewAddress("Level 15 520 collins street 3000")
	if err != nil {
		t.Error(err)
	}
	// addressParts.LoadAddressString("2243-2247 Carlingford Rd, Carling Ford SA 2118")
	// addressParts.LoadAddressString("L 12 20-21 PARK AVE N, EILDON VIC 3713")
	//addressParts.LoadAddressString("UNIT 1 2 LITTLE HILL STREET TWEED HEADS NSW 2485")
	addressParts.ProcessAddress()

	spew.Dump(addressParts)
}

func TestAddresses(t *testing.T) {

	for _, testAddress := range testAddresses {
		addressParts := AddressParts{}
		addressParts.LoadAddressString(testAddress.AddressString)
		addressParts.ProcessAddress()

		var hasError bool
		if testAddress.LevelType != addressParts.LevelType {
			hasError = true
			errorExpectedString(t, testAddress.LevelType, addressParts.LevelType)
		}
		if testAddress.LevelNumber != addressParts.LevelNumber {
			hasError = true
			errorExpectedInt(t, testAddress.LevelNumber, addressParts.LevelNumber)
		}
		if testAddress.FlatType != addressParts.FlatType {
			hasError = true
			errorExpectedString(t, testAddress.FlatType, addressParts.FlatType)
		}
		if testAddress.FlatNumber != addressParts.FlatNumber {
			hasError = true
			errorExpectedInt(t, testAddress.FlatNumber, addressParts.FlatNumber)
		}
		if testAddress.FlatNumberSuffix != addressParts.FlatNumberSuffix {
			hasError = true
			errorExpectedString(t, testAddress.FlatNumberSuffix, addressParts.FlatNumberSuffix)
		}
		if testAddress.StreetNumber != addressParts.StreetNumber {
			hasError = true
			errorExpectedInt(t, testAddress.StreetNumber, addressParts.StreetNumber)
		}
		if testAddress.StreetNumberEnd != addressParts.StreetNumberEnd {
			hasError = true
			errorExpectedInt(t, testAddress.StreetNumberEnd, addressParts.StreetNumberEnd)
		}
		if testAddress.StreetNumberSuffix != addressParts.StreetNumberSuffix {
			hasError = true
			errorExpectedString(t, testAddress.StreetNumberSuffix, addressParts.StreetNumberSuffix)
		}
		if testAddress.StreetName != addressParts.StreetName {
			hasError = true
			errorExpectedString(t, testAddress.StreetName, addressParts.StreetName)
		}
		if testAddress.StreetType != addressParts.StreetType {
			hasError = true
			errorExpectedString(t, testAddress.StreetType, addressParts.StreetType)
		}
		if testAddress.Suburb != addressParts.Suburb {
			hasError = true
			errorExpectedString(t, testAddress.Suburb, addressParts.Suburb)
		}
		if testAddress.PostCode != addressParts.PostCode {
			hasError = true
			errorExpectedInt(t, testAddress.PostCode, addressParts.PostCode)
		}
		if testAddress.State != addressParts.State {
			hasError = true
			errorExpectedString(t, testAddress.State, addressParts.State)
		}
		for _, val := range addressParts.AddressStringParts {
			if val != "" {
				hasError = true
				t.Errorf("Address Leftovers: %s", val)
			}
		}

		if hasError {
			spew.Dump(addressParts)
		}

	}

}
