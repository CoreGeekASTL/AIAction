// Copyright (c) Huawei Technologies Co., Ltd. 2012-2019. All rights reserved."
package flagutil

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want *flag.FlagSet
	}{
		{
			name: "",
			args: args{
				obj: &struct {
					Name string `flag:"name"`
					A    int    `flag:"a"`
					B    bool   `flag:"b"`
					C    struct {
						D string `flag:"d"`
					} `flag:"c"`
				}{
					Name: "xx",
					A:    3,
					C: struct {
						D string `flag:"d"`
					}{D: "ssss"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{"-name=123", "-c-d=456"}
			Parse(tt.args.obj)
			fmt.Println(tt.args.obj)
		})
	}
}
