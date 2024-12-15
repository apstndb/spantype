package typector

import (
	"fmt"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

func CodeToSimpleType(code sppb.TypeCode) *sppb.Type {
	return &sppb.Type{Code: code}
}

func ElemCodeToArrayType(code sppb.TypeCode) *sppb.Type {
	return ElemTypeToArrayType(CodeToSimpleType(code))
}

func ElemTypeToArrayType(typ *sppb.Type) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_ARRAY, ArrayElementType: typ}
}

func StructTypeFieldsToStructType(fields []*sppb.StructType_Field) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_STRUCT, StructType: &sppb.StructType{Fields: fields}}
}

func FQNToProtoType(fqn string) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_PROTO, ProtoTypeFqn: fqn}
}

func FQNToEnumType(fqn string) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_ENUM, ProtoTypeFqn: fqn}
}

func NameCodeToStructType(name string, code sppb.TypeCode) *sppb.Type {
	return NameTypeToStructType(name, CodeToSimpleType(code))
}

func NameTypeToStructType(name string, typ *sppb.Type) *sppb.Type {
	return StructTypeFieldsToStructType([]*sppb.StructType_Field{
		NameTypeToStructTypeField(name, typ),
	})
}

func NameCodeToStructTypeField(name string, code sppb.TypeCode) *sppb.StructType_Field {
	return NameTypeToStructTypeField(name, CodeToSimpleType(code))
}

func NameTypeToStructTypeField(name string, typ *sppb.Type) *sppb.StructType_Field {
	return &sppb.StructType_Field{Name: name, Type: typ}
}

func CodeToUnnamedStructTypeField(code sppb.TypeCode) *sppb.StructType_Field {
	return NameTypeToStructTypeField("", CodeToSimpleType(code))
}

func TypeToUnnamedStructTypeField(typ *sppb.Type) *sppb.StructType_Field {
	return &sppb.StructType_Field{Type: typ}
}

func NameTypeSlicesToStructType(names []string, types []*sppb.Type) (*sppb.Type, error) {
	fields, err := NameTypeSlicesToStructTypeFields(names, types)
	if err != nil {
		return nil, err
	}
	return StructTypeFieldsToStructType(fields), nil
}

func MustNameTypeSlicesToStructType(names []string, types []*sppb.Type) *sppb.Type {
	return must(NameTypeSlicesToStructType(names, types))
}

func NameTypeSlicesToStructTypeFields(names []string, types []*sppb.Type) ([]*sppb.StructType_Field, error) {
	if len(names) != len(types) {
		return nil, fmt.Errorf("length mismatch: len(names)=%d, len(types)=%d", len(names), len(types))
	}

	var fields []*sppb.StructType_Field
	for i := range names {
		fields = append(fields, NameTypeToStructTypeField(names[i], types[i]))
	}
	return fields, nil
}

func MustNameTypeSlicesToStructTypeFields(names []string, types []*sppb.Type) []*sppb.StructType_Field {
	return must(NameTypeSlicesToStructTypeFields(names, types))
}

func NameCodeSlicesToStructType(names []string, codes []sppb.TypeCode) (*sppb.Type, error) {
	fields, err := NameCodeSlicesToStructTypeFields(names, codes)
	if err != nil {
		return nil, err
	}
	return StructTypeFieldsToStructType(fields), nil
}

func MustNameCodeSlicesToStructType(names []string, codes []sppb.TypeCode) *sppb.Type {
	return must(NameCodeSlicesToStructType(names, codes))
}

func NameCodeSlicesToStructTypeFields(names []string, codes []sppb.TypeCode) ([]*sppb.StructType_Field, error) {
	var types []*sppb.Type
	for _, code := range codes {
		types = append(types, CodeToSimpleType(code))
	}

	return NameTypeSlicesToStructTypeFields(names, types)
}

func MustNameCodeSlicesToStructTypeFields(names []string, codes []sppb.TypeCode) []*sppb.StructType_Field {
	return must(NameCodeSlicesToStructTypeFields(names, codes))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
