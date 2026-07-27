// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package flagutil 递归解析结构体的 flag tag，并注册为命令行参数。
package flagutil

import (
	"flag"
	"fmt"
	"reflect"
)

// Parse 递归解析结构体的 flag tag，并注册为命令行参数
func Parse(obj interface{}) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	// 创建一个与 obj 类型相同的空对象
	newObj := reflect.New(v.Type()).Elem()

	// 递归解析结构体字段，并注册 flag
	parseStruct(v, newObj, "")
	flag.Parse()
	// 将解析后的 newObj 的值赋值给原始 obj
	v.Set(newObj)
	return newObj
}

// parseStruct 递归解析结构体字段，并注册 flag
func parseStruct(defaultObj, newObj reflect.Value, prefix string) {
	for i := 0; i < newObj.NumField(); i++ {
		field := newObj.Type().Field(i)
		fieldValue := newObj.Field(i)
		defaultFieldValue := defaultObj.Field(i)

		// 获取 flag tag
		flagName := field.Tag.Get("flag")
		fullFlagName := flagName

		if flagName == "" {
			fullFlagName = prefix
		} else if prefix == "" {
			fullFlagName = flagName
		} else {
			fullFlagName = prefix + "-" + flagName
		}
		// 结构体递归解析
		if fieldValue.Kind() == reflect.Struct {
			parseStruct(defaultFieldValue, fieldValue, fullFlagName)
			continue
		}
		// 非结构体无flag不注册
		if flagName == "" {
			continue
		}

		// 获取字段的描述
		desc := field.Tag.Get("desc")

		// 检查字段是否可寻址
		if !fieldValue.CanAddr() {
			continue
		}

		// 创建一个可寻址的指针
		fieldPtr := reflect.NewAt(field.Type, fieldValue.Addr().UnsafePointer()).Interface()

		// 根据类型注册 flag
		switch field.Type.Kind() {
		case reflect.String:
			flag.StringVar(fieldPtr.(*string), fullFlagName, defaultFieldValue.String(), desc)
		case reflect.Int, reflect.Int64:
			flag.IntVar(fieldPtr.(*int), fullFlagName, int(defaultFieldValue.Int()), desc)
		case reflect.Uint, reflect.Uint64:
			flag.UintVar(fieldPtr.(*uint), fullFlagName, uint(defaultFieldValue.Int()), desc)
		case reflect.Bool:
			flag.BoolVar(fieldPtr.(*bool), fullFlagName, defaultFieldValue.Bool(), desc)
		case reflect.Float64:
			flag.Float64Var(fieldPtr.(*float64), fullFlagName, defaultFieldValue.Float(), desc)
		default:
			panic(fmt.Sprintf("不支持的类型: %s", field.Type))
		}
	}
}
