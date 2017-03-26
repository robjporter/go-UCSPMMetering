package functions

import (
	"github.com/robjporter/go-functions/as"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_CurrentMonthName(t *testing.T) {
	now := time.Now()
	_, currentMonth, _ := now.Date()
	nextMonth := currentMonth + 1
	previousMonth := currentMonth - 1
	Convey("Should equal the month today", t, func() {
		So(CurrentMonthName(), ShouldEqual, currentMonth.String())
	})
	Convey("Should not equal next month", t, func() {
		So(CurrentMonthName(), ShouldNotEqual, nextMonth.String())
	})
	Convey("Should not equal previous month", t, func() {
		So(CurrentMonthName(), ShouldNotEqual, previousMonth.String())
	})
}

func Test_CurrentYear(t *testing.T) {
	now := time.Now()
	currentYear, _, _ := now.Date()
	nextYear := currentYear + 1
	previoiusYear := currentYear - 1
	Convey("Should equal the year today", t, func() {
		So(CurrentYear(), ShouldEqual, as.ToString(currentYear))
	})
	Convey("Should not equal next year", t, func() {
		So(CurrentYear(), ShouldNotEqual, as.ToString(nextYear))
	})
	Convey("Should not equal previous year", t, func() {
		So(CurrentYear(), ShouldNotEqual, as.ToString(previoiusYear))
	})
}

func Test_IsYear(t *testing.T) {
	now := time.Now()
	currentYear, _, _ := now.Date()
	Convey("Is Valid year with string", t, func() {
		So(IsYear("test"), ShouldEqual, as.ToString(currentYear))
	})
	Convey("Is Valid year with invalid year number", t, func() {
		So(IsYear("4444"), ShouldEqual, as.ToString(currentYear))
	})
	Convey("Is Valid year with valid year number", t, func() {
		So(IsYear("2016"), ShouldEqual, "2016")
	})
}

func Test_isValidYear(t *testing.T) {
	Convey("Is Valid year with string", t, func() {
		So(isValidYear("test"), ShouldEqual, false)
	})
	Convey("Is Valid year with invalid date 1999", t, func() {
		So(isValidYear("1999"), ShouldEqual, false)
	})
	Convey("Is Valid year with invalid date 2031", t, func() {
		So(isValidYear("2031"), ShouldEqual, false)
	})
	Convey("Is Valid year with valid date 2000", t, func() {
		So(isValidYear("2000"), ShouldEqual, true)
	})
	Convey("Is Valid year with valid date 2030", t, func() {
		So(isValidYear("2030"), ShouldEqual, true)
	})
}

func Test_isNumber(t *testing.T) {
	Convey("Is valid number with string", t, func() {
		So(isNumber("test"), ShouldEqual, false)
	})
	Convey("Is valid number with number and string", t, func() {
		So(isNumber("4test"), ShouldEqual, false)
	})
	Convey("Is valid number with positive string", t, func() {
		So(isNumber("4"), ShouldEqual, true)
		So(isNumber("44"), ShouldEqual, true)
		So(isNumber("4444"), ShouldEqual, true)
		So(isNumber("4444444444444"), ShouldEqual, true)
	})
	Convey("Is valid number with negative string", t, func() {
		So(isNumber("-4"), ShouldEqual, true)
		So(isNumber("-44"), ShouldEqual, true)
		So(isNumber("-4444"), ShouldEqual, true)
		So(isNumber("-4444444444444"), ShouldEqual, true)
	})
}

func Test_IsMonth(t *testing.T) {
	Convey("Is valid month with invalid string", t, func() {
		So(IsMonth("test"), ShouldEqual, "")
	})
	Convey("Is valid month with valid string", t, func() {
		So(IsMonth("jan"), ShouldEqual, "January")
		So(IsMonth("feb"), ShouldEqual, "February")
		So(IsMonth("mar"), ShouldEqual, "March")
		So(IsMonth("apr"), ShouldEqual, "April")
		So(IsMonth("january"), ShouldEqual, "January")
		So(IsMonth("february"), ShouldEqual, "February")
		So(IsMonth("march"), ShouldEqual, "March")
		So(IsMonth("april"), ShouldEqual, "April")
		So(IsMonth("may"), ShouldEqual, "May")
		So(IsMonth("june"), ShouldEqual, "June")
		So(IsMonth("july"), ShouldEqual, "July")
		So(IsMonth("august"), ShouldEqual, "August")
		So(IsMonth("september"), ShouldEqual, "September")
		So(IsMonth("october"), ShouldEqual, "October")
		So(IsMonth("november"), ShouldEqual, "November")
		So(IsMonth("december"), ShouldEqual, "December")
	})
}

func Test_monthContains(t *testing.T) {
	Convey("Month name contains with invalid string", t, func() {
		So(monthContains("test", "jan", "january"), ShouldEqual, false)
		So(monthContains("ja", "jan", "january"), ShouldEqual, false)
	})
	Convey("Month name contains with start string", t, func() {
		So(monthContains("jan", "jan", "january"), ShouldEqual, true)
		So(monthContains("may", "may", "may"), ShouldEqual, true)
	})
	Convey("Month name contains with end string", t, func() {
		So(monthContains("january", "jan", "january"), ShouldEqual, true)
		So(monthContains("may", "may", "may"), ShouldEqual, true)
	})
	Convey("Month name contains with randon string length", t, func() {
		So(monthContains("janua", "jan", "january"), ShouldEqual, true)
		So(monthContains("januay", "jan", "january"), ShouldEqual, false)
		So(monthContains("ma", "may", "may"), ShouldEqual, false)
	})
}

func Test_GetTimestampStartOfMonth(t *testing.T) {
	Convey("Get start of month timestamp with invalid month name and valid year", t, func() {
		So(GetTimestampStartOfMonth("test", 2016), ShouldEqual, 0)
	})
	Convey("Get start of month timestamp with invalid month name and invalid year", t, func() {
		So(GetTimestampStartOfMonth("test", 4444), ShouldEqual, 0)
	})
	Convey("Get start of month timestamp with valid month name and invalid year", t, func() {
		So(GetTimestampStartOfMonth("january", 4444), ShouldEqual, 0)
	})
	Convey("Get start of month timestamp with valid month name and valid year", t, func() {
		So(GetTimestampStartOfMonth("january", 2017), ShouldEqual, 1483228800)
		So(GetTimestampStartOfMonth("february", 2017), ShouldEqual, 1485907200)
		So(GetTimestampStartOfMonth("march", 2017), ShouldEqual, 1488326400)
		So(GetTimestampStartOfMonth("april", 2017), ShouldEqual, 1491001200)
		So(GetTimestampStartOfMonth("may", 2017), ShouldEqual, 1493593200)
		So(GetTimestampStartOfMonth("june", 2017), ShouldEqual, 1496271600)
		So(GetTimestampStartOfMonth("july", 2017), ShouldEqual, 1498863600)
		So(GetTimestampStartOfMonth("august", 2017), ShouldEqual, 1501542000)
		So(GetTimestampStartOfMonth("september", 2017), ShouldEqual, 1504220400)
		So(GetTimestampStartOfMonth("october", 2017), ShouldEqual, 1506812400)
		So(GetTimestampStartOfMonth("november", 2017), ShouldEqual, 1509494400)
		So(GetTimestampStartOfMonth("december", 2017), ShouldEqual, 1512086400)
	})
}

func Test_GetTimestampEndOfMonth(t *testing.T) {
	Convey("Get end of month timestamp with invalid month name and valid year", t, func() {
		So(GetTimestampEndOfMonth("test", 2016), ShouldEqual, 0)
	})
	Convey("Get end of month timestamp with invalid month name and invalid year", t, func() {
		So(GetTimestampEndOfMonth("test", 4444), ShouldEqual, 0)
	})
	Convey("Get end of month timestamp with valid month name and invalid year", t, func() {
		So(GetTimestampEndOfMonth("january", 4444), ShouldEqual, 0)
	})
	Convey("Get end of month timestamp with valid month name and valid year", t, func() {
		So(GetTimestampEndOfMonth("january", 2017), ShouldEqual, 1485907199)
		So(GetTimestampEndOfMonth("february", 2017), ShouldEqual, 1488326399)
		So(GetTimestampEndOfMonth("march", 2017), ShouldEqual, 1491001199)
		So(GetTimestampEndOfMonth("april", 2017), ShouldEqual, 1493593199)
		So(GetTimestampEndOfMonth("may", 2017), ShouldEqual, 1496271599)
		So(GetTimestampEndOfMonth("june", 2017), ShouldEqual, 1498863599)
		So(GetTimestampEndOfMonth("july", 2017), ShouldEqual, 1501541999)
		So(GetTimestampEndOfMonth("august", 2017), ShouldEqual, 1504220399)
		So(GetTimestampEndOfMonth("september", 2017), ShouldEqual, 1506812399)
		So(GetTimestampEndOfMonth("october", 2017), ShouldEqual, 1509494399)
		So(GetTimestampEndOfMonth("november", 2017), ShouldEqual, 1512086399)
		So(GetTimestampEndOfMonth("december", 2017), ShouldEqual, 1514764799)
	})
}

func Test_GetStartOfMonth(t *testing.T) {
	Convey("Get end of month day with invalid month name and valid year", t, func() {
		So(GetStartOfMonth("test", 2016), ShouldEqual, "")
	})
	Convey("Get end of month day with invalid month name and invalid year", t, func() {
		So(GetStartOfMonth("test", 4444), ShouldEqual, "")
	})
	Convey("Get end of month day with valid month name and invalid year", t, func() {
		So(GetStartOfMonth("january", 4444), ShouldEqual, "")
	})
	Convey("Get end of month day with valid month name and valid year", t, func() {
		So(GetStartOfMonth("january", 2017), ShouldStartWith, "2017-01-01 00:00:00")
		So(GetStartOfMonth("february", 2017), ShouldStartWith, "2017-02-01 00:00:00")
		So(GetStartOfMonth("march", 2017), ShouldStartWith, "2017-03-01 00:00:00")
		So(GetStartOfMonth("april", 2017), ShouldStartWith, "2017-04-01 00:00:00")
		So(GetStartOfMonth("may", 2017), ShouldStartWith, "2017-05-01 00:00:00")
		So(GetStartOfMonth("june", 2017), ShouldStartWith, "2017-06-01 00:00:00")
		So(GetStartOfMonth("july", 2017), ShouldStartWith, "2017-07-01 00:00:00")
		So(GetStartOfMonth("august", 2017), ShouldStartWith, "2017-08-01 00:00:00")
		So(GetStartOfMonth("september", 2017), ShouldStartWith, "2017-09-01 00:00:00")
		So(GetStartOfMonth("october", 2017), ShouldStartWith, "2017-10-01 00:00:00")
		So(GetStartOfMonth("november", 2017), ShouldStartWith, "2017-11-01 00:00:00")
		So(GetStartOfMonth("december", 2017), ShouldStartWith, "2017-12-01 00:00:00")
	})
}

func Test_GetEndOfMonth(t *testing.T) {
	Convey("Get end of month day with invalid month name and valid year", t, func() {
		So(GetEndOfMonth("test", 2016), ShouldEqual, "")
	})
	Convey("Get end of month day with invalid month name and invalid year", t, func() {
		So(GetEndOfMonth("test", 4444), ShouldEqual, "")
	})
	Convey("Get end of month day with valid month name and invalid year", t, func() {
		So(GetEndOfMonth("january", 4444), ShouldEqual, "")
	})
	Convey("Get end of month day with valid month name and valid year", t, func() {
		So(GetEndOfMonth("january", 2017), ShouldStartWith, "2017-01-31 23:59:59.999999999")
		So(GetEndOfMonth("february", 2017), ShouldStartWith, "2017-02-28 23:59:59.999999999")
		So(GetEndOfMonth("march", 2017), ShouldStartWith, "2017-03-31 23:59:59.999999999")
		So(GetEndOfMonth("april", 2017), ShouldStartWith, "2017-04-30 23:59:59.999999999")
		So(GetEndOfMonth("may", 2017), ShouldStartWith, "2017-05-31 23:59:59.999999999")
		So(GetEndOfMonth("june", 2017), ShouldStartWith, "2017-06-30 23:59:59.999999999")
		So(GetEndOfMonth("july", 2017), ShouldStartWith, "2017-07-31 23:59:59.999999999")
		So(GetEndOfMonth("august", 2017), ShouldStartWith, "2017-08-31 23:59:59.999999999")
		So(GetEndOfMonth("september", 2017), ShouldStartWith, "2017-09-30 23:59:59.999999999")
		So(GetEndOfMonth("october", 2017), ShouldStartWith, "2017-10-31 23:59:59.999999999")
		So(GetEndOfMonth("november", 2017), ShouldStartWith, "2017-11-30 23:59:59.999999999")
		So(GetEndOfMonth("december", 2017), ShouldStartWith, "2017-12-31 23:59:59.999999999")
	})
}

func Test_getMonthPos(t *testing.T) {
	Convey("Get month position with invalid string", t, func() {
		So(getMonthPos("test"), ShouldEqual, 0)
	})
	Convey("Get month position with spelling mistake", t, func() {
		So(getMonthPos("janruary"), ShouldEqual, 0)
	})
	Convey("Get month position with valid month names", t, func() {
		So(getMonthPos("january"), ShouldEqual, 1)
		So(getMonthPos("february"), ShouldEqual, 2)
		So(getMonthPos("march"), ShouldEqual, 3)
		So(getMonthPos("april"), ShouldEqual, 4)
		So(getMonthPos("may"), ShouldEqual, 5)
		So(getMonthPos("june"), ShouldEqual, 6)
		So(getMonthPos("july"), ShouldEqual, 7)
		So(getMonthPos("august"), ShouldEqual, 8)
		So(getMonthPos("september"), ShouldEqual, 9)
		So(getMonthPos("october"), ShouldEqual, 10)
		So(getMonthPos("november"), ShouldEqual, 11)
		So(getMonthPos("december"), ShouldEqual, 12)
	})
}
