package config

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {

			Convey("There should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The values should be set to the expected defaults", func() {
				So(cfg.BindAddr, ShouldEqual, ":23800")
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthCheckTimeout, ShouldEqual, 2*time.Second)
				So(cfg.MongoConfig.Database, ShouldEqual, "identities")
				So(cfg.MongoConfig.IdentityCollection, ShouldEqual, "identities")
				So(cfg.MongoConfig.TokenCollection, ShouldEqual, "tokens")
				So(cfg.MongoConfig.BindAddr, ShouldEqual, "localhost:27017")
			})
		})
	})
}
