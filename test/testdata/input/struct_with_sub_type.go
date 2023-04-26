package input

type StructWithSubType struct {
	A       string
	SubType StructWithSubType_SubType
}

type StructWithSubType_SubType struct {
	AA string
}
