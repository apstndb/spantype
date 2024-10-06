package spantype

import (
	"fmt"
	"strings"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

type StructMode int

const (
	// StructModeBase formats `STRUCT` type as `STRUCT`
	StructModeBase StructMode = iota
	// StructModeRecursive formats `STRUCT` type with field types. e.g. `STRUCT<INT64, STRUCT<INT64>>`
	StructModeRecursive
	// StructModeRecursiveWithName formats `STRUCT` type with field types with field name. e.g. `STRUCT<n INT64, s STRUCT<n INT64>>`
	StructModeRecursiveWithName
)

type ProtoEnumMode int

const (
	// ProtoEnumModeBase formats `PROTO` and `ENUM` type as `PROTO` and `ENUM`
	ProtoEnumModeBase ProtoEnumMode = iota
	// ProtoEnumModeLeaf formats `PROTO` and `ENUM` type without package name. e.g. `ProtoType`, `EnumType`
	ProtoEnumModeLeaf
	// ProtoEnumModeFull formats `PROTO` and `ENUM` type as full qualified name. e.g. `examples.ProtoType`, `examples.EnumType`
	ProtoEnumModeFull
)

type ArrayMode int

const (
	// ArrayModeBase formats `ARRAY` type as `ARRAY.`
	ArrayModeBase ArrayMode = iota
	// ArrayModeRecursive formats `ARRAY` type with element type. e.g. `ARRAY<INT64>`
	ArrayModeRecursive
)

// FormatOption is a option for FormatType, and FormatStructFields.
type FormatOption struct {
	Struct StructMode
	Proto  ProtoEnumMode
	Enum   ProtoEnumMode
	Array  ArrayMode
}

var (
	// FormatOptionSimplest is a FormatOption for FormatTypeSimplest.
	FormatOptionSimplest = FormatOption{
		Struct: StructModeBase,
		Proto:  ProtoEnumModeBase,
		Enum:   ProtoEnumModeBase,
		Array:  ArrayModeBase,
	}
	// FormatOptionSimple is a FormatOption for FormatTypeSimple.
	FormatOptionSimple = FormatOption{
		Struct: StructModeBase,
		Proto:  ProtoEnumModeLeaf,
		Enum:   ProtoEnumModeLeaf,
		Array:  ArrayModeRecursive,
	}
	// FormatOptionNormal is a FormatOption for FormatTypeNormal.
	FormatOptionNormal = FormatOption{
		Struct: StructModeRecursive,
		Proto:  ProtoEnumModeLeaf,
		Enum:   ProtoEnumModeLeaf,
		Array:  ArrayModeRecursive,
	}
	// FormatOptionVerbose is a FormatOption for FormatTypeVerbose.
	FormatOptionVerbose = FormatOption{
		Struct: StructModeRecursiveWithName,
		Proto:  ProtoEnumModeFull,
		Enum:   ProtoEnumModeFull,
		Array:  ArrayModeRecursive,
	}
)

func lastCut(s, sep string) (before string, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return "", s, false
}

// FormatType formats Cloud Spanner type using the given FormatOption.
func FormatType(typ *sppb.Type, opts FormatOption) string {
	code := typ.GetCode()
	switch code {
	case sppb.TypeCode_ARRAY:
		if opts.Array == ArrayModeBase {
			break
		}
		return fmt.Sprintf("ARRAY<%v>", FormatType(typ.GetArrayElementType(), opts))
	case sppb.TypeCode_PROTO:
		return FormatProtoEnum(typ, opts.Proto)
	case sppb.TypeCode_ENUM:
		return FormatProtoEnum(typ, opts.Enum)
	case sppb.TypeCode_STRUCT:
		if opts.Struct == StructModeBase {
			break
		}
		return fmt.Sprintf("STRUCT<%v>", FormatStructFields(typ.GetStructType().GetFields(), opts))
	}

	return FormatTypeCode(code)
}

// FormatProtoEnum formats `PROTO` or `ENUM` type using ProtoEnumMode.
// It panics when the input type is not `PROTO` or `ENUM`.
func FormatProtoEnum(typ *sppb.Type, mode ProtoEnumMode) string {
	if typ.GetCode() != sppb.TypeCode_PROTO && typ.GetCode() != sppb.TypeCode_ENUM {
		panic(fmt.Sprintf("precondition failed: TypeCode must be PROTO or ENUM, but %v", typ))
	}

	switch mode {
	case ProtoEnumModeLeaf:
		_, after, _ := lastCut(typ.GetProtoTypeFqn(), ".")
		return after
	case ProtoEnumModeFull:
		return typ.GetProtoTypeFqn()
	default:
		return FormatTypeCode(typ.GetCode())
	}
}

// FormatTypeCode formats sppb.TypeCode, but it formats unknown type code as `UNKNOWN(int32(code))`. e.g. `UNKNOWN(-1)`
func FormatTypeCode(code sppb.TypeCode) string {
	if name, ok := sppb.TypeCode_name[int32(code)]; ok {
		return name
	} else {
		return fmt.Sprintf("UNKNOWN(%v)", int32(code))
	}
}

// FormatStructFields formats Cloud Spanner struct fields or `metadata.rowType` using the given FormatOption.
func FormatStructFields(fields []*sppb.StructType_Field, opts FormatOption) string {
	var fieldsStr []string
	for _, field := range fields {
		typeStr := FormatType(field.GetType(), opts)
		if opts.Struct == StructModeRecursiveWithName && field.GetName() != "" {
			fieldsStr = append(fieldsStr, fmt.Sprintf("%v %v", field.GetName(), typeStr))
		} else {
			fieldsStr = append(fieldsStr, fmt.Sprintf("%v", typeStr))
		}
	}
	return strings.Join(fieldsStr, ", ")
}

// FormatTypeSimplest formats Cloud Spanner type as simplest format.
// e.g. `INT64`, `ARRAY, `PROTO`, `ENUM` `STRUCT`
func FormatTypeSimplest(typ *sppb.Type) string {
	return FormatType(typ, FormatOptionSimplest)
}

// FormatTypeSimple formats Cloud Spanner type as simple format.
// e.g. `INT64`, `ARRAY<INT64>`, `ProtoType`, `EnumType` `STRUCT`
func FormatTypeSimple(typ *sppb.Type) string {
	return FormatType(typ, FormatOptionSimple)
}

// FormatTypeNormal formats Cloud Spanner type as normal format.
// e.g. `INT64`, `ARRAY<INT64>`, `ProtoType`, `EnumType`, `STRUCT<INT64>`
func FormatTypeNormal(typ *sppb.Type) string {
	return FormatType(typ, FormatOptionNormal)
}

// FormatTypeVerbose formats Cloud Spanner type as verbose format.
// e.g. `INT64`, `ARRAY<INT64>`, `examples.ProtoType`, `examples.EnumType`, `STRUCT<n INT64>`
func FormatTypeVerbose(typ *sppb.Type) string {
	return FormatType(typ, FormatOptionVerbose)
}
