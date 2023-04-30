package _misc

import "time"

type CustomType struct {
	Foo string
	Bar string
}

type Policy struct {
	ID     string
	Name   string
	Author string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
