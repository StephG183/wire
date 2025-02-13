package wire

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockExchangeRate creates a ExchangeRate
func mockExchangeRate() *ExchangeRate {
	eRate := NewExchangeRate()
	eRate.ExchangeRate = "1,2345"
	return eRate
}

// TestMockExchangeRate validates mockExchangeRate
func TestMockExchangeRate(t *testing.T) {
	eRate := mockExchangeRate()

	require.NoError(t, eRate.Validate(), "mockExchangeRate does not validate and will break other tests")
}

// TestExchangeRate validates ExchangeRate
func TestExchangeRateNumeric(t *testing.T) {
	eRate := mockExchangeRate()
	eRate.ExchangeRate = "1,--0.00"

	err := eRate.Validate()

	require.EqualError(t, err, fieldError("ExchangeRate", ErrNonAmount, eRate.ExchangeRate).Error())
}

// TestParseExchangeRateWrongLength parses a wrong ExchangeRate record length
func TestParseExchangeRateWrongLength(t *testing.T) {
	var line = "{3720}1,2345"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseExchangeRate()

	require.EqualError(t, err, r.parseError(fieldError("ExchangeRate", ErrRequireDelimiter)).Error())

	_, err = r.Read()

	require.EqualError(t, err, r.parseError(fieldError("ExchangeRate", ErrRequireDelimiter)).Error())
}

// TestParseExchangeRateReaderParseError parses a wrong ExchangeRate reader parse error
func TestParseExchangeRateReaderParseError(t *testing.T) {
	var line = "{3720}1,2345Z     *"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseExchangeRate()

	require.EqualError(t, err, r.parseError(fieldError("ExchangeRate", ErrNonAmount, "1,2345Z")).Error())

	_, err = r.Read()

	require.EqualError(t, err, r.parseError(fieldError("ExchangeRate", ErrNonAmount, "1,2345Z")).Error())
}

// TestExchangeRateTagError validates a ExchangeRate tag
func TestExchangeRateTagError(t *testing.T) {
	eRate := mockCurrencyInstructedAmount()
	eRate.tag = "{9999}"

	err := eRate.Validate()

	require.EqualError(t, err, fieldError("tag", ErrValidTagForType, eRate.tag).Error())
}

// TestStringErrorExchangeRateVariableLength parses using variable length
func TestStringErrorExchangeRateVariableLength(t *testing.T) {
	var line = "{3720}"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseExchangeRate()
	require.NoError(t, err)

	line = "{3720}123         NNN"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseExchangeRate()
	require.ErrorContains(t, err, ErrRequireDelimiter.Error())

	line = "{3720}123***"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseExchangeRate()
	require.ErrorContains(t, err, r.parseError(NewTagMaxLengthErr(errors.New(""))).Error())

	line = "{3720}123*"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseExchangeRate()
	require.NoError(t, err)
}

// TestStringExchangeRateOptions validates Format() formatted according to the FormatOptions
func TestStringExchangeRateOptions(t *testing.T) {
	var line = "{3720}123*"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseExchangeRate()
	require.NoError(t, err)

	record := r.currentFEDWireMessage.ExchangeRate
	require.Equal(t, "{3720}123         *", record.String())
	require.Equal(t, "{3720}123*", record.Format(FormatOptions{VariableLengthFields: true}))
	require.Equal(t, record.String(), record.Format(FormatOptions{VariableLengthFields: false}))
}
