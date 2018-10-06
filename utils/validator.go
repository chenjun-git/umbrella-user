package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	libphone "github.com/nyaruka/phonenumbers"
)

func ValidEmail(strEmail string) error {
	if strEmail == "" {
		return errors.New("empty string")
	}
	if govalidator.IsEmail(strEmail) {
		return nil
	} else {
		return errors.New("invalid email")
	}
}

// var v string = "+86 13681454478"
func ValidPhone(str string) error {
	if str == "" {
		return errors.New("empty string")
	}

	if str[0] != '+' {
		return errors.New("phone format error, expect {+CC}{Space}{Number}")
	}
	pieces := strings.Split(str, " ")
	if len(pieces) < 2 {
		return errors.New("phone format error, expect {+CC}{Space}{Number}")
	}

	countryCodeStr := pieces[0]
	countryCode, err := strconv.Atoi(strings.TrimLeft(countryCodeStr, "+"))
	if err != nil {
		return errors.New("invalid country code")
	}

	region := libphone.GetRegionCodeForCountryCode(countryCode)
	if region == libphone.UNKNOWN_REGION {
		return errors.New("unknown country code region")
	}

	phoneNumber, err := libphone.Parse(str, region)
	if err != nil {
		return errors.New("parse phone number failed")
	}

	if !libphone.IsValidNumberForRegion(phoneNumber, region) {
		return errors.New("mismatch phone region")
	}

	return nil
}