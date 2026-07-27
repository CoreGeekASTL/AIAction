// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package driver 代理gauss db driver，适配beego orm
package driver

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"

	gaussdb "gitee.com/opengauss/openGauss-connector-go-pq"
)

var extractTableNameError = errors.New("failed to extract the table name from the SQL statement")

type Decorator struct {
	driver.Driver
}

func (d Decorator) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &ConnDecorator{Conn: conn}, nil
}

type ConnDecorator struct {
	driver.Conn
}

func (cd *ConnDecorator) Prepare(query string) (driver.Stmt, error) {
	tableName, err := extractTableName(query)
	if err == nil {
		query = strings.Replace(query, tableName, strings.Trim(tableName, `"`), -1)
	}
	return cd.Conn.Prepare(query)
}

func extractTableNameForTag(tag, sql string) (string, error) {
	startIdx := len(tag)
	sql = strings.TrimSpace(sql[startIdx:])
	// 获取表名
	tableNameEndIdx := strings.Index(sql, " ")
	if tableNameEndIdx == -1 {
		return sql, extractTableNameError
	}
	return sql[:tableNameEndIdx], nil
}

// 从 SQL 语句中提取表名
func extractTableName(sql string) (string, error) {
	// 清理 SQL 语句中的多余空格，并统一为大写
	sql = strings.ToUpper(strings.TrimSpace(sql))

	if strings.HasPrefix(sql, "SELECT") {
		startIdx := strings.Index(sql, "FROM")
		if startIdx == -1 {
			return sql, extractTableNameError
		}

		return extractTableNameForTag("FROM", strings.TrimSpace(sql[startIdx:]))
	}

	if strings.HasPrefix(sql, "INSERT INTO") {
		return extractTableNameForTag("INSERT INTO", sql)
	}

	if strings.HasPrefix(sql, "UPDATE") {
		return extractTableNameForTag("UPDATE", sql)
	}

	if strings.HasPrefix(sql, "DELETE FROM") {
		return extractTableNameForTag("DELETE FROM", sql)
	}

	// 如果不匹配任何 SQL 类型，返回错误
	return "", extractTableNameError
}

func init() {
	d := &Decorator{
		Driver: &gaussdb.Driver{},
	}
	sql.Register("gaussdb_1", d)
}
