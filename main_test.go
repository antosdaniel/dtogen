package main

import (
	"fmt"
	"testing"

	"github.com/antosdaniel/mappergen/internal/generator"
	"github.com/stretchr/testify/assert"
)

func Test_getInput(t *testing.T) {
	type args struct {
		src string
		dst string
		out string
	}
	tests := []struct {
		name string
		args args
		want *generator.Input
		err  error
	}{
		{
			name: "src is required",
			args: args{
				src: "",
				dst: "",
			},
			err: fmt.Errorf("src: can not be empty"),
		},
		{
			name: "dst is required",
			args: args{
				src: "src.type",
				dst: "",
			},
			err: fmt.Errorf("dst: can not be empty"),
		},
		{
			name: "src does not follow pattern",
			args: args{
				src: "type",
				dst: "dst.type",
			},
			err: fmt.Errorf("src: requires \"<import path>.<type>\" pattern"),
		},
		{
			name: "dst does not follow pattern",
			args: args{
				src: "src.type",
				dst: "type",
			},
			err: fmt.Errorf("dst: requires \"<import path>.<type>\" pattern"),
		},
		{
			name: "passes",
			args: args{
				src: "./src/path.type,pkg/src.type,github.com/pkg.type",
				dst: "dst.type",
				out: "./path/out.go",
			},
			want: &generator.Input{
				Src: []generator.TypeInput{
					{
						ImportPath: "./src/path",
						Type:       "type",
					},
					{
						ImportPath: "pkg/src",
						Type:       "type",
					},
					{
						ImportPath: "github.com/pkg",
						Type:       "type",
					},
				},
				Dst: generator.TypeInput{
					ImportPath: "dst",
					Type:       "type",
				},
				OutputPkgPath: "./path",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*src = tt.args.src
			*dst = tt.args.dst
			*out = tt.args.out
			result, err := getInput()

			if tt.err != nil {
				if assert.NotNil(t, err) {
					assert.Equal(t, tt.err.Error(), err.Error())
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}
