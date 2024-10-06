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
		var structTypeStrs []string
		for _, v := range typ.GetStructType().GetFields() {
			if opts.Struct == StructModeRecursiveWithName && v.GetName() != "" {
				structTypeStrs = append(structTypeStrs, fmt.Sprintf("%v %v", v.GetName(), FormatType(v.GetType(), opts)))
			} else {
				structTypeStrs = append(structTypeStrs, fmt.Sprintf("%v", FormatType(v.GetType(), opts)))
			}
		}
		return fmt.Sprintf("STRUCT<%v>", strings.Join(structTypeStrs, ", "))
	}

	if name, ok := sppb.TypeCode_name[int32(code)]; ok {
		return name
	} else {
		return "UNKNOWN"
	}
}

func FormatTypeSimple(typ *sppb.Type) string {
	return FormatType(typ, FormatOption{
		Struct: StructModeBase,
		Proto:  ProtoModeLeaf,
		Enum:   EnumModeLeaf,
		Array:  ArrayModeRecursive,
	})
}

func FormatTypeSimplest(typ *sppb.Type) string {
	return FormatType(typ, FormatOption{
		Struct: StructModeBase,
		Proto:  ProtoModeBase,
		Enum:   EnumModeBase,
		Array:  ArrayModeBase,
	})
}

func FormatTypeNormal(typ *sppb.Type) string {
	return FormatType(typ, FormatOption{
		Struct: StructModeRecursive,
		Proto:  ProtoModeLeaf,
		Enum:   EnumModeLeaf,
		Array:  ArrayModeRecursive,
	})
}

func FormatTypeVerbose(typ *sppb.Type) string {
	return FormatType(typ, FormatOption{
		Struct: StructModeRecursiveWithName,
		Proto:  ProtoModeFull,
		Enum:   EnumModeFull,
		Array:  ArrayModeRecursive,
	})
}
