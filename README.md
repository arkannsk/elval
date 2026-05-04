# ElVal — Lightning Fast Go Validator with Code Generation & OpenAPI Support

[![Go Version](https://img.shields.io/github/go-mod/go-version/arkannsk/elval)](https://golang.org/)
[![License](https://img.shields.io/github/license/arkannsk/elval)](LICENSE)
[![Benchmarks](https://img.shields.io/badge/benchmarks-440_allocs%2Fop-brightgreen)](BENCHMARKS.md)

[//]: # ([![Go Reference]&#40;https://pkg.go.dev/badge/github.com/arkannsk/elval.svg&#41;]&#40;https://pkg.go.dev/github.com/arkannsk/elval&#41;)

[//]: # ([![Go Report Card]&#40;https://goreportcard.com/badge/github.com/arkannsk/elval&#41;]&#40;https://goreportcard.com/report/github.com/arkannsk/elval&#41;)

**ElVal** is a code generation-based validator for Go that eliminates reflection overhead entirely. By generating
type-safe validation code at build time, ElVal achieves **6x faster validation** with **zero memory allocations** at
runtime.

Additionally, ElVal generates **OpenAPI 3.0 schemas** directly from your struct annotations, supporting complex types,
external libraries (via stubs), polymorphism, custom type rewriting, and file/stream uploads.

## Table of Contents

* [ElVal — Lightning Fast Go Validator with Code Generation &amp; OpenAPI Support](#elval--lightning-fast-go-validator-with-code-generation--openapi-support)
    * [Table of Contents](#table-of-contents)
    * [Features](#features)
    * [Quick Start](#quick-start)
        * [1\. Define your struct with annotations](#1-define-your-struct-with-annotations)
        * [2\. Generate validation code](#2-generate-validation-code)
        * [3\. Use the generated Validate method](#3-use-the-generated-validate-method)
    * [Commands](#commands)
        * [generate — Generate validation code](#generate--generate-validation-code)
        * [lint — Validate annotations](#lint--validate-annotations)
        * [Go Generate Integration](#go-generate-integration)
    * [Diagnostics &amp; Linting](#diagnostics--linting)
        * [Output Format](#output-format)
        * [Output Control Flags](#output-control-flags)
        * [Example: go generate with diagnostics](#example-go-generate-with-diagnostics)
    * [Annotations](#annotations)
        * [Required &amp; Optional](#required--optional)
        * [String Validators](#string-validators)
        * [Numeric Validators](#numeric-validators)
        * [Comparison Validators](#comparison-validators)
        * [Enum Validators](#enum-validators)
        * [Date Validators (time\.Time)](#date-validators-timetime)
        * [Duration Validators (time\.Duration)](#duration-validators-timeduration)
        * [Slice Validators](#slice-validators)
        * [URL Validators](#url-validators)
    * [File &amp; Stream Uploads](#file--stream-uploads)
        * [Auto\-detection of Standard Types](#auto-detection-of-standard-types)
        * [Custom Types: Explicit Annotations](#custom-types-explicit-annotations)
        * [Ignoring Structs](#ignoring-structs)
    * [Nested Structures](#nested-structures)
    * [Custom Validators](#custom-validators)
        * [Register validator](#register-validator)
        * [Use in struct](#use-in-struct)
    * [Decorators](#decorators)
    * [OpenAPI](#openapi)
        * [OpenAPI Annotations](#openapi-annotations)
            * [External Types &amp; Stubs](#external-types--stubs)
            * [Reference the Stub in your API Struct:](#reference-the-stub-in-your-api-struct)
            * [Polymorphism (Discriminator &amp; OneOf):](#polymorphism-discriminator--oneof)
            * [Type Rewriting](#type-rewriting)
            * [File &amp; Stream Fields](#file--stream-fields)
            * [Nested Structures &amp; Generics](#nested-structures--generics)
            * [Structs Without Body Fields](#structs-without-body-fields)
    * [Performance](#performance)

## Features

- 🚀 **6x faster** than traditional validators
- 💾 **Zero memory allocations** at runtime
- 🔧 **Code generation** - no reflection overhead
- 📝 **Simple annotations** in comments
- 🎯 **Type-safe** validation
- 🔌 **Extensible** with custom validators
- 🧩 **Decorators** for auto-populating fields
- 📦 **No external dependencies** for validation
- 📁 **File & stream support** with OpenAPI auto-mapping

## Quick Start

```bash
go install github.com/arkannsk/elval/cmd/elval-gen@latest
```

### 1. Define your struct with annotations

```go
package user

//go:generate elval-gen generate -input . -openapi

type User struct {
	// @evl:validate required
	// @evl:validate min:2
	// @evl:validate max:50
	Name string

	// @evl:validate required
	// @evl:validate pattern:email
	Email string

	// @evl:validate min:18
	// @evl:validate max:120
	Age int

	// @evl:validate optional
	// @evl:validate pattern:phone
	Phone string
}

// File upload example
type UploadRequest struct {
	// @evl:validate required
	// @oa:description User avatar image
	Avatar *os.File

	// @evl:validate optional
	// @oa:description Raw data stream
	Payload io.Reader
}
```

### 2. Generate validation code

```bash
go generate ./...
```

This creates two files:

1. `user.gen.go`: Contains the `Validate()` method.
2. `user.oa.gen.go`: Contains the `OaSchema()` method for OpenAPI.

### 3. Use the generated Validate method

```go
package main

func main() {
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	if err := user.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}

	// OpenAPI Schema
	schema := user.OaSchema()
	fmt.Printf("Schema Ref: %s\n", schema.Ref)
}
```

## Commands

### generate — Generate validation code

```bash
# Basic generation
elval-gen generate -i ./user

# With OpenAPI schemas
elval-gen generate -i ./user -openapi

# Verbose output
elval-gen generate -i ./user -v

# Short form
elval-gen gen -i ./user -openapi -v
```

### lint — Validate annotations

```bash
# Check annotations without generating code
elval-gen lint -i ./user

# Show warnings
elval-gen lint -i ./user -v

# Treat warnings as errors
elval-gen lint -i ./user -Werror -v

# Recursive check with excludes
elval-gen lint -i . -r -exclude vendor,testdata -v
```

### Go Generate Integration

```go
//go:generate elval-gen generate -i .

type User struct {
// @evl:validate required
// @evl:validate min:2
Name string
}
```

## Diagnostics & Linting

ElVal provides detailed diagnostics with precise locations, hints, and colored output.

### Output Format

```
[error] user.go:24:2: (validator) @evl:pattern User.Email: directive 'pattern' is not supported for type int
 [Hint]: available: required, optional, min, max, enum, url, http_url...

[warning] user.go:31:2: (validator) @evl:min User.Age: directive 'min' requires one parameter
 [Hint]: example: @evl:validate min:value
```

| Component               | Description                          |
|-------------------------|--------------------------------------|
| `[error]` / `[warning]` | Severity level (colored in terminal) |
| `user.go:24:2`          | File:line:column location            |
| `(validator)`           | Component that issued the diagnostic |
| `@evl:pattern`          | Directive type                       |
| `User.Email`            | Context: struct.field                |
| `[Hint]:`               | Suggestion for fixing the issue      |

### Output Control Flags

| Flag              | Description                                 |
|-------------------|---------------------------------------------|
| `-v`, `--verbose` | Show all warnings and debug information     |
| `--no-color`      | Disable colored output (for logs/pipes)     |
| `-Werror`         | Treat warnings as errors (exit with code 1) |

### Example: `go generate` with diagnostics

```go
//go:generate elval-gen generate -i . -openapi -v

type Product struct {
// @evl:validate required
// @evl:validate pattern:email  ← error: pattern not for int
ID int
}
```

**Output**:

```
[error] product.go:10:2: (validator) @evl:pattern Product.ID: directive 'pattern' is not supported for type int
 [Hint]: available: required, optional, min, max, enum, gt, gte, lt, lte, not-zero

❌ 1 error(s)
Files: generated 0, skipped 1
exit status 1
```

> 💡 **Rule**: Directives with errors or warnings **are excluded** from generated code. This guarantees validation won't
> run with incorrect parameters.

## Annotations

### Required & Optional

| Directive  | Description                            | Example                  |
|------------|----------------------------------------|--------------------------|
| `required` | Field is required                      | `@evl:validate required` |
| `optional` | Field is optional (skips empty values) | `@evl:validate optional` |

### String Validators

| Directive              | Description                   | Example                                |
|------------------------|-------------------------------|----------------------------------------|
| `min:{n}`              | Minimum string length         | `@evl:validate min:3`                  |
| `max:{n}`              | Maximum string length         | `@evl:validate max:50`                 |
| `len:{n}`              | Exact string length           | `@evl:validate len:10`                 |
| `pattern:email`        | Email address validation      | `@evl:validate pattern:email`          |
| `pattern:phone`        | Phone number validation       | `@evl:validate pattern:phone`          |
| `pattern:uuid`         | UUID validation               | `@evl:validate pattern:uuid`           |
| `pattern:{regex}`      | Regular expression validation | `@evl:validate pattern:^[A-Z]{3}-\d+$` |
| `contains:{substr}`    | String must contain substring | `@evl:validate contains:admin`         |
| `starts_with:{prefix}` | String must start with prefix | `@evl:validate starts_with:https://`   |
| `ends_with:{suffix}`   | String must end with suffix   | `@evl:validate ends_with:.go`          |

### Numeric Validators

| Directive  | Description              | Example                  |
|------------|--------------------------|--------------------------|
| `min:{n}`  | Minimum value            | `@evl:validate min:18`   |
| `max:{n}`  | Maximum value            | `@evl:validate max:99`   |
| `gt:{n}`   | Greater than             | `@evl:validate gt:0`     |
| `gte:{n}`  | Greater than or equal to | `@evl:validate gte:18`   |
| `lt:{n}`   | Less than                | `@evl:validate lt:100`   |
| `lte:{n}`  | Less than or equal to    | `@evl:validate lte:100`  |
| `not-zero` | Non-zero value           | `@evl:validate not-zero` |

### Comparison Validators

| Directive     | Description                  | Example                       |
|---------------|------------------------------|-------------------------------|
| `eq:{value}`  | Equal to specified value     | `@evl:validate eq:"active"`   |
| `neq:{value}` | Not equal to specified value | `@evl:validate neq:"deleted"` |

### Enum Validators

| Directive            | Description             | Example                                      |
|----------------------|-------------------------|----------------------------------------------|
| `enum:{v1},{v2},...` | Value from allowed list | `@evl:validate enum:active,inactive,pending` |

### Date Validators (time.Time)

| Directive       | Description           | Example                           |
|-----------------|-----------------------|-----------------------------------|
| `after:{date}`  | Date after specified  | `@evl:validate after:2020-01-01`  |
| `before:{date}` | Date before specified | `@evl:validate before:2025-12-31` |
| `not-zero`      | Non-zero date         | `@evl:validate not-zero`          |

### Duration Validators (time.Duration)

| Directive        | Description       | Example                  |
|------------------|-------------------|--------------------------|
| `min:{duration}` | Minimum duration  | `@evl:validate min:1s`   |
| `max:{duration}` | Maximum duration  | `@evl:validate max:24h`  |
| `not-zero`       | Non-zero duration | `@evl:validate not-zero` |

### Slice Validators

| Directive  | Description           | Example                  |
|------------|-----------------------|--------------------------|
| `required` | Slice cannot be nil   | `@evl:validate required` |
| `not-zero` | Slice cannot be empty | `@evl:validate not-zero` |
| `min:{n}`  | Minimum slice size    | `@evl:validate min:1`    |
| `max:{n}`  | Maximum slice size    | `@evl:validate max:10`   |
| `len:{n}`  | Exact slice size      | `@evl:validate len:3`    |

### URL Validators

| Directive  | Description                | Example                  |
|------------|----------------------------|--------------------------|
| `url`      | Any valid URL              | `@evl:validate url`      |
| `http_url` | HTTP or HTTPS URL          | `@evl:validate http_url` |
| `dsn`      | Database connection string | `@evl:validate dsn`      |

## File & Stream Uploads

ElVal automatically maps file and stream types to standard OpenAPI schemas (compliant
with [OpenAPI 3.2.0 §4.14.7](https://spec.openapis.org/oas/v3.2.0.html#considerations-for-file-uploads)).

### Auto-detection of Standard Types

| Go Type          | OpenAPI Schema                 | Swagger UI             |
|------------------|--------------------------------|------------------------|
| `*os.File`       | `type: string, format: binary` | 📄 Choose file         |
| `io.Reader`      | `type: string, format: binary` | 📄 Choose file         |
| `io.ReadCloser`  | `type: string, format: binary` | 📄 Choose file         |
| `[]byte`         | `type: string, format: byte`   | 📝 Text input (base64) |
| `multipart.File` | `type: string, format: binary` | 📄 Choose file         |

**Example**:

```go
type UploadRequest struct {
// @oa:description User avatar image
Avatar *os.File // → format: binary

// @oa:description Raw payload stream
Payload io.Reader  // → format: binary

// @oa:description Base64-encoded thumbnail
Thumbnail []byte // → format: byte
}
```

### Custom Types: Explicit Annotations

If your type implements `io.Reader` but isn't a standard type, use semantic annotations:

| Annotation        | Effect                         | Example                  |
|-------------------|--------------------------------|--------------------------|
| `@oa:file`        | `type: string, format: binary` | For files in `multipart` |
| `@oa:stream`      | `type: string, format: binary` | For raw streams          |
| `@oa:format:byte` | `type: string, format: byte`   | For base64 data          |

**Example**:

```go
type CustomReader struct { /* ... */ }
func (c *CustomReader) Read(p []byte) (int, error) { /* ... */ }

type UploadRequest struct {
// @oa:file
// @oa:description "Custom file reader"
Avatar *CustomReader // → format: binary
}
```

> 💡 **Note**: Field-level annotation (`@oa:file`) takes priority over type auto-detection. This lets you explicitly
> declare intent even if the type isn't recognized automatically.

### Ignoring Structs

To completely exclude a struct from OpenAPI schema generation, use `@oa:ignore` at the type level:

```go
// @oa:ignore  ← entire struct is ignored
type InternalMarker struct {}

type Response struct {
// This field will be skipped in schema because type is ignored
Marker InternalMarker
}
```

> ⚠️ If a field with an ignored type has an explicit annotation (`@oa:file`, `@oa:stream`), the field annotation takes
> priority — the field will be generated.

## Nested Structures

```go
type Address struct {
// @evl:validate required
// @evl:validate min:2
City string
}

type User struct {
Name    string
Address Address // automatically validated
}
```

## Custom Validators

### Register validator

```go
package main

func init() {
	validator.RegisterCustom("x-color", func(value any, params string) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string")
		}

		validColors := map[string]bool{"red": true, "green": true, "blue": true}
		if !validColors[str] {
			return fmt.Errorf("invalid color: %s", str)
		}
		return nil
	})
}
```

### Use in struct

```go
type Product struct {
// @evl:validate x-color
Color string

// @evl:validate x-even
Count int
}
```

## Decorators

Auto-populate fields from various sources:

| Decorator              | Description                         | Example                              |
|------------------------|-------------------------------------|--------------------------------------|
| `ctx-get:{key}`        | Get value from context              | `@evl:decor ctx-get:user_id`         |
| `httpctx-get:{header}` | Get value from HTTP header          | `@evl:decor httpctx-get:X-User-Role` |
| `env-get:{var}`        | Get value from environment variable | `@evl:decor env-get:APP_ENV`         |
| `time-now`             | Set current time                    | `@evl:decor time-now`                |
| `uuid-gen`             | Generate UUID                       | `@evl:decor uuid-gen`                |

```go
type Request struct {
// @evl:decor ctx-get:user_id
UserID string

// @evl:decor httpctx-get:X-Request-ID
RequestID string

// @evl:decor env-get:APP_ENV
Environment string

// @evl:decor time-now
CreatedAt time.Time
}

// Apply decorators
if err := req.Decorate(ctx); err != nil {
// handle error
}
```

## OpenAPI

Generate OpenAPI 3 schemas from your validation annotations:

```bash
elval-gen -input . -openapi
```

### OpenAPI Annotations

| Annotation                  | Description                      | Example                       |
|-----------------------------|----------------------------------|-------------------------------|
| `@oa:title`                 | Schema title                     | `@oa:title "User Name"`       |
| `@oa:description`           | Schema description               | `@oa:description "Full name"` |
| `@oa:example`               | Example value                    | `@oa:example "John"`          |
| `@oa:format`                | OpenAPI format                   | `@oa:format email`            |
| `@oa:file`                  | Mark field as file for multipart | `@oa:file`                    |
| `@oa:stream`                | Mark field as raw stream         | `@oa:stream`                  |
| `@oa:format:{binary\|byte}` | Explicitly set binary format     | `@oa:format:byte`             |
| `@oa:ignore`                | Exclude struct/field from schema | `// @oa:ignore`               |

#### External Types & Stubs

When using third-party types (e.g., `geojson.Feature`), you can define local "stub" structs to document them properly.
This allows you to control how external types appear in your API documentation.

```go
package docs

// FeatureDocs documents geojson.Feature
// @oa:description "GeoJSON Feature with geometry and properties"
type FeatureDocs struct {
	Geometry   any                    `json:"geometry"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}
```

#### Reference the Stub in your API Struct:

```go
package api

import "github.com/paulmach/orb/geojson"

type CreateLocationRequest struct {
	// Point to your local stub using its GlobalRef
	// @oa:rewrite.ref "github.com/myorg/api/docs.FeatureDocs"
	Feature geojson.Feature `json:"feature"`

	UserID string `json:"user_id"`
}
```

#### Polymorphism (Discriminator & OneOf):

Support for discriminated unions using `@oa:discriminator` and `@oa:oneOf`. This is useful for modeling heterogeneous
collections like GeoJSON geometries.

```go
package docs

// PointDocs represents a Point geometry
// @oa:description "A single geographic coordinate"
type PointDocs struct {
	// @oa:enum "Point"
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

// PolygonDocs represents a Polygon geometry
// @oa:description "A closed geometric shape"
type PolygonDocs struct {
	// @oa:enum "Polygon"
	Type        string       `json:"type"`
	Coordinates [][3]float64 `json:"coordinates"`
}

// FeatureDocs documents a GeoJSON Feature
// @oa:description "GeoJSON Feature with geometry and properties"
// @oa:discriminator.propertyName "geometry.type"
// @oa:discriminator.mapping "Point:PointFeature"
// @oa:discriminator.mapping "Polygon:PolygonFeature"
type FeatureDocs struct {
	// @oa:title "Geometry"
	// @oa:description "The geometric shape"
	// @oa:oneOf "PointDocs,PolygonDocs"
	Geometry any `json:"geometry"`

	Properties map[string]interface{} `json:"properties,omitempty"`
}
```

#### Type Rewriting

If you want to simplify a complex type in the documentation (e.g., show a JSON blob as a simple string), use
`@oa:rewrite.type`.

Supported rewrite types: `string`, `integer`, `number`, `boolean`, `array`.

```go
type UploadRequest struct {
// Technically json.RawMessage, but documented as string
// @oa:rewrite.type string
// @oa:description "JSON payload as a base64 encoded string"
Payload json.RawMessage `json:"payload"`
}
```

#### File & Stream Fields

ElVal automatically generates correct schemas for file fields:

```go
type UploadRequest struct {
// @oa:description User avatar
Avatar *os.File // → {"type": "string", "format": "binary"}
}
```

**Generated schema**:

```yaml
UploadRequest:
  type: object
  properties:
    avatar:
      type: string
      format: binary
      description: User avatar
```

> 🔄 **How it works with nooa/router**:  
> When a spec consumer (`nooa`, Swagger UI, Postman) sees `format: binary`, it automatically:
> 1. Wraps `requestBody` in `multipart/form-data`
> 2. Maps property name (`avatar`) → `Content-Disposition: name="avatar"`
> 3. Handles size/MIME validation (if configured at runtime)
>
> `elval` doesn't know about the runtime — it only generates the standard OpenAPI schema.

#### Nested Structures & Generics

ElVal automatically validates nested structures. It also supports generic wrappers like `Option[T]` or `Result[T]`.

```go
type Product struct {
// Validates the inner Review struct if present
Reviews []model.Option[Review]
}
```

#### Structs Without Body Fields

If a struct has no fields for the request body (only `@oa:in` parameters or `@oa:ignore`), a minimal schema is
generated:

```go
// @oa:ignore
type EmptyMarker struct {}

type Response struct {
// @oa:in path id
ID string
}
```

```yaml
# EmptyMarker is not generated (ignored)
# Response has only parameters, no body:
parameters:
  - name: id
    in: path
    required: true
    schema: { type: string }
```

## Performance

```bash
goos: linux
goarch: amd64
pkg: github.com/arkannsk/elval/test/benchmark
cpu: AMD Ryzen 7 3700X 8-Core Processor             
BenchmarkElvalManual-16                   183387              8545 ns/op            6188 B/op         82 allocs/op
BenchmarkElvalGenerated-16               2627594               440.4 ns/op             0 B/op          0 allocs/op
BenchmarkPlayground-16                    454696              2562 ns/op             389 B/op         14 allocs/op
BenchmarkPlaygroundWithCache-16           483514              2727 ns/op             390 B/op         14 allocs/op
```
