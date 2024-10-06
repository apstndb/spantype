package spantype_test

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/apstndb/spantype"
	"testing"
)

func TestFormatType(t *testing.T) {
	for _, tt := range []struct {
		desc         string
		typ          *sppb.Type
		wantSimplest string
		wantSimple   string
		wantVerbose  string
		wantNormal   string
	}{
		{
			desc:         "UNKNOWN",
			typ:          &sppb.Type{Code: -1},
			wantSimplest: "-1",
			wantSimple:   "UNKNOWN",
			wantNormal:   "UNKNOWN(-1)",
			wantVerbose:  "UNKNOWN(-1)",
		},
		{
			desc:         "TYPE_CODE_UNSPECIFIED",
			typ:          &sppb.Type{Code: sppb.TypeCode_TYPE_CODE_UNSPECIFIED},
			wantSimplest: "TYPE_CODE_UNSPECIFIED",
			wantSimple:   "TYPE_CODE_UNSPECIFIED",
			wantNormal:   "TYPE_CODE_UNSPECIFIED",
			wantVerbose:  "TYPE_CODE_UNSPECIFIED",
		},
		{
			desc:         "BOOL",
			typ:          &sppb.Type{Code: sppb.TypeCode_BOOL},
			wantSimplest: "BOOL",
			wantSimple:   "BOOL",
			wantNormal:   "BOOL",
			wantVerbose:  "BOOL",
		},
		{
			desc:         "INT64",
			typ:          &sppb.Type{Code: sppb.TypeCode_INT64},
			wantSimplest: "INT64",
			wantSimple:   "INT64",
			wantNormal:   "INT64",
			wantVerbose:  "INT64",
		},
		{
			desc:         "FLOAT64",
			typ:          &sppb.Type{Code: sppb.TypeCode_FLOAT64},
			wantSimplest: "FLOAT64",
			wantSimple:   "FLOAT64",
			wantNormal:   "FLOAT64",
			wantVerbose:  "FLOAT64",
		},
		{
			desc:         "FLOAT32",
			typ:          &sppb.Type{Code: sppb.TypeCode_FLOAT32},
			wantSimplest: "FLOAT32",
			wantSimple:   "FLOAT32",
			wantNormal:   "FLOAT32",
			wantVerbose:  "FLOAT32",
		},
		{
			desc:         "TIMESTAMP",
			typ:          &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			wantSimplest: "TIMESTAMP",
			wantSimple:   "TIMESTAMP",
			wantNormal:   "TIMESTAMP",
			wantVerbose:  "TIMESTAMP",
		},
		{
			desc:         "DATE",
			typ:          &sppb.Type{Code: sppb.TypeCode_DATE},
			wantSimplest: "DATE",
			wantSimple:   "DATE",
			wantNormal:   "DATE",
			wantVerbose:  "DATE",
		},
		{
			desc:         "STRING",
			typ:          &sppb.Type{Code: sppb.TypeCode_STRING},
			wantSimplest: "STRING",
			wantSimple:   "STRING",
			wantNormal:   "STRING",
			wantVerbose:  "STRING",
		},
		{
			desc:         "BYTES",
			typ:          &sppb.Type{Code: sppb.TypeCode_BYTES},
			wantSimplest: "BYTES",
			wantSimple:   "BYTES",
			wantNormal:   "BYTES",
			wantVerbose:  "BYTES",
		},
		// ARRAY
		{
			desc: "ARRAY",
			typ: &sppb.Type{
				Code:             sppb.TypeCode_ARRAY,
				ArrayElementType: &sppb.Type{Code: sppb.TypeCode_INT64},
			},
			wantSimplest: "ARRAY",
			wantSimple:   "ARRAY<INT64>",
			wantNormal:   "ARRAY<INT64>",
			wantVerbose:  "ARRAY<INT64>",
		},
		{
			desc:         "NUMERIC",
			typ:          &sppb.Type{Code: sppb.TypeCode_NUMERIC},
			wantSimplest: "NUMERIC",
			wantSimple:   "NUMERIC",
			wantNormal:   "NUMERIC",
			wantVerbose:  "NUMERIC",
		},
		{
			desc:         "JSON",
			typ:          &sppb.Type{Code: sppb.TypeCode_JSON},
			wantSimplest: "JSON",
			wantSimple:   "JSON",
			wantNormal:   "JSON",
			wantVerbose:  "JSON",
		},
		// STRUCT
		{
			desc: "STRUCT with name",
			typ: &sppb.Type{Code: sppb.TypeCode_STRUCT, StructType: &sppb.StructType{
				Fields: []*sppb.StructType_Field{
					{
						Name: "arr",
						Type: &sppb.Type{
							Code: sppb.TypeCode_ARRAY,
							ArrayElementType: &sppb.Type{
								Code: sppb.TypeCode_STRUCT,
								StructType: &sppb.StructType{Fields: []*sppb.StructType_Field{
									{
										Name: "n",
										Type: &sppb.Type{Code: sppb.TypeCode_INT64},
									},
								}},
							},
						},
					},
				},
			}},
			wantSimplest: "STRUCT",
			wantSimple:   "STRUCT",
			wantNormal:   "STRUCT<ARRAY<STRUCT<INT64>>>",
			wantVerbose:  "STRUCT<arr ARRAY<STRUCT<n INT64>>>",
		},
		{
			desc: "STRUCT without name",
			typ: &sppb.Type{Code: sppb.TypeCode_STRUCT, StructType: &sppb.StructType{
				Fields: []*sppb.StructType_Field{
					{
						Type: &sppb.Type{
							Code: sppb.TypeCode_ARRAY,
							ArrayElementType: &sppb.Type{
								Code: sppb.TypeCode_STRUCT,
								StructType: &sppb.StructType{Fields: []*sppb.StructType_Field{
									{
										Type: &sppb.Type{Code: sppb.TypeCode_INT64},
									},
								}},
							},
						},
					},
				},
			}},
			wantSimplest: "STRUCT",
			wantSimple:   "STRUCT",
			wantNormal:   "STRUCT<ARRAY<STRUCT<INT64>>>",
			wantVerbose:  "STRUCT<ARRAY<STRUCT<INT64>>>",
		},
		// PROTO
		{
			desc: "PROTO without package",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_PROTO,
				ProtoTypeFqn: "ProtoType",
			},
			wantSimplest: "PROTO",
			wantSimple:   "ProtoType",
			wantNormal:   "ProtoType",
			wantVerbose:  "ProtoType",
		},
		{
			desc: "PROTO",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_PROTO,
				ProtoTypeFqn: "examples.ProtoType",
			},
			wantSimplest: "PROTO",
			wantSimple:   "ProtoType",
			wantNormal:   "ProtoType",
			wantVerbose:  "examples.ProtoType",
		},
		// ENUM
		{
			desc: "ENUM",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_ENUM,
				ProtoTypeFqn: "examples.EnumType",
			},
			wantSimplest: "ENUM",
			wantSimple:   "EnumType",
			wantNormal:   "EnumType",
			wantVerbose:  "examples.EnumType",
		},
		{
			desc: "ENUM without package",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_ENUM,
				ProtoTypeFqn: "EnumType",
			},
			wantSimplest: "ENUM",
			wantSimple:   "EnumType",
			wantNormal:   "EnumType",
			wantVerbose:  "EnumType",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			if got := spantype.FormatTypeSimple(tt.typ); tt.wantSimple != got {
				t.Errorf("FormatTypeSimple want: %v, got: %v", tt.wantSimple, got)
			}
			if got := spantype.FormatTypeSimplest(tt.typ); tt.wantSimplest != got {
				t.Errorf("FormatTypeSimplest want: %v, got: %v", tt.wantSimplest, got)
			}
			if got := spantype.FormatTypeNormal(tt.typ); tt.wantNormal != got {
				t.Errorf("FormatTypeNormal want: %v, got: %v", tt.wantNormal, got)
			}
			if got := spantype.FormatTypeVerbose(tt.typ); tt.wantVerbose != got {
				t.Errorf("FormatTypeVerbose want: %v, got: %v", tt.wantVerbose, got)
			}
		})
	}
}
