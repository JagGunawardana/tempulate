package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckFiles(t *testing.T) {
	Convey("Check files", t, func() {
		contents := []byte("junk")
		ioutil.WriteFile("one", contents, 0644)
		ioutil.WriteFile("two", contents, 0644)
		ioutil.WriteFile("three", contents, 0644)
		Convey("All present", func() {
			err := checkFiles([]string{"one", "two", "three"})
			So(err, ShouldBeNil)
		})
		Convey("One missing", func() {
			err := checkFiles([]string{"one", "two", "three", "five"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Failed to find files: five")
		})
		Convey("All missing", func() {
			err := checkFiles([]string{"nothing", "is", "there"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Failed to find files: nothing,is,there")
		})
		Reset(func() {
			os.Remove("one")
			os.Remove("two")
			os.Remove("three")
		})
	})
}
