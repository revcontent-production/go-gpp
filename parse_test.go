package gpp

import (
	"fmt"
	"github.com/revcontent-production/go-gpp/constants"
	"github.com/revcontent-production/go-gpp/sections"
	"github.com/revcontent-production/go-gpp/sections/uspca"
	"github.com/revcontent-production/go-gpp/sections/uspva"
	"testing"

	"github.com/stretchr/testify/assert"
)

type gppTestData struct {
	description   string
	gppString     string
	expected      GppContainer
	expectedError []error
}

func TestParse(t *testing.T) {
	testData := map[string]gppTestData{
		"gpp-tcf": {
			description: "GPP string with EU TCF V2",
			gppString:   "DBABM~CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{2},
				Sections: []Section{GenericSection{sectionID: 2,
					value: "CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA"}},
			},
		},
		"gpp-tcf-valid-quantum": { // header is valid base64 quantum, should gracefully decode correctly
			description: "GPP string with EU TCF V2",
			gppString:   "DBABMA~CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{2},
				Sections: []Section{GenericSection{sectionID: 2,
					value: "CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA"}},
			},
		},
		"gpp-tcf-usp": {
			description: "GPP string with EU TCF v2 and US Privacy",
			gppString:   "DBACNY~CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA~1YNN",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{2, 6},
				Sections: []Section{GenericSection{sectionID: 2,
					value: "CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA"},
					GenericSection{sectionID: 6,
						value: "1YNN"}},
			},
		},
		"gpp-tcfca-usp": {
			description: "GPP string with Canadian TCF and US Privacy",
			gppString:   "DBABjw~CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA~1YNN",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{5, 6},
				Sections: []Section{GenericSection{sectionID: 5,
					value: "CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA"},
					GenericSection{sectionID: 6,
						value: "1YNN"}},
			},
		},
		"gpp-uspca": {
			description: "GPP string with USPCA",
			gppString:   "DBABBgA~xlgWEYCZAA",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{8},
				Sections: []Section{uspca.USPCA{
					CoreSegment: uspca.USPCACoreSegment{
						Version:                     49,
						SaleOptOutNotice:            2,
						SharingOptOutNotice:         1,
						SensitiveDataLimitUseNotice: 1,
						SaleOptOut:                  2,
						SharingOptOut:               0,
						SensitiveDataProcessing: []byte{
							0, 1, 1, 2, 0, 1, 0, 1, 2,
						},
						KnownChildSensitiveDataConsents: []byte{
							0, 0,
						},
						PersonalDataConsents:    0,
						MspaCoveredTransaction:  2,
						MspaOptOutOptionMode:    1,
						MspaServiceProviderMode: 2,
					},
					GPCSegment: sections.CommonUSGPCSegment{
						SubsectionType: 1,
						Gpc:            false,
					},
					SectionID: 8,
					Value:     "xlgWEYCZAA"},
				},
			},
		},
		"gpp-uspva": {
			description: "GPP string with USPVA",
			gppString:   "DBABRgA~bSFgmiU",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{9},
				Sections: []Section{uspva.USPVA{
					CoreSegment: sections.CommonUSCoreSegment{
						Version:                         27,
						SharingNotice:                   1,
						SaleOptOutNotice:                0,
						TargetedAdvertisingOptOutNotice: 2,
						SaleOptOut:                      0,
						TargetedAdvertisingOptOut:       1,
						SensitiveDataProcessing: []byte{
							1, 2, 0, 0, 2, 1, 2, 2,
						},
						KnownChildSensitiveDataConsents: []byte{0},
						MspaCoveredTransaction:          2,
						MspaOptOutOptionMode:            1,
						MspaServiceProviderMode:         1,
					},
					SectionID: constants.SectionUSPVA,
					Value:     "bSFgmiU"},
				},
			},
		},
		"gpp-tcf-tcfca-usp": {
			description: "GPP string with EU TCF v2, Canadian TCF and US Privacy",
			gppString:   "DBACOe~CPrzpwAPrzpwAEXahAENC7CwAP_AAH_AACiQIgAB4C5GQCFDeHpdAJsUAAQDQMhAAKAgAAQFgYABCBoAAIwCAAAwAACCAAoCAAIAIABBAAEAAAAAAAEAQAAAAAEAAEAAAAAAIAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAyAAAAAIAEEAAAAACAAEAAAgAABAAAgAAAAAAAAAAAAAIAAAAAAAAAAAEEQIGAACwAKgAXAAyAByAEAAQgAkABkADQAHIAPIAfAB_AEQARQAmABPACkAF8AMQAZgA0AB-AEIAKMAUoAyIBlAGWAOeAdwB3gEDgIOAhABEQCLAE7AKCAU8AtIBdQDFAGvAOoAvMBkwDLAGfANVAfuBBQCIAAAA~CPrzpwAPrzpwAEXahAENC7CgAf-AAP-AAAiAAHgLkZAIUN4el0AmxQABANAyEAAoCAABAWBgAEIGgAAjAIAADAAAIIACgIAAgAgAEEAAQAAAAAAAQBAAAAAAQAAQAAAAAAgAAAAAAAAAAAAAAgAAAAAAAAAAAAAAADIAAAAAgAQQAAAAAIAAQAACAAAEAACAAAAAAAAAAAAAAgAAAAAAAAAAAQRAgYAALAAqABcADIAHIAQABCACQAGQANAAcgA8gB8AH8ARABFACYAE8AKQAXwAxABmADQAH4AQgAowBSgDIgGUAZYA54B3AHeAQOAg4CEAERAIsATsAoIBTwC0gF1AMUAa8A6gC8wGTAMsAZ8A1UB-4EFAIgA~1-N-",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{2, 5, 6},
				Sections: []Section{GenericSection{sectionID: 2,
					value: "CPrzpwAPrzpwAEXahAENC7CwAP_AAH_AACiQIgAB4C5GQCFDeHpdAJsUAAQDQMhAAKAgAAQFgYABCBoAAIwCAAAwAACCAAoCAAIAIABBAAEAAAAAAAEAQAAAAAEAAEAAAAAAIAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAyAAAAAIAEEAAAAACAAEAAAgAABAAAgAAAAAAAAAAAAAIAAAAAAAAAAAEEQIGAACwAKgAXAAyAByAEAAQgAkABkADQAHIAPIAfAB_AEQARQAmABPACkAF8AMQAZgA0AB-AEIAKMAUoAyIBlAGWAOeAdwB3gEDgIOAhABEQCLAE7AKCAU8AtIBdQDFAGvAOoAvMBkwDLAGfANVAfuBBQCIAAAA"},
					GenericSection{sectionID: 5,
						value: "CPrzpwAPrzpwAEXahAENC7CgAf-AAP-AAAiAAHgLkZAIUN4el0AmxQABANAyEAAoCAABAWBgAEIGgAAjAIAADAAAIIACgIAAgAgAEEAAQAAAAAAAQBAAAAAAQAAQAAAAAAgAAAAAAAAAAAAAAgAAAAAAAAAAAAAAADIAAAAAgAQQAAAAAIAAQAACAAAEAACAAAAAAAAAAAAAAgAAAAAAAAAAAQRAgYAALAAqABcADIAHIAQABCACQAGQANAAcgA8gB8AH8ARABFACYAE8AKQAXwAxABmADQAH4AQgAowBSgDIgGUAZYA54B3AHeAQOAg4CEAERAIsATsAoIBTwC0gF1AMUAa8A6gC8wGTAMsAZ8A1UB-4EFAIgA"},
					GenericSection{sectionID: 6,
						value: "1-N-"}},
			},
		},
		"gpp-empty": {
			description: "GPP string no data",
			gppString:   "DBAA",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{},
				Sections:     []Section{},
			},
		},
		"gpp-tcf-error": {
			description: "GPP string with EU TCF V2",
			gppString:   "DBGBM~CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA",
			expected: GppContainer{
				Version:      1,
				SectionTypes: []constants.SectionID{2},
				Sections: []Section{GenericSection{sectionID: 2,
					value: "CPXxRfAPXxRfAAfKABENB-CgAAAAAAAAAAYgAAAAAAAA"}},
			},
			expectedError: []error{fmt.Errorf("error parsing GPP header, section identifiers: error reading an int offset value in a Range(Fibonacci) entry(1): error reading bit 4 of Integer(Fibonacci): expected 1 bit at bit 32, but the byte array was only 4 bytes long")},
		},
		"gpp-uspca-error": {
			description:   "GPP string with USPCA",
			gppString:     "DBABBgA~xlgWE",
			expectedError: []error{fmt.Errorf("error parsing uspca consent string: unable to set field CoreSegment.SensitiveDataProcessing due to parse error: expected 2 bits to start at bit 32, but the byte array was only 4 bytes long")},
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			result, err := Parse(test.gppString)

			if len(test.expectedError) == 0 {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, result)
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func TestFailFastHeaderValidate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := failFastHeaderValidate("DBABM")
		assert.NoError(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		err := failFastHeaderValidate("")
		assert.EqualError(t, err, "error parsing GPP header, should be at least 4 bytes long")
	})

	t.Run("short", func(t *testing.T) {
		err := failFastHeaderValidate("DB")
		assert.EqualError(t, err, "error parsing GPP header, should be at least 4 bytes long")
	})

	t.Run("invalid-type", func(t *testing.T) {
		err := failFastHeaderValidate("AAAA")
		assert.EqualError(t, err, "error parsing GPP header, header must have type=3")
	})
}

