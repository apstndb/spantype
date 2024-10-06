package spantype

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"fmt"
	"strings"
)

type StructMode int

const (
	StructModeBase StructMode = iota
	StructModeRecursive
	StructModeRecursiveWithName
)

type ProtoMode int

const (
	ProtoModeBase ProtoMode = iota
	ProtoModeLeaf
	ProtoModeFull
)

type EnumMode int

const (
	EnumModeBase EnumMode = iota
	EnumModeLeaf
	EnumModeFull
)

type ArrayMode int

const (
	ArrayModeBase ArrayMode = iota
	ArrayModeRecursive
)

type FormatOption struct {
	Struct StructMode
	Proto  ProtoMode
	Enum   EnumMode
	Array  ArrayMode
}

var (
	FormatTypeOptionSimple = FormatOption{
		Struct: StructModeBase,
		Proto:  ProtoModeLeaf,
		Enum:   EnumModeLeaf,
		Array:  ArrayModeRecursive,
	}
	FormatTypeOptionSimplest = FormatOption{
		Struct: StructModeBase,
		Proto:  ProtoModeBase,
		Enum:   EnumModeBase,
		Array:  ArrayModeBase,
	}
	FormatOptionNormal = FormatOption{
		Struct: StructModeRecursive,
		Proto:  ProtoModeLeaf,
		Enum:   EnumModeLeaf,
		Array:  ArrayModeRecursive,
	}
	FormatTypeOptionVerbose = FormatOption{
		Struct: StructModeRecursiveWithName,
		Proto:  ProtoModeFull,
		Enum:   EnumModeFull,
		Array:  ArrayModeRecursive,
	}
)

type FormatRowOption struct {
	TypeOption FormatOption
}

func lastCut(s, sep string) (before string, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func FormatType(typ *sppb.Type, opts FormatOption) string {
	code := typ.GetCode()
	switch code {
	case sppb.TypeCode_ARRAY:
		if opts.Array == ArrayModeBase {
			break
		}
		return fmt.Sprintf("ARRAY<%v>", FormatType(typ.GetArrayElementType(), opts))
	case sppb.TypeCode_PROTO:
		switch opts.Proto {
		case ProtoModeBase:
			break
		case ProtoModeLeaf:
			fqn := typ.GetProtoTypeFqn()
			_, after, _ := lastCut(fqn, ".")
			return after
		case ProtoModeFull:
			return typ.GetProtoTypeFqn()
		}
	case sppb.TypeCode_ENUM:
		switch opts.Enum {
		case EnumModeBase:
			break
		case EnumModeLeaf:
			fqn := typ.GetProtoTypeFqn()
			_, after, _ := lastCut(fqn, ".")
			return after
		case EnumModeFull:
			return typ.GetProtoTypeFqn()
		}
	case sppb.TypeCode_STRUCT:
		if opts.Struct == StructModeBase {
			break
		}
		return fmt.Sprintf("STRUCT<%v>", FormatStructFields(typ.GetStructType().GetFields(), opts))
	}

	if name, ok := sppb.TypeCode_name[int32(code)]; ok {
		return name
	} else {
		return "UNKNOWN"
	}
}

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

func FormatTypeSimple(typ *sppb.Type) string {
	return FormatType(typ, FormatTypeOptionSimple)
}

func FormatTypeSimplest(typ *sppb.Type) string {
	return FormatType(typ, FormatTypeOptionSimplest)
}

func FormatTypeNormal(typ *sppb.Type) string {
	return FormatType(typ, FormatOptionNormal)
}

func FormatTypeVerbose(typ *sppb.Type) string {
	return FormatType(typ, FormatTypeOptionVerbose)
}
