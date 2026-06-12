package spantype

import (
	"fmt"
	"strings"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

// FormatTypePostgreSQL renders a [sppb.Type] using PostgreSQL-dialect
// spellings from
// https://docs.cloud.google.com/spanner/docs/reference/postgresql/data-types
// (supported types: bool, bytea, date, float4, float8, bigint, interval,
// jsonb, numeric, timestamptz, text, uuid, oid, and `T[]` array
// declarations). [sppb.TypeCode_JSON] is rendered "jsonb" only when the type
// carries [sppb.TypeAnnotationCode_PG_JSONB] (otherwise "json", a label that
// preserves the wire distinction even though plain JSON is not a PostgreSQL
// data type), and [sppb.TypeCode_INT64] with
// [sppb.TypeAnnotationCode_PG_OID] is rendered "oid". A nil type renders as
// the empty string.
//
// [sppb.TypeCode_PROTO] and [sppb.TypeCode_ENUM] are not
// PostgreSQL-interface types: they render as "proto" and "enum" — not as
// "bytea" or "text" — so wire types are never mislabeled if they appear in
// metadata. [sppb.TypeCode_STRUCT] does not appear in PostgreSQL-dialect
// metadata in practice today; the provisional rendering uses GoogleSQL
// STRUCT<...> declaration syntax with field types spelled using the
// PostgreSQL names above, and may change once official guidance for
// composite types exists (no stability guarantee for the STRUCT branch).
//
// This function originated as FormatPostgreSQLType in
// github.com/apstndb/spanpg, where its spellings are pinned against
// PostgreSQL-dialect Spanner databases by integration probes.
func FormatTypePostgreSQL(typ *sppb.Type) string {
	if typ == nil {
		return ""
	}
	return formatTypePostgreSQLImpl(typ)
}

func formatTypePostgreSQLImpl(typ *sppb.Type) string {
	// Recursive positions (array elements, struct field types) may pass nil;
	// the explicit check keeps the behavior independent of getter nil-safety.
	if typ == nil {
		return "unknown"
	}
	switch typ.GetCode() {
	case sppb.TypeCode_ARRAY:
		elem := typ.GetArrayElementType()
		if elem == nil {
			return "array/*nil element type*/"
		}
		return formatTypePostgreSQLImpl(elem) + "[]"

	case sppb.TypeCode_STRUCT:
		// Provisional: PG dialect rarely exposes STRUCT today; see FormatTypePostgreSQL godoc.
		fields := typ.GetStructType().GetFields()
		if len(fields) == 0 {
			return "STRUCT<>"
		}
		return "STRUCT<" + formatStructFieldsWith(fields, formatTypePostgreSQLImpl) + ">"

	case sppb.TypeCode_BOOL:
		return "bool"

	case sppb.TypeCode_INT64:
		if typ.GetTypeAnnotation() == sppb.TypeAnnotationCode_PG_OID {
			return "oid"
		}
		return "bigint"

	case sppb.TypeCode_FLOAT32:
		return "float4"

	case sppb.TypeCode_FLOAT64:
		return "float8"

	case sppb.TypeCode_STRING:
		return "text"

	case sppb.TypeCode_BYTES:
		return "bytea"

	case sppb.TypeCode_TIMESTAMP:
		return "timestamptz"

	case sppb.TypeCode_DATE:
		return "date"

	case sppb.TypeCode_NUMERIC:
		return "numeric"

	case sppb.TypeCode_JSON:
		if typ.GetTypeAnnotation() == sppb.TypeAnnotationCode_PG_JSONB {
			return "jsonb"
		}
		return "json"

	case sppb.TypeCode_INTERVAL:
		return "interval"

	case sppb.TypeCode_UUID:
		return "uuid"

	case sppb.TypeCode_PROTO:
		return "proto"

	case sppb.TypeCode_ENUM:
		return "enum"

	case sppb.TypeCode_TYPE_CODE_UNSPECIFIED:
		return "unknown"

	default:
		return fmt.Sprintf("unknown(%d)", typ.GetCode())
	}
}

// formatStructFieldsWith joins STRUCT fields as `name type, ...` (unnamed
// fields render the type only), spelling each field type with format.
// Unlike [FormatStructFields] it always shows available field names.
func formatStructFieldsWith(fields []*sppb.StructType_Field, format func(*sppb.Type) string) string {
	var b strings.Builder
	for i, f := range fields {
		if i > 0 {
			b.WriteString(", ")
		}
		if n := f.GetName(); n != "" {
			b.WriteString(n)
			b.WriteByte(' ')
		}
		b.WriteString(format(f.GetType()))
	}
	return b.String()
}
