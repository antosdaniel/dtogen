package generator

import (
	"go/ast"
	"strings"
)

type SrcID int

type Mapper struct {
	src SrcMappedStructs
	dst MappedStruct

	mappings Mappings
}

func (m Mapper) Src() SrcMappedStructs {
	return m.src
}

func (m Mapper) Dst() MappedStruct {
	return m.dst
}

func (m Mapper) Mappings() Mappings {
	return m.mappings
}

type SrcMappedStructs map[SrcID]MappedStruct

type MappedStruct struct {
	srcID    SrcID
	pkg      Package
	typeName string
}

func (ms MappedStruct) Pkg() Package {
	return ms.pkg
}

func (ms MappedStruct) TypeName() string {
	return ms.typeName
}

type Mappings []Mapping

type Mapping struct {
	dstField string
	src      MappingSrc
}

func (m Mapping) DstField() string {
	return m.dstField
}

func (m Mapping) Src() MappingSrc {
	return m.src
}

type MappingSrc interface {
	Type() SrcType
	SrcID() SrcID
}

type SrcType string

const (
	SrcType_Field  SrcType = "field"
	SrcType_Method SrcType = "method"
)

type MappingSrcField struct {
	srcID     SrcID
	fieldName string
}

func (MappingSrcField) Type() SrcType {
	return SrcType_Field
}

func (f MappingSrcField) SrcID() SrcID {
	return f.srcID
}

func (f MappingSrcField) FieldName() string {
	return f.fieldName
}

type MappingSrcMethod struct {
	srcID      SrcID
	methodName string
}

func (MappingSrcMethod) Type() SrcType {
	return SrcType_Method
}

func (m MappingSrcMethod) SrcID() SrcID {
	return m.srcID
}

func (m MappingSrcMethod) MethodName() string {
	return m.methodName
}

func newMapper(
	src []*ParsedStruct,
	dst *ParsedStruct,
) Mapper {
	m := Mapper{
		src: map[SrcID]MappedStruct{},
		dst: MappedStruct{pkg: dst.Package, typeName: dst.Name},
	}
	for i, s := range src {
		srcID := SrcID(i)
		m.src[srcID] = MappedStruct{srcID: srcID, pkg: s.Package, typeName: s.Name}
	}

	for _, d := range dst.Fields {
		m.mappings = append(m.mappings, newMapping(d, src[0].Fields)) // TODO: fix when we support multiple sources
	}

	return m
}

func newMapping(dstField ParsedField, srcFields ParsedFields) Mapping {
	fieldMapping, ok := findFieldMapping(dstField.Name(), srcFields)
	if ok {
		return Mapping{
			dstField: dstField.Name(),
			src:      fieldMapping,
		}
	}

	return Mapping{dstField: dstField.Name()}
}

func findFieldMapping(dstFieldName string, fields ParsedFields) (MappingSrcField, bool) {
	dst := strings.ToLower(dstFieldName)
	for _, f := range fields {
		if !ast.IsExported(f.Name()) {
			continue
		}
		if dst != strings.ToLower(f.Name()) {
			continue
		}

		return MappingSrcField{
			srcID:     0, // TODO: fix when we support multiple sources
			fieldName: f.Name(),
		}, true
	}

	return MappingSrcField{}, false
}
