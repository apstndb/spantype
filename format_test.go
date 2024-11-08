package spantype_test

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/apstndb/spantype"
	"testing"
)

func TestFormatType(t *testing.T) {
	for _, tt := range []struct {
		desc            string
		typ             *sppb.Type
		wantSimplest    string
		wantSimple      string
		wantVerbose     string
		wantMoreVerbose string
		wantNormal      string
	}{
		{
			desc:            "UNKNOWN",
			typ:             &sppb.Type{Code: -1},
			wantSimplest:    "-1",
			wantSimple:      "UNKNOWN",
			wantNormal:      "UNKNOWN(-1)",
			wantVerbose:     "UNKNOWN(-1)",
			wantMoreVerbose: "UNKNOWN(-1)",
		},
		{
			desc:            "TYPE_CODE_UNSPECIFIED",
			typ:             &sppb.Type{Code: sppb.TypeCode_TYPE_CODE_UNSPECIFIED},
			wantSimplest:    "TYPE_CODE_UNSPECIFIED",
			wantSimple:      "TYPE_CODE_UNSPECIFIED",
			wantNormal:      "TYPE_CODE_UNSPECIFIED",
			wantVerbose:     "TYPE_CODE_UNSPECIFIED",
			wantMoreVerbose: "TYPE_CODE_UNSPECIFIED",
		},
		{
			desc:            "BOOL",
			typ:             &sppb.Type{Code: sppb.TypeCode_BOOL},
			wantSimplest:    "BOOL",
			wantSimple:      "BOOL",
			wantNormal:      "BOOL",
			wantVerbose:     "BOOL",
			wantMoreVerbose: "BOOL",
		},
		{
			desc:            "INT64",
			typ:             &sppb.Type{Code: sppb.TypeCode_INT64},
			wantSimplest:    "INT64",
			wantSimple:      "INT64",
			wantNormal:      "INT64",
			wantVerbose:     "INT64",
			wantMoreVerbose: "INT64",
		},
		{
			desc:            "FLOAT64",
			typ:             &sppb.Type{Code: sppb.TypeCode_FLOAT64},
			wantSimplest:    "FLOAT64",
			wantSimple:      "FLOAT64",
			wantNormal:      "FLOAT64",
			wantVerbose:     "FLOAT64",
			wantMoreVerbose: "FLOAT64",
		},
		{
			desc:            "FLOAT32",
			typ:             &sppb.Type{Code: sppb.TypeCode_FLOAT32},
			wantSimplest:    "FLOAT32",
			wantSimple:      "FLOAT32",
			wantNormal:      "FLOAT32",
			wantVerbose:     "FLOAT32",
			wantMoreVerbose: "FLOAT32",
		},
		{
			desc:            "TIMESTAMP",
			typ:             &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			wantSimplest:    "TIMESTAMP",
			wantSimple:      "TIMESTAMP",
			wantNormal:      "TIMESTAMP",
			wantVerbose:     "TIMESTAMP",
			wantMoreVerbose: "TIMESTAMP",
		},
		{
			desc:            "DATE",
			typ:             &sppb.Type{Code: sppb.TypeCode_DATE},
			wantSimplest:    "DATE",
			wantSimple:      "DATE",
			wantNormal:      "DATE",
			wantVerbose:     "DATE",
			wantMoreVerbose: "DATE",
		},
		{
			desc:            "STRING",
			typ:             &sppb.Type{Code: sppb.TypeCode_STRING},
			wantSimplest:    "STRING",
			wantSimple:      "STRING",
			wantNormal:      "STRING",
			wantVerbose:     "STRING",
			wantMoreVerbose: "STRING",
		},
		{
			desc:            "BYTES",
			typ:             &sppb.Type{Code: sppb.TypeCode_BYTES},
			wantSimplest:    "BYTES",
			wantSimple:      "BYTES",
			wantNormal:      "BYTES",
			wantVerbose:     "BYTES",
			wantMoreVerbose: "BYTES",
		},
		// ARRAY
		{
			desc: "ARRAY",
			typ: &sppb.Type{
				Code:             sppb.TypeCode_ARRAY,
				ArrayElementType: &sppb.Type{Code: sppb.TypeCode_INT64},
			},
			wantSimplest:    "ARRAY",
			wantSimple:      "ARRAY<INT64>",
			wantNormal:      "ARRAY<INT64>",
			wantVerbose:     "ARRAY<INT64>",
			wantMoreVerbose: "ARRAY<INT64>",
		},
		{
			desc:            "NUMERIC",
			typ:             &sppb.Type{Code: sppb.TypeCode_NUMERIC},
			wantSimplest:    "NUMERIC",
			wantSimple:      "NUMERIC",
			wantNormal:      "NUMERIC",
			wantVerbose:     "NUMERIC",
			wantMoreVerbose: "NUMERIC",
		},
		{
			desc:            "JSON",
			typ:             &sppb.Type{Code: sppb.TypeCode_JSON},
			wantSimplest:    "JSON",
			wantSimple:      "JSON",
			wantNormal:      "JSON",
			wantVerbose:     "JSON",
			wantMoreVerbose: "JSON",
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
			wantSimplest:    "STRUCT",
			wantSimple:      "STRUCT",
			wantNormal:      "STRUCT<ARRAY<STRUCT<INT64>>>",
			wantVerbose:     "STRUCT<arr ARRAY<STRUCT<n INT64>>>",
			wantMoreVerbose: "STRUCT<arr ARRAY<STRUCT<n INT64>>>",
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
			wantSimplest:    "STRUCT",
			wantSimple:      "STRUCT",
			wantNormal:      "STRUCT<ARRAY<STRUCT<INT64>>>",
			wantVerbose:     "STRUCT<ARRAY<STRUCT<INT64>>>",
			wantMoreVerbose: "STRUCT<ARRAY<STRUCT<INT64>>>",
		},
		// PROTO
		{
			desc: "PROTO without package",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_PROTO,
				ProtoTypeFqn: "ProtoType",
			},
			wantSimplest:    "PROTO",
			wantSimple:      "ProtoType",
			wantNormal:      "ProtoType",
			wantVerbose:     "ProtoType",
			wantMoreVerbose: "PROTO<ProtoType>",
		},
		{
			desc: "PROTO",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_PROTO,
				ProtoTypeFqn: "examples.ProtoType",
			},
			wantSimplest:    "PROTO",
			wantSimple:      "ProtoType",
			wantNormal:      "ProtoType",
			wantVerbose:     "examples.ProtoType",
			wantMoreVerbose: "PROTO<examples.ProtoType>",
		},
		// ENUM
		{
			desc: "ENUM",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_ENUM,
				ProtoTypeFqn: "examples.EnumType",
			},
			wantSimplest:    "ENUM",
			wantSimple:      "EnumType",
			wantNormal:      "EnumType",
			wantVerbose:     "examples.EnumType",
			wantMoreVerbose: "ENUM<examples.EnumType>",
		},
		{
			desc: "ENUM without package",
			typ: &sppb.Type{
				Code:         sppb.TypeCode_ENUM,
				ProtoTypeFqn: "EnumType",
			},
			wantSimplest:    "ENUM",
			wantSimple:      "EnumType",
			wantNormal:      "EnumType",
			wantVerbose:     "EnumType",
			wantMoreVerbose: "ENUM<EnumType>",
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
			if got := spantype.FormatTypeMoreVerbose(tt.typ); tt.wantMoreVerbose != got {
				t.Errorf("FormatTypeMoreVerbose want: %v, got: %v", tt.wantMoreVerbose, got)
			}
		})
	}
}

func TestFormatTypeCode(t *testing.T) {
	tests := []struct {
		desc        string
		code        sppb.TypeCode
		want        string
		shouldPanic bool
		mode        spantype.UnknownMode
	}{
		{
			desc:        "UNKNOWN should panic",
			code:        -1,
			shouldPanic: true,
			mode:        spantype.UnknownModePanic,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			defer func() {
				if rec := recover(); rec != nil && !tt.shouldPanic {
					t.Errorf("FormatTypeCode should not panic: %v", rec)
				}
			}()
			if got := spantype.FormatTypeCode(tt.code, tt.mode); tt.want != got {
				t.Errorf("FormatTypeCode want: %v, got: %v", tt.want, got)
			}
			if tt.shouldPanic {
				t.Errorf("FormatTypeCode should panic")
			}
		})
	}
}
