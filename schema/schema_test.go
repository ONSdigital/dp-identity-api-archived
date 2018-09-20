package schema

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIdentity_Validate(t *testing.T) {
	Convey("should not return error if identity is valid", t, func() {
		i := &Identity{
			Name:     "Bucky O'Hare",
			Email:    "captain@TheRighteousIndignation.com",
			Password: "S.P.A.C.E",
		}

		err := i.Validate()
		So(err, ShouldBeNil)
	})

	Convey("should error if identity is nil", t, func() {
		var i *Identity = nil
		err := i.Validate()
		So(err, ShouldResemble, ErrIdentityNil)
	})

	Convey("should error if identity.name is nil", t, func() {
		i := &Identity{}
		err := i.Validate()
		So(err, ShouldResemble, ErrNameValidation)
	})

	Convey("should error if identity.email is nil", t, func() {
		i := &Identity{Name: "Bucky O'Hare"}
		err := i.Validate()
		So(err, ShouldResemble, ErrEmailValidation)
	})

	Convey("should error if identity.password is nil", t, func() {
		i := &Identity{Name: "Bucky O'Hare", Email: "captain@TheRighteousIndignation.com"}
		err := i.Validate()
		So(err, ShouldResemble, ErrPasswordValidation)
	})
}
