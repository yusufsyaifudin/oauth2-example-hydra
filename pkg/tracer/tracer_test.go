package tracer

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNullCloser_Close(t *testing.T) {
	convey.Convey("Null Closer", t, func() {
		convey.Convey("Should always return nil on close", func() {
			closer := nullCloser{}
			convey.So(closer, convey.ShouldNotBeNil)

			err := closer.Close()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestNew(t *testing.T) {
	convey.Convey("New tracer", t, func() {
		convey.Convey("Should return not nil when disabled", func() {
			tracer, closer := New(false, "TRACER", "localhost:5775", 1)
			convey.So(tracer, convey.ShouldNotBeNil)
			convey.So(closer, convey.ShouldNotBeNil)
		})

		convey.Convey("Should using default service name when enabled and empty service name", func() {
			tracer, closer := New(true, "", "localhost:5775", 1)
			convey.So(tracer, convey.ShouldNotBeNil)
			convey.So(closer, convey.ShouldNotBeNil)
		})

		convey.Convey("Should fallback to null closer when host port is not valid", func() {
			tracer, closer := New(true, "TRACER", "//\\\\//se.org", 1)
			convey.So(tracer, convey.ShouldNotBeNil)
			convey.So(closer, convey.ShouldNotBeNil)
			convey.So(closer, convey.ShouldResemble, new(nullCloser))
		})

		convey.Convey("Should return not nil when enabled", func() {
			tracer, closer := New(true, "TRACER", "localhost:5775", 1)
			convey.So(tracer, convey.ShouldNotBeNil)
			convey.So(closer, convey.ShouldNotBeNil)
		})

	})
}
