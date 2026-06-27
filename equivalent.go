package spantype

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

// EquivalentTypes reports whether a and b are Spanner-equivalent [cloud.google.com/go/spanner/apiv1/spannerpb.Type]
// metadata for identity retyping: a value of either type can carry the same wire
// payload when only Type metadata differs in ways that do not affect type identity.
//
// EquivalentTypes is stricter than wire-encoding compatibility: [sppb.Type.TypeAnnotation]
// and [sppb.Type.ProtoTypeFqn] are part of type identity even when serialization is
// unchanged.
//
// Rules follow [google.spanner.v1.Type] shape
// (https://github.com/googleapis/googleapis/blob/master/google/spanner/v1/type.proto)
// and GoogleSQL supertypes for composite types
// (https://docs.cloud.google.com/spanner/docs/reference/standard-sql/conversion_rules#supertypes):
//
//   - Scalars: same TypeCode, TypeAnnotation, and ProtoTypeFqn; container fields
//     must be absent. Unknown protobuf fields are intentionally ignored.
//   - ARRAY: equivalent [sppb.Type.ArrayElementType]; [sppb.Type.StructType] must be absent.
//   - STRUCT: same number of fields with pairwise equivalent field types by
//     position; field names are not compared. [sppb.Type.ArrayElementType] must be absent.
//
// This is narrower than the Cast table in the same conversion-rules page:
// types that can CAST to one another (for example INT64 and FLOAT64) are not
// necessarily equivalent. Use EquivalentTypes only when semantics require an
// identity path (no value conversion), such as ARRAY/STRUCT layout checks or
// gcvctor.WithEquivalentType.
func EquivalentTypes(a, b *sppb.Type) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.GetCode() != b.GetCode() {
		return false
	}
	if a.GetCode() == sppb.TypeCode_TYPE_CODE_UNSPECIFIED {
		return false
	}
	if a.GetTypeAnnotation() != b.GetTypeAnnotation() {
		return false
	}
	if a.GetProtoTypeFqn() != b.GetProtoTypeFqn() {
		return false
	}
	switch a.GetCode() {
	case sppb.TypeCode_ARRAY:
		if a.GetStructType() != nil || b.GetStructType() != nil {
			return false
		}
		if a.GetArrayElementType() == nil || b.GetArrayElementType() == nil {
			return false
		}
		return EquivalentTypes(a.GetArrayElementType(), b.GetArrayElementType())
	case sppb.TypeCode_STRUCT:
		if a.GetArrayElementType() != nil || b.GetArrayElementType() != nil {
			return false
		}
		aStruct := a.GetStructType()
		bStruct := b.GetStructType()
		if aStruct == nil || bStruct == nil {
			return false
		}
		aFields := aStruct.GetFields()
		bFields := bStruct.GetFields()
		if len(aFields) != len(bFields) {
			return false
		}
		for i := range aFields {
			if aFields[i] == nil || bFields[i] == nil {
				if aFields[i] != bFields[i] {
					return false
				}
				continue
			}
			if !EquivalentTypes(aFields[i].GetType(), bFields[i].GetType()) {
				return false
			}
		}
		return true
	case sppb.TypeCode_PROTO, sppb.TypeCode_ENUM:
		if a.GetProtoTypeFqn() == "" {
			return false
		}
		if a.GetArrayElementType() != nil || b.GetArrayElementType() != nil {
			return false
		}
		if a.GetStructType() != nil || b.GetStructType() != nil {
			return false
		}
		return true
	case sppb.TypeCode_BOOL,
		sppb.TypeCode_INT64,
		sppb.TypeCode_FLOAT64,
		sppb.TypeCode_FLOAT32,
		sppb.TypeCode_TIMESTAMP,
		sppb.TypeCode_DATE,
		sppb.TypeCode_STRING,
		sppb.TypeCode_BYTES,
		sppb.TypeCode_NUMERIC,
		sppb.TypeCode_JSON,
		sppb.TypeCode_INTERVAL,
		sppb.TypeCode_UUID:
		if a.GetArrayElementType() != nil || b.GetArrayElementType() != nil {
			return false
		}
		if a.GetStructType() != nil || b.GetStructType() != nil {
			return false
		}
		return true
	default:
		return false
	}
}