// go test -bench="^BenchmarkParse$" -benchmem .
// BenchmarkParse-8          625084              1912 ns/op            1472 B/op         48 allocs/op (Apple M1 Pro)
func BenchmarkParse(b *testing.B) {
	const gppString = "DBABrGA~DSJgmkoZJSA.YA~BlgWEYCY.QA~BSFgmiU~bSFgmJQ.YA~BWJYJllA~bSFgmSZQ.YA"
	for i := 0; i < b.N; i++ {
		_, err := Parse(gppString)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestParse1(t *testing.T) {
	result, err := Parse("DBABBg~BUoAAAA")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Sections))

	result, err = Parse("DBABTYA~1YNY~000001010101010000101010000000000000000000000000000000010101")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Sections))

	result, err = Parse("DBABM~")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Sections))

	result, err = Parse("DBABBg~BAAAAAA")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Sections))
}

func TestParse2(t *testing.T) {
	//result, err := Parse("DBABRg~BVoAAA")
	result, err := Parse("DBABRgA~bSFgmiU")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Sections))

	result, err = Parse("DBABFg~BVKAAAA")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Sections))
}

func TestParse3(t *testing.T) {
	result, err := Parse("DBACOe~CPzyhEAPzyhEAEXdtAENDXCwAP_AAH_AACiQI9AB4C5EQCFDcHJNAIoUAAQDQIhAACAgAAABgYAACBoAAIwAAAAwAAAAAAoCAAAAIABAAAEAAAAAAAEAAAAAAAEAAEAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAiAAAAAIAEEAAAAACAAEAAAgAABAAAgAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAQQUgSgALAAqABcADIAHAAQAAkABkADQAHIAPAAfQBEAEUAJgATwApABfADMAGgAPwAhABSwDKAMsAc8A7gDvAIHAQcBCACLAFPALqAvMBkwDLAGfANVAfuBBQAA~CPzyhEAPzyhEAEXdtAENDXCgAf-AAP-AAAj0AHgLkRAIUNwck0AihQABANAiEAAICAAAAGBgAAIGgAAjAAAADAAAAAACgIAAAAgAEAAAQAAAAAAAQAAAAAAAQAAQAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACIAAAAAgAQQAAAAAIAAQAACAAAEAACAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAABBBSBKAAsACoAFwAMgAcABAACQAGQANAAcgA8AB9AEQARQAmABPACkAF8AMwAaAA_ACEAFLAMoAywBzwDuAO8AgcBBwEIAIsAU8AuoC8wGTAMsAZ8A1UB-4EFA~1YN-")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(result.Sections))
}
