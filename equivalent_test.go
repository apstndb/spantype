package spantype_test

import (
	"testing"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/apstndb/spantype"
	"github.com/apstndb/spantype/typector"
	"google.golang.org/protobuf/proto"
)

func TestEquivalentTypesScalar(t *testing.T) {
	t.Parallel()

	a := typector.Int64()
	b := typector.CodeToSimpleType(sppb.TypeCode_INT64)
	if !spantype.EquivalentTypes(a, b) {
		t.Fatalf("INT64 types should be equivalent")
	}
	if spantype.EquivalentTypes(a, typector.String()) {
		t.Fatalf("INT64 and STRING should not be equivalent")
	}
}

func TestEquivalentTypesArray(t *testing.T) {
	t.Parallel()

	elem := typector.Int64()
	a := typector.ElemTypeToArrayType(elem)
	b := typector.ElemTypeToArrayType(typector.CodeToSimpleType(sppb.TypeCode_INT64))
	if !spantype.EquivalentTypes(a, b) {
		t.Fatalf("ARRAY<INT64> types should be equivalent")
	}
}

func TestEquivalentTypesStructIgnoresFieldNames(t *testing.T) {
	t.Parallel()

	fieldType := typector.Int64()
	a := typector.NameTypeToStructType("a", fieldType)
	b := typector.NameTypeToStructType("b", fieldType)
	if proto.Equal(a, b) {
		t.Fatalf("proto.Equal should distinguish STRUCT field names")
	}
	if !spantype.EquivalentTypes(a, b) {
		t.Fatalf("EquivalentTypes should ignore field names")
	}
}

func TestEquivalentTypesPGNumericAnnotation(t *testing.T) {
	t.Parallel()

	if spantype.EquivalentTypes(typector.Numeric(), typector.PGNumeric()) {
		t.Fatalf("GoogleSQL NUMERIC and PG_NUMERIC should not be equivalent")
	}
}

func TestEquivalentTypesProtoFQNMismatch(t *testing.T) {
	t.Parallel()

	a := typector.FQNToProtoType("examples.Foo")
	b := typector.FQNToProtoType("examples.Bar")
	if spantype.EquivalentTypes(a, b) {
		t.Fatalf("PROTO types with different FQN should not be equivalent")
	}
}

func TestEquivalentTypesEnumFQNMismatch(t *testing.T) {
	t.Parallel()

	a := typector.FQNToEnumType("examples.Color")
	b := typector.FQNToEnumType("examples.Size")
	if spantype.EquivalentTypes(a, b) {
		t.Fatalf("ENUM types with different FQN should not be equivalent")
	}
}

func TestEquivalentTypesArrayContainerAnnotationMismatch(t *testing.T) {
	t.Parallel()

	elem := typector.Int64()
	a := &sppb.Type{
		Code:             sppb.TypeCode_ARRAY,
		ArrayElementType: elem,
	}
	b := &sppb.Type{
		Code:             sppb.TypeCode_ARRAY,
		ArrayElementType: typector.CodeToSimpleType(sppb.TypeCode_INT64),
		TypeAnnotation:   sppb.TypeAnnotationCode_PG_NUMERIC,
	}
	if spantype.EquivalentTypes(a, b) {
		t.Fatalf("ARRAY types with different container TypeAnnotation should not be equivalent")
	}
}
