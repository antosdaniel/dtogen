package _misc

// RegisteredType Example custom type, which should be used as-is on generated DTOs.
type RegisteredType struct {
	Foo string
	Bar string
}

// NonRegisteredType Example custom type, which will be assumed to exist in destination package.
type NonRegisteredType struct {
	Foo string
	Bar string
}
