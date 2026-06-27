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

func TestEquivalentTypesScalarSupertypesNotEquivalent(t *testing.T) {
	t.Parallel()

	if spantype.EquivalentTypes(typector.Int64(), typector.Float64()) {
		t.Fatalf("INT64 and FLOAT64 share a GoogleSQL supertype but are not identity-equivalent")
	}
	if spantype.EquivalentTypes(typector.Int64(), typector.Numeric()) {
		t.Fatalf("INT64 and NUMERIC are coercible/supertypable but not identity-equivalent")
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

func TestEquivalentTypesArrayNegative(t *testing.T) {
	t.Parallel()

	a := typector.ElemTypeToArrayType(typector.Int64())
	b := typector.ElemTypeToArrayType(typector.String())
	if spantype.EquivalentTypes(a, b) {
		t.Fatalf("ARRAY<INT64> and ARRAY<STRING> should not be equivalent")
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

func TestEquivalentTypesStructNegative(t *testing.T) {
	t.Parallel()

	t.Run("field type mismatch", func(t *testing.T) {
		a := typector.NameTypeToStructType("a", typector.Int64())
		b := typector.NameTypeToStructType("a", typector.String())
		if spantype.EquivalentTypes(a, b) {
			t.Fatalf("STRUCT with different field types should not be equivalent")
		}
	})

	t.Run("field count mismatch", func(t *testing.T) {
		a := typector.NameTypeToStructType("a", typector.Int64())
		b := typector.MustNameCodeSlicesToStructType([]string{"a", "b"}, []sppb.TypeCode{
			sppb.TypeCode_INT64,
			sppb.TypeCode_INT64,
		})
		if spantype.EquivalentTypes(a, b) {
			t.Fatalf("STRUCT with different field counts should not be equivalent")
		}
	})

	t.Run("field order matters", func(t *testing.T) {
		a := typector.MustNameCodeSlicesToStructType([]string{"a", "b"}, []sppb.TypeCode{
			sppb.TypeCode_INT64,
			sppb.TypeCode_STRING,
		})
		b := typector.MustNameCodeSlicesToStructType([]string{"a", "b"}, []sppb.TypeCode{
			sppb.TypeCode_STRING,
			sppb.TypeCode_INT64,
		})
		if spantype.EquivalentTypes(a, b) {
			t.Fatalf("STRUCT field order should matter for equivalence")
		}
	})
}

func TestEquivalentTypesNestedArrayStructFieldNames(t *testing.T) {
	t.Parallel()

	innerA := typector.NameTypeToStructType("a", typector.Int64())
	innerB := typector.NameTypeToStructType("b", typector.Int64())
	a := typector.ElemTypeToArrayType(innerA)
	b := typector.ElemTypeToArrayType(innerB)
	if !spantype.EquivalentTypes(a, b) {
		t.Fatalf("ARRAY<STRUCT> with different field names should be equivalent")
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

func TestEquivalentTypesUnspecified(t *testing.T) {
	t.Parallel()

	unspecified := &sppb.Type{}
	if spantype.EquivalentTypes(unspecified, unspecified) {
		t.Fatalf("TYPE_CODE_UNSPECIFIED should not be equivalent")
	}
}

func TestEquivalentTypesNilAndMalformed(t *testing.T) {
	t.Parallel()

	t.Run("nil inputs", func(t *testing.T) {
		if !spantype.EquivalentTypes(nil, nil) {
			t.Error("nil types should be equivalent")
		}
		if spantype.EquivalentTypes(typector.Int64(), nil) {
			t.Error("non-nil and nil types should not be equivalent")
		}
		if spantype.EquivalentTypes(nil, typector.Int64()) {
			t.Error("nil and non-nil types should not be equivalent")
		}
	})

	t.Run("malformed scalar container fields", func(t *testing.T) {
		malformedScalar := &sppb.Type{
			Code:             sppb.TypeCode_INT64,
			ArrayElementType: typector.Int64(),
		}
		if spantype.EquivalentTypes(malformedScalar, typector.Int64()) {
			t.Error("malformed scalar with ArrayElementType should not be equivalent to valid scalar")
		}
	})

	t.Run("malformed array", func(t *testing.T) {
		missingElem := &sppb.Type{Code: sppb.TypeCode_ARRAY}
		validArray := typector.ElemTypeToArrayType(typector.Int64())
		if spantype.EquivalentTypes(missingElem, missingElem) {
			t.Error("ARRAY without element type should not be equivalent")
		}
		if spantype.EquivalentTypes(missingElem, validArray) {
			t.Error("ARRAY without element type should not be equivalent to valid array")
		}

		malformedArray := &sppb.Type{
			Code:             sppb.TypeCode_ARRAY,
			ArrayElementType: typector.Int64(),
			StructType:       &sppb.StructType{},
		}
		if spantype.EquivalentTypes(malformedArray, validArray) {
			t.Error("malformed array with StructType should not be equivalent to valid array")
		}
	})

	t.Run("malformed struct", func(t *testing.T) {
		missingStruct := &sppb.Type{Code: sppb.TypeCode_STRUCT}
		validStruct := typector.NameTypeToStructType("a", typector.Int64())
		if spantype.EquivalentTypes(missingStruct, missingStruct) {
			t.Error("STRUCT without StructType should not be equivalent")
		}
		if spantype.EquivalentTypes(missingStruct, validStruct) {
			t.Error("STRUCT without StructType should not be equivalent to valid struct")
		}

		malformedStruct := &sppb.Type{
			Code:             sppb.TypeCode_STRUCT,
			StructType:       &sppb.StructType{Fields: []*sppb.StructType_Field{typector.NameTypeToStructTypeField("a", typector.Int64())}},
			ArrayElementType: typector.Int64(),
		}
		if spantype.EquivalentTypes(malformedStruct, validStruct) {
			t.Error("malformed struct with ArrayElementType should not be equivalent to valid struct")
		}
	})

	t.Run("malformed proto and enum", func(t *testing.T) {
		emptyProto := &sppb.Type{Code: sppb.TypeCode_PROTO}
		if spantype.EquivalentTypes(emptyProto, emptyProto) {
			t.Error("PROTO without FQN should not be equivalent")
		}

		emptyEnum := &sppb.Type{Code: sppb.TypeCode_ENUM}
		if spantype.EquivalentTypes(emptyEnum, emptyEnum) {
			t.Error("ENUM without FQN should not be equivalent")
		}
	})
}
