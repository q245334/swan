package sensitivity

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProfileDataRenderer(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	Convey("While checking value in slice", t, func() {
		sliceOfValues := []string{"a", "b"}
		Convey("I should receive true if it exists", func() {
			So(isValueInSlice("a", sliceOfValues), ShouldEqual, true)
		})
		Convey("I should receive false if it does not exist", func() {
			So(isValueInSlice("c", sliceOfValues), ShouldEqual, false)
		})
	})
	Convey("While creating list of unique values from map", t, func() {
		sliceOfValues := []string{"1", "2"}
		Convey("If key is present in map", func() {
			Convey("and value is not in list, I should receive a list with this value", func() {
				elem := map[string]string{"a": "3"}
				returnedList := createUniqueList("a", elem, sliceOfValues)
				So(len(returnedList), ShouldEqual, 1)
				So(returnedList[0], ShouldEqual, "3")
			})
			Convey("and value is in list, I should receive empty list", func() {
				elem := map[string]string{"a": "1"}
				returnedList := createUniqueList("a", elem, sliceOfValues)
				So(len(returnedList), ShouldEqual, 0)
			})
		})
		Convey("If key is not present in map", func() {
			Convey("I should receive empty list", func() {
				elem := map[string]string{"a": "1"}
				returnedList := createUniqueList("c", elem, sliceOfValues)
				So(len(returnedList), ShouldEqual, 0)
			})
		})

	})

	Convey("While calculating average of string values", t, func() {
		Convey("from empty slice, I should receive error", func() {
			slice := []string{}
			value, err := calculateAverage(slice)
			So(value, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
		Convey("from not empty slice, I should not receive error and I should receive proper value", func() {
			slice := []string{"1", "2", "3"}
			value, err := calculateAverage(slice)
			So(*value, ShouldEqual, 2)
			So(err, ShouldBeNil)
		})
	})
	Convey("While getting load point number", t, func() {
		Convey("from correct PhaseID, I should receive proper value and no error", func() {
			loadPoint := 1
			value, err := getLoadPointNumber(fmt.Sprintf("loadpoint_id_%d", loadPoint))
			So(*value, ShouldEqual, loadPoint)
			So(err, ShouldBeNil)
		})
		Convey("from not correct PhaseID, I should receive error", func() {
			loadPoint := 1
			_, err := getLoadPointNumber(fmt.Sprintf("loadpoint_id", loadPoint))
			So(err, ShouldNotBeNil)
		})
	})
}