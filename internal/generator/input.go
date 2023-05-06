package generator

type Input struct {
	Src []TypeInput
	Dst TypeInput
}

type TypeInput struct {
	ImportPath string
	Type       string
}
