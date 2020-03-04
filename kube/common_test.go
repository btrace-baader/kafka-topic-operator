package kube

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMergeMaps(t *testing.T) {
	Convey("Merging two maps.", t, func() {
		Convey("Positive tests.", func() {
			m1 := map[string]string{
				"key1": "value1",
			}
			m2 := map[string]string{
				"key2": "value2",
			}
			m3, e := mergeMaps(m1, m2)
			So(e, ShouldEqual, nil)
			So(m3, ShouldContainKey, "key1")
			So(m3, ShouldContainKey, "key2")
			So(m3["key1"], ShouldEqual, "value1")
			So(m3["key2"], ShouldEqual, "value2")

		})
		Convey("Negative tests.", func() {
			m1 := map[string]string{
				"key1": "value1",
			}
			m2 := map[string]string{
				"key1": "value1",
				"key2": "value2",
			}
			_, e := mergeMaps(m1, m2)
			So(e, ShouldNotEqual, nil)
		})
	})
}

func TestRemoveEmpty(t *testing.T) {
	Convey("Removing empty entries from map.", t, func() {
		Convey("Negative tests.", func() {
			Convey("Non empty key value", func() {
				m1 := map[string]string{
					"key1": "value1",
					"key2": "value2",
				}
				m2 := removeEmpty(m1)
				So(m2["key1"], ShouldEqual, "value1")
				So(m2["key2"], ShouldEqual, "value2")
			})
			Convey("Empty key, non-empty value.", func() {
				m1 := map[string]string{
					"key1": "value1",
					"":     "value2",
				}
				m2 := removeEmpty(m1)
				So(m2, ShouldContainKey, "")
				So(m2, ShouldContainKey, "key1")
				So(m2["key1"], ShouldEqual, "value1")
				So(m2[""], ShouldEqual, "value2")
			})
		})
		Convey("Positive tests. empty, key value.", func() {
			Convey("empty key value", func() {
				m1 := map[string]string{
					"key1": "value1",
					"":     "",
				}
				m2 := removeEmpty(m1)
				So(m2, ShouldNotContainKey, "")
				So(m2, ShouldContainKey, "key1")
				So(m2["key1"], ShouldEqual, "value1")
			})
			Convey("non-empty key, empty value", func() {
				m1 := map[string]string{
					"key1": "value1",
					"key2": "",
				}
				m2 := removeEmpty(m1)
				So(m2, ShouldContainKey, "key1")
				So(m2, ShouldNotContainKey, "key2")
				So(m2["key1"], ShouldEqual, "value1")
			})

		})

	})
}
