package spantype

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

// EquivalentTypes reports whether a and b are Spanner-equivalent [cloud.google.com/go/spanner/apiv1/spannerpb.Type]
// metadata for identity retyping: a value of either type can carry the same wire
// payload when only Type metadata differs.
//
// Rules follow [google.spanner.v1.Type] shape
// (https://github.com/googleapis/googleapis/blob/master/google/spanner/v1/type.proto)
// and GoogleSQL supertypes for composite types
// (https://docs.cloud.google.com/spanner/docs/reference/standard-sql/conversion_rules#supertypes):
//
//   - Scalars: proto.Equal, including [sppb.Type.TypeAnnotation] and
//     [sppb.Type.ProtoTypeFqn] for PROTO and ENUM.
//   - ARRAY: equivalent [sppb.Type.ArrayElementType].
//   - STRUCT: same number of fields with pairwise equivalent field types by
//     position; field names are not compared.
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
	if a.GetTypeAnnotation() != b.GetTypeAnnotation() {
		return false
	}
	if a.GetProtoTypeFqn() != b.GetProtoTypeFqn() {
		return false
	}
	switch a.GetCode() {
	case sppb.TypeCode_ARRAY:
		return EquivalentTypes(a.GetArrayElementType(), b.GetArrayElementType())
	case sppb.TypeCode_STRUCT:
		aStruct := a.GetStructType()
		bStruct := b.GetStructType()
		if aStruct == nil || bStruct == nil {
			return aStruct == bStruct
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
	default:
		// Code, TypeAnnotation, and ProtoTypeFqn already match. For scalar
		// codes, reject malformed container fields without proto.Equal so
		// unknown protobuf fields from newer servers do not break identity checks.
		if a.GetArrayElementType() != nil || b.GetArrayElementType() != nil {
			return false
		}
		if a.GetStructType() != nil || b.GetStructType() != nil {
			return false
		}
		return true
	}
}
