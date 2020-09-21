package cmd

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetInputStr(t *testing.T) {
	q1 := "01234567891234567890123456789"
	Convey("sdfdsf12-1", t, func() {
		So(getInputStr(q1), ShouldEqual, "0123456789290123456789")
	})
	// 表示断言成
	q2 := "表示断言成表示断言成dfdf我的sdf否则表示失败表示失败"
	Convey("sdfdsf12-2", t, func() {
		So(getInputStr(q2), ShouldEqual, "表示断言成表示断言成29否则表示失败表示失败")
	})
	q3 := "表示断言成表示断言成我的dfdf我的sdf我的否则表示失败表示失败"
	Convey("sdfdsf12-3", t, func() {
		So(getInputStr(q3), ShouldEqual, "表示断言成表示断言成33否则表示失败表示失败")
	})
	q4 := "表示断言成表示断言"
	Convey("sdfdsf12-4", t, func() {
		So(getInputStr(q4), ShouldEqual, "表示断言成表示断言")
	})
	q5 := "012345678"
	Convey("sdfdsf12-5", t, func() {
		So(getInputStr(q5), ShouldEqual, "012345678")
	})
}

func TestGetSha256(t *testing.T) {
	Convey("debug test", t, func() {
		So(getSha256("012345678"), ShouldEqual, "36f50957f5e0b6ee3ef455674da35a86667f3314209dc1514c510fe95e840831")
	})
}

func TestTransform(t *testing.T) {
	ret, _ := transform("Non-Motor Vehicle")
	log.Printf("%#v", ret)
}
