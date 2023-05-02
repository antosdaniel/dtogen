package main

import (
	"testing"

	"github.com/antosdaniel/dtogen/internal/generator"
	"github.com/stretchr/testify/assert"
)

func Test_getFields(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want generator.FieldsInput
	}{
		{
			name: "selected fields",
			args: []string{"one", "two", "three"},
			want: generator.FieldsInput{
				{Name: "one"},
				{Name: "two"},
				{Name: "three"},
			},
		},
		{
			name: "field is renamed",
			args: []string{"one", "two=twooo", "three"},
			want: generator.FieldsInput{
				{Name: "one"},
				{Name: "two", RenameTo: "twooo"},
				{Name: "three"},
			},
		},
		{
			name: "field has override type",
			args: []string{"one", "two.MyType", "three"},
			want: generator.FieldsInput{
				{Name: "one"},
				{Name: "two", OverrideTypeTo: "MyType"},
				{Name: "three"},
			},
		},
		{
			name: "field is renamed and has override type",
			args: []string{"one", "two=twooo.MyType", "three"},
			want: generator.FieldsInput{
				{Name: "one"},
				{Name: "two", RenameTo: "twooo", OverrideTypeTo: "MyType"},
				{Name: "three"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFields(tt.args)

			assert.Equal(t, tt.want, result)
		})
	}
}
