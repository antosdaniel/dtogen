package _misc

import "time"

type Key string

type CustomType struct {
	Foo string
	Bar string
}

type Policy struct {
	ID              string
	Name            string
	AuthorFirstName string
	AuthorLastName  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
