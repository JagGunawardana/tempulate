package munge

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMunge(t *testing.T) {
	noParams := []string{}
	Convey("Munge", t, func() {
		Convey("Null template", func() {
			out, err := MungeFile("Hello world", noParams)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello world")
		})
		Convey("Env var set", func() {
			os.Setenv("TEMPJUNK", "314159265")
			out, err := MungeFile("Hello {{ envdef \"TEMPJUNK\" \"\"}}", noParams)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello 314159265")
		})
		Convey("Env var unset use default", func() {
			out, err := MungeFile("Hello {{ envdef \"TEMPJUNK\" \"my default\"}}", noParams)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello my default")
		})
		Reset(func() {
			os.Unsetenv("TEMPJUNK")
		})
	})
	Convey("Params", t, func() {
		const (
			paramFile1 = "param1.json"
			paramFile2 = "param2.yaml"
		)
		createParam := func(name string, contents string) {
			err := ioutil.WriteFile(name, []byte(contents), 0664)
			So(err, ShouldBeNil)
		}
		Convey("Single value JSON", func() {
			createParam(paramFile1, `{"thing": "world"}`)
			out, err := MungeFile(`Hello {{ value "$.thing" }}`, []string{paramFile1})
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello world")
		})
		Convey("Single value YAML", func() {
			createParam(paramFile1, `thing: world`)
			out, err := MungeFile(`Hello {{ value "$.thing" }}`, []string{paramFile1})
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello world")
		})
		Convey("Single value YAML from second file", func() {
			createParam(paramFile1, `{"notathing": "not world"}`)
			createParam(paramFile2, `thing: world`)
			out, err := MungeFile(`Hello {{ value "$.thing" }}`, []string{paramFile1, paramFile2})
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello world")
		})
		Convey("Single value YAML from first file", func() {
			createParam(paramFile2, `{"notathing": "not world"}`)
			createParam(paramFile1, `thing: world`)
			out, err := MungeFile(`Hello {{ value "$.thing" }}`, []string{paramFile1, paramFile2})
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello world")
		})
		Convey("Single value JSON in both, picks first file in list", func() {
			createParam(paramFile1, `{"mything": "not world"}`)
			createParam(paramFile2, `{"mything": "not of this world"}`)
			out, err := MungeFile(`Hello {{ value "$.mything" }}`, []string{paramFile1, paramFile2})
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello not world")
		})
		Convey("Single value YAML in both, picks first file in list", func() {
			createParam(paramFile1, `mything: not world`)
			createParam(paramFile2, `mything: not of this world`)
			out, err := MungeFile(`Hello {{ value "$.mything" }}`, []string{paramFile1, paramFile2})
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "Hello not world")
		})
		Reset(func() {
			os.Remove(paramFile1)
			os.Remove(paramFile2)
		})
	})
}
