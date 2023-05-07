package generator

type Input struct {
	Src           []TypeInput
	Dst           TypeInput
	OutputPkgPath string
}

type TypeInput struct {
	ImportPath string
	Type       string
}
