package spantype

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/protobuf/proto"
)

// EquivalentTypes reports whether a and b are Spanner-equivalent type metadata.
// Scalar types require proto.Equal metadata. ARRAY types require equivalent
// element types. STRUCT types require the same number of fields with pairwise
// equivalent field types; field names are not compared. This matches identity
// CAST and ARRAY cast equivalence in GoogleSQL semantic layers.
func EquivalentTypes(a, b *sppb.Type) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if a.GetCode() != b.GetCode() {
		return false
	}
	switch a.GetCode() {
	case sppb.TypeCode_ARRAY:
		return EquivalentTypes(a.GetArrayElementType(), b.GetArrayElementType())
	case sppb.TypeCode_STRUCT:
		aStruct := a.GetStructType()
		bStruct := b.GetStructType()
		if aStruct == nil || bStruct == nil {
			return aStruct == nil && bStruct == nil
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
		return proto.Equal(a, b)
	}
}
