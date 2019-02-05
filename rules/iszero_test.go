package rules

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsZero(t *testing.T) {
	Convey("test is zero", t, func() {
		for _, val := range []interface{}{
			[]string(nil),
			(*int)(nil),
			0,
			"",
			false,
			uint(0),
			uintptr(0),
			float64(0),
		} {
			So(isZero(reflect.ValueOf(val)), ShouldBeTrue)
		}
		testStr := ""
		testInt := 0
		for _, val := range []interface{}{
			1,
			&testInt,
			&testStr,
			[]string{"test"},
		} {
			So(isZero(reflect.ValueOf(val)), ShouldBeFalse)
		}
	})
}
