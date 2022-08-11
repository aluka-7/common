package common_test

import (
	"testing"

	"github.com/aluka-7/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPage(t *testing.T) {
	Convey("Test Page", t, func() {
		page := common.DefaultPagination()
		Convey("Test Set Total Record", func() {
			page.SetTotalRecord(10)
			limit, offset := page.Limit()
			So(limit, ShouldEqual, 20)
			So(offset, ShouldEqual, 0)
		})
		Convey("Test Set Page Number", func() {
			page.SetPageNumber(2)
			limit, offset := page.Limit()
			So(limit, ShouldEqual, 20)
			So(offset, ShouldEqual, 20)
		})
	})
	Convey("Test Query", t, func() {
		query := common.Query{}
		Convey("Test Mark page with default", func() {
			page := query.MarkPage()
			limit, offset := page.Limit()
			So(limit, ShouldEqual, 20)
			So(offset, ShouldEqual, 0)
		})
	})
}
