package tokentest

import (
	"github.com/ONSdigital/dp-identity-api/token"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type scenario = struct {
	desc    string
	inputH  int64
	inputM  int64
	inputS  int64
	expectH int64
	expectM int64
	expectS int64
}

func TestNewExpiryHelper(t *testing.T) {
	Convey("should set provide hour, min and sec values for valid input", t, func() {
		helper := token.NewExpiryHelper(23, 59, 59)
		So(helper.GetExpiryHour(), ShouldEqual, 23)
		So(helper.GetExpiryMin(), ShouldEqual, 59)
		So(helper.GetExpirySec(), ShouldEqual, 59)
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
			helper := token.NewExpiryHelper(s.inputH, s.inputM, s.expectS)

			So(helper.GetExpiryHour(), ShouldEqual, s.expectH)
			So(helper.GetExpiryMin(), ShouldEqual, s.expectM)
			So(helper.GetExpirySec(), ShouldEqual, s.expectS)

			t.Logf("scenario: %d, description: %s successful", i, s.desc)
		}
	})
}
