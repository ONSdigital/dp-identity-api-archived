package token

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type scenario = struct {
	desc    string
	inputH  int
	inputM  int
	inputS  int
	expectH int
	expectM int
	expectS int
}

func TestNewExpiryHelper(t *testing.T) {
	Convey("should set provide hour, min and sec values for valid input", t, func() {
		helper := NewExpiryHelper(23, 59, 59)
		So(helper.expiryHour, ShouldEqual, 23)
		So(helper.expiryMinute, ShouldEqual, 59)
		So(helper.expirySecond, ShouldEqual, 59)
	})
}

func TestNewExpiryHelperInvalidInput(t *testing.T) {
	scenarios := []scenario{
		{desc: "hour < 0", inputH: -1, inputM: 0, inputS: 0, expectH: 0, expectM: 0, expectS: 0},
		{desc: "min < 0", inputH: 0, inputM: -1, inputS: 0, expectH: 0, expectM: 0, expectS: 0},
		{desc: "sec < 0", inputH: 0, inputM: -1, inputS: 0, expectH: 0, expectM: 0, expectS: 0},
		{desc: "hour > 23", inputH: 100, inputM: 0, inputS: 0, expectH: 0, expectM: 0, expectS: 0},
		{desc: "min > 59", inputH: 0, inputM: 60, inputS: 0, expectH: 0, expectM: 0, expectS: 0},
		{desc: "sec > 59", inputH: 0, inputM: 60, inputS: 0, expectH: 0, expectM: 0, expectS: 0},
	}

	Convey("should set default value if input is invalid", t, func() {

		for i, s := range scenarios {
			helper := NewExpiryHelper(s.inputH, s.inputM, s.expectS)

			So(helper.expiryHour, ShouldEqual, s.expectH)
			So(helper.expiryMinute, ShouldEqual, s.expectM)
			So(helper.expirySecond, ShouldEqual, s.expectS)

			t.Logf("scenario: %d, description: %s successful", i, s.desc)
		}
	})
}
