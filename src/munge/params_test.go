package munge

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonPath(t *testing.T) {
	Convey("Path functions", t, func() {
		inJson := `{"name": "Niamh"}`
		inYaml := `
name: Niamh
children:
  - Seamus
  - Lorcan
  - Sinead
        `
		Convey("Single top JSON value", func() {
			out, err := ExtractPath(inJson, "$.name")
			So(err, ShouldBeNil)
			So(out, ShouldHaveSameTypeAs, "a string")
			So(out.(string), ShouldEqual, "Niamh")
		})
		Convey("Single top JSON value not present errors", func() {
			_, err := ExtractPath(inJson, "$.noname")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "key error")
		})
		Convey("Single top YAML value", func() {
			out, err := ExtractPath(inYaml, "$.name")
			So(err, ShouldBeNil)
			So(out, ShouldHaveSameTypeAs, "a string")
			So(out.(string), ShouldEqual, "Niamh")
		})
		Convey("YAML list value", func() {
			out, err := ExtractPath(inYaml, "$.children[1]")
			So(err, ShouldBeNil)
			So(out, ShouldHaveSameTypeAs, "a string")
			So(out.(string), ShouldEqual, "Lorcan")
		})
	})
}
