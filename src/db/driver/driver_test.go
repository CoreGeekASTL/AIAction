// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package driver

import (
	"database/sql/driver"
	"testing"

	convey "github.com/smartystreets/goconvey/convey"
)

var _ driver.Driver = &DoNothingDriver{}
var _ driver.Conn = &DoNothingConnect{}

type DoNothingDriver struct {
}

func (d *DoNothingDriver) Open(name string) (driver.Conn, error) {
	return &DoNothingConnect{}, nil
}

type DoNothingConnect struct {
	PrepareInput string
}

func (d *DoNothingConnect) Prepare(query string) (driver.Stmt, error) {
	d.PrepareInput = query
	return nil, nil
}

func (d *DoNothingConnect) Close() error {
	return nil
}

func (d *DoNothingConnect) Begin() (driver.Tx, error) {
	return nil, nil
}

func TestDriver(t *testing.T) {
	convey.Convey("test driver", t, func() {
		d := Decorator{
			Driver: &DoNothingDriver{},
		}

		conn, err := d.Open("test")
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("test sql invalid", func() {
			_, err = conn.Prepare(`call test"`)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("test table invalid", func() {
			_, err = conn.Prepare(`select * from"t_test"`)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("test select", func() {
			_, err = conn.Prepare(`select t.* from "t_test" t`)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("test insert", func() {
			_, err = conn.Prepare(`insert into "t_test"() values()`)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("test update", func() {
			_, err = conn.Prepare(`update "t_test" set 1=1 where 1=1`)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("test delete", func() {
			_, err = conn.Prepare(`delete from "t_test" where 1=1`)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
