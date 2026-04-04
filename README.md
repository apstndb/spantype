# spantype

`github.com/apstndb/spantype` provides two related packages for working with Cloud Spanner types:

- `spantype`: format `google.spanner.v1.Type` values for logs, errors, and debugging.
- `typector`: construct `*spannerpb.Type` and `*spannerpb.StructType_Field` values for tests and helpers.

[![Go Reference](https://pkg.go.dev/badge/github.com/apstndb/spantype.svg)](https://pkg.go.dev/github.com/apstndb/spantype)

## Packages

### `spantype`

The root package formats Spanner types with configurable verbosity.

| Function | Intended use |
| --- | --- |
| `FormatTypeSimplest` | Very compact summaries such as schema overviews |
| `FormatTypeSimple` | Compact logs where array element types still matter |
| `FormatTypeNormal` | Default structured output without field names |
| `FormatTypeVerbose` | Human-facing diagnostics with struct field names |
| `FormatTypeMoreVerbose` | Errors and debugging where `PROTO` / `ENUM` kind should stay explicit |

If you need custom behavior, call `FormatType` with `FormatOption`.

### `typector`

`typector` is a constructor helper package for building Spanner type values.

- Use `CodeToSimpleType` when you already have a `sppb.TypeCode`.
- Use shorthand constructors such as `Int64()`, `String()`, and `UUID()` for common scalar types.
- Use `ElemCodeToArrayType` / `ElemTypeToArrayType` for arrays.
- Use `FQNToProtoType` / `FQNToEnumType` for `PROTO` and `ENUM`, which require a fully-qualified name.
- Prefer `...Code...` forms when your input is a type code, and `...Type...` forms when you already have `*sppb.Type`.

## CLI Example

[`./cmd/spantype`](./cmd/spantype) is a small example program that reads protobuf JSON from stdin:

```shell
echo '{"fields":[{"name":"n","type":{"code":"INT64"}}]}' | go run ./cmd/spantype --mode=verbose
```

Supported modes are `simplest`, `simple`, `normal`, `verbose`, and `more`.

## Development

```shell
go test ./...
go build ./...
go vet ./...
```
