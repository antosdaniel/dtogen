package input

type StructWithInternalStruct struct {
	A        string
	Internal struct {
		B string
		C string
	}
}
