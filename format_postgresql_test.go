package spantype_test

import (
	"testing"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"

	"github.com/apstndb/spantype"
	"github.com/apstndb/spantype/typector"
)

func pgType(code sppb.TypeCode, ann sppb.TypeAnnotationCode) *sppb.Type {
	return &sppb.Type{Code: code, TypeAnnotation: ann}
}

func TestFormatTypePostgreSQL(t *testing.T) {
	t.Parallel()

	none := sppb.TypeAnnotationCode_TYPE_ANNOTATION_CODE_UNSPECIFIED
	tests := []struct {
		desc string
		typ  *sppb.Type
		want string
	}{
		{"nil", nil, ""},
		{"bool", pgType(sppb.TypeCode_BOOL, none), "bool"},
		{"bigint", pgType(sppb.TypeCode_INT64, none), "bigint"},
		{"oid", pgType(sppb.TypeCode_INT64, sppb.TypeAnnotationCode_PG_OID), "oid"},
		{"float4", pgType(sppb.TypeCode_FLOAT32, none), "float4"},
		{"float8", pgType(sppb.TypeCode_FLOAT64, none), "float8"},
		{"text", pgType(sppb.TypeCode_STRING, none), "text"},
		{"bytea", pgType(sppb.TypeCode_BYTES, none), "bytea"},
		{"timestamptz", pgType(sppb.TypeCode_TIMESTAMP, none), "timestamptz"},
		{"date", pgType(sppb.TypeCode_DATE, none), "date"},
		{"numeric plain", pgType(sppb.TypeCode_NUMERIC, none), "numeric"},
		{"numeric pg annotation", pgType(sppb.TypeCode_NUMERIC, sppb.TypeAnnotationCode_PG_NUMERIC), "numeric"},
		{"jsonb", pgType(sppb.TypeCode_JSON, sppb.TypeAnnotationCode_PG_JSONB), "jsonb"},
		{"plain json keeps the wire distinction", pgType(sppb.TypeCode_JSON, none), "json"},
		{"interval", pgType(sppb.TypeCode_INTERVAL, none), "interval"},
		{"uuid", pgType(sppb.TypeCode_UUID, none), "uuid"},
		{"proto never mislabeled", pgType(sppb.TypeCode_PROTO, none), "proto"},
		{"enum never mislabeled", pgType(sppb.TypeCode_ENUM, none), "enum"},
		{"unspecified", pgType(sppb.TypeCode_TYPE_CODE_UNSPECIFIED, none), "unknown"},
		{"unknown code", pgType(sppb.TypeCode(9999), none), "unknown(9999)"},
		{"array", typector.ElemCodeToArrayType(sppb.TypeCode_INT64), "bigint[]"},
		{"array of jsonb", typector.ElemTypeToArrayType(pgType(sppb.TypeCode_JSON, sppb.TypeAnnotationCode_PG_JSONB)), "jsonb[]"},
		{"nested array", typector.ElemTypeToArrayType(typector.ElemCodeToArrayType(sppb.TypeCode_STRING)), "text[][]"},
		{"array with nil element type", &sppb.Type{Code: sppb.TypeCode_ARRAY}, "array/*nil element type*/"},
		{"empty struct", &sppb.Type{Code: sppb.TypeCode_STRUCT, StructType: &sppb.StructType{}}, "STRUCT<>"},
		{
			"struct with named and unnamed fields",
			typector.StructTypeFieldsToStructType([]*sppb.StructType_Field{
				typector.NameCodeToStructTypeField("x", sppb.TypeCode_INT64),
				{Type: pgType(sppb.TypeCode_STRING, none)},
			}),
			"STRUCT<x bigint, text>",
		},
		{
			"nested struct",
			typector.StructTypeFieldsToStructType([]*sppb.StructType_Field{
				typector.NameTypeToStructTypeField("inner", typector.StructTypeFieldsToStructType([]*sppb.StructType_Field{
					typector.NameCodeToStructTypeField("x", sppb.TypeCode_INT64),
					typector.NameCodeToStructTypeField("y", sppb.TypeCode_STRING),
				})),
			}),
			"STRUCT<inner STRUCT<x bigint, y text>>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			if got := spantype.FormatTypePostgreSQL(tt.typ); got != tt.want {
				t.Errorf("FormatTypePostgreSQL() = %q, want %q", got, tt.want)
			}
		})
	}
}
