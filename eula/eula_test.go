package eula

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func Test_DisplayEula(t *testing.T) {
	Convey("Check EULA Display", t, func() {
		So(DisplayEULA(), ShouldContainSubstring, "Software License Agreement")
	})
}

func Test_askForConfirmation(t *testing.T) {
	Convey("Check user input decline EULA", t, func() {
		So(AskForConfirmation("", strings.NewReader("n")), ShouldEqual, false)
		So(AskForConfirmation("", strings.NewReader("N")), ShouldEqual, false)
		So(AskForConfirmation("", strings.NewReader("no")), ShouldEqual, false)
		So(AskForConfirmation("", strings.NewReader("No")), ShouldEqual, false)
		So(AskForConfirmation("", strings.NewReader("NO")), ShouldEqual, false)
		So(AskForConfirmation("", strings.NewReader("nO")), ShouldEqual, false)
	})
	Convey("Check user input accept EULA", t, func() {
		So(AskForConfirmation("", strings.NewReader("y")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("Y")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("yes")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("yEs")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("yeS")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("yES")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("Yes")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("YeS")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("YEs")), ShouldEqual, true)
		So(AskForConfirmation("", strings.NewReader("YES")), ShouldEqual, true)
	})
}

func Test_postString(t *testing.T) {
	data := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	Convey("Does string element exist in string slice - invalid element", t, func() {
		So(posString(data, "onee"), ShouldEqual, -1)
		So(posString(data, "on"), ShouldEqual, -1)
		So(posString(data, "wdwdw"), ShouldEqual, -1)
		So(posString(data, ""), ShouldEqual, -1)
	})
	Convey("Does string element exist in string slice - valid element", t, func() {
		So(posString(data, "one"), ShouldEqual, 0)
	})
}

func Test_containsString(t *testing.T) {
	data := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	Convey("Does string element exist in string slice - invalid element", t, func() {
		So(containsString(data, "onee"), ShouldEqual, false)
		So(containsString(data, "on"), ShouldEqual, false)
		So(containsString(data, "wdwdw"), ShouldEqual, false)
		So(containsString(data, ""), ShouldEqual, false)
	})
	Convey("Does string element exist in string slice - valid element", t, func() {
		So(containsString(data, "one"), ShouldEqual, true)
	})
}
