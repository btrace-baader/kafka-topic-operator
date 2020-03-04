package controllers

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestContainsString(t *testing.T) {
	Convey("Test if a slice contains a string", t, func() {
		Convey("Positive tests", func() {
			Convey("match non-empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := "abc"
				b := containsString(a, str)
				So(b, ShouldEqual, true)
			})
			Convey("match empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := ""
				b := containsString(a, str)
				So(b, ShouldEqual, true)
			})
		})
		Convey("Negative tests", func() {
			Convey("match non-empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := "trs"
				b := containsString(a, str)
				So(b, ShouldNotEqual, true)
			})
			Convey("match empty", func() {
				a := []string{"abc", "a", "zyx"}
				str := ""
				b := containsString(a, str)
				So(b, ShouldNotEqual, true)
			})
		})
	})
}

func TestRemoveString(t *testing.T) {
	Convey("Test if a string a removed from slice", t, func() {
		Convey("Positive tests", func() {
			Convey("remove non-empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := "abc"
				b := removeString(a, str)
				So(b, ShouldNotContain, "abc")
			})
			Convey("remove empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := ""
				b := removeString(a, str)
				So(b, ShouldNotContain, "")
			})
		})
		Convey("Negative tests", func() {
			Convey("check non-empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := "abcd"
				b := removeString(a, str)
				So(b, ShouldContain, "abc")
			})
			Convey("check empty", func() {
				a := []string{"abc", "a", "zyx", ""}
				str := "xyz"
				b := removeString(a, str)
				So(b, ShouldContain, "")
			})
		})
	})
}
