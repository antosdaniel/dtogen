package _misc

// RegisteredValueType Example custom type, which should be used as-is on generated DTOs.
type RegisteredValueType struct {
	Foo string
	Bar string
}

// NonRegisteredValueType Example custom type, which will be assumed to exist in destination package.
type NonRegisteredValueType struct {
	Foo string
	Bar string
}
