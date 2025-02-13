package wire

import (
	"errors"
	"strings"
	"testing"

	"github.com/moov-io/base"
	"github.com/stretchr/testify/require"
)

// mockOriginator creates a Originator
func mockOriginator() *Originator {
	o := NewOriginator()
	o.Personal.IdentificationCode = PassportNumber
	o.Personal.Identifier = "1234"
	o.Personal.Name = "Name"
	o.Personal.Address.AddressLineOne = "Address One"
	o.Personal.Address.AddressLineTwo = "Address Two"
	o.Personal.Address.AddressLineThree = "Address Three"
	return o
}

// TestMockOriginator validates mockOriginator
func TestMockOriginator(t *testing.T) {
	o := mockOriginator()

	require.NoError(t, o.Validate(), "mockOriginator does not validate and will break other tests")
}

// TestOriginatorIdentificationCodeValid validates Originator IdentificationCode
func TestOriginatorIdentificationCodeValid(t *testing.T) {
	o := mockOriginator()
	o.Personal.IdentificationCode = "Baseball Card ID"

	err := o.Validate()

	if !base.Match(err, ErrIdentificationCode) {
		t.Errorf("%T: %s", err, err)
	}
}

// TestOriginatorIdentifierAlphaNumeric validates Originator Identifier is alphanumeric
func TestOriginatorIdentifierAlphaNumeric(t *testing.T) {
	o := mockOriginator()
	o.Personal.Identifier = "®"

	err := o.Validate()

	require.EqualError(t, err, fieldError("Identifier", ErrNonAlphanumeric, o.Personal.Identifier).Error())
}

// TestOriginatorNameAlphaNumeric validates Originator Name is alphanumeric
func TestOriginatorNameAlphaNumeric(t *testing.T) {
	o := mockOriginator()
	o.Personal.Name = "®"

	err := o.Validate()

	require.EqualError(t, err, fieldError("Name", ErrNonAlphanumeric, o.Personal.Name).Error())
}

// TestOriginatorAddressLineOneAlphaNumeric validates Originator AddressLineOne is alphanumeric
func TestOriginatorAddressLineOneAlphaNumeric(t *testing.T) {
	o := mockOriginator()
	o.Personal.Address.AddressLineOne = "®"

	err := o.Validate()

	require.EqualError(t, err, fieldError("AddressLineOne", ErrNonAlphanumeric, o.Personal.Address.AddressLineOne).Error())
}

// TestOriginatorAddressLineTwoAlphaNumeric validates Originator AddressLineTwo is alphanumeric
func TestOriginatorAddressLineTwoAlphaNumeric(t *testing.T) {
	o := mockOriginator()
	o.Personal.Address.AddressLineTwo = "®"

	err := o.Validate()

	require.EqualError(t, err, fieldError("AddressLineTwo", ErrNonAlphanumeric, o.Personal.Address.AddressLineTwo).Error())
}

// TestOriginatorAddressLineThreeAlphaNumeric validates Originator AddressLineThree is alphanumeric
func TestOriginatorAddressLineThreeAlphaNumeric(t *testing.T) {
	o := mockOriginator()
	o.Personal.Address.AddressLineThree = "®"

	err := o.Validate()

	require.EqualError(t, err, fieldError("AddressLineThree", ErrNonAlphanumeric, o.Personal.Address.AddressLineThree).Error())
}

// TestOriginatorIdentificationCodeWithNoIdentifier validates Originator Identifier is required
// when IdentificationCode is present
func TestOriginatorIdentificationCodeWithNoIdentifier(t *testing.T) {
	o := mockOriginator()
	o.Personal.IdentificationCode = "D"
	o.Personal.Identifier = ""

	err := o.Validate()

	require.EqualError(t, err, fieldError("Identifier", ErrFieldRequired).Error())
}

// TestOriginatorIdentifierWithNoIdentificationCode validates Originator IdentificationCode
// is required when Identifier is present
func TestOriginatorIdentifierWithNoIdentificationCode(t *testing.T) {
	o := mockOriginator()
	o.Personal.IdentificationCode = ""
	o.Personal.Identifier = "1234567890ABC"

	err := o.Validate()

	require.EqualError(t, err, fieldError("IdentificationCode", ErrFieldRequired).Error())
}

// TestParseOriginatorWrongLength parses a wrong Originator record length
func TestParseOriginatorWrongLength(t *testing.T) {
	var line = "{5000}11234                              Name                               Address One                        Address Two                        Address Three                    "
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseOriginator()

	require.EqualError(t, err, r.parseError(fieldError("Identifier", ErrRequireDelimiter)).Error())
}

// TestParseOriginatorReaderParseError parses a wrong Originator reader parse error
func TestParseOriginatorReaderParseError(t *testing.T) {
	var line = "{5000}11234                              *®ame                               *Address One                        *Address Two                        *Address Three                     *"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseOriginator()

	require.EqualError(t, err, r.parseError(fieldError("Name", ErrNonAlphanumeric, "®ame")).Error())

	_, err = r.Read()

	require.EqualError(t, err, r.parseError(fieldError("Name", ErrNonAlphanumeric, "®ame")).Error())
}

// TestOriginatorTagError validates a Originator tag
func TestOriginatorTagError(t *testing.T) {
	o := mockOriginator()
	o.tag = "{9999}"

	require.EqualError(t, o.Validate(), fieldError("tag", ErrValidTagForType, o.tag).Error())
}

// TestStringOriginatorVariableLength parses using variable length
func TestStringOriginatorVariableLength(t *testing.T) {
	var line = "{5000}B1*"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseOriginator()
	require.NoError(t, err)

	line = "{5000}B1                                                                                                                                                                             NNN"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseOriginator()
	require.ErrorContains(t, err, ErrRequireDelimiter.Error())

	line = "{5000}B1*******"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseOriginator()
	require.ErrorContains(t, err, r.parseError(NewTagMaxLengthErr(errors.New(""))).Error())

	line = "{5000}B1*"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseOriginator()
	require.NoError(t, err)
}

// TestStringOriginatorOptions validates Format() formatted according to the FormatOptions
func TestStringOriginatorOptions(t *testing.T) {
	var line = "{5000}B1*"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseOriginator()
	require.NoError(t, err)

	record := r.currentFEDWireMessage.Originator
	require.Equal(t, "{5000}B1                                 *                                   *                                   *                                   *                                   *", record.String())
	require.Equal(t, "{5000}B1*", record.Format(FormatOptions{VariableLengthFields: true}))
	require.Equal(t, record.String(), record.Format(FormatOptions{VariableLengthFields: false}))
}
