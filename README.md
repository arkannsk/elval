# ElVal — Lightning Fast Go Validator with Code Generation & OpenAPI Support

[![Go Version](https://img.shields.io/github/go-mod/go-version/arkannsk/elval)](https://golang.org/)
[![License](https://img.shields.io/github/license/arkannsk/elval)](LICENSE)
[![Benchmarks](https://img.shields.io/badge/benchmarks-440_allocs%2Fop-brightgreen)](BENCHMARKS.md)

[//]: # ([![Go Reference]&#40;https://pkg.go.dev/badge/github.com/arkannsk/elval.svg&#41;]&#40;https://pkg.go.dev/github.com/arkannsk/elval&#41;)
[//]: # ([![Go Report Card]&#40;https://goreportcard.com/badge/github.com/arkannsk/elval&#41;]&#40;https://goreportcard.com/report/github.com/arkannsk/elval&#41;)


**ElVal** is a code generation-based validator for Go that eliminates reflection overhead entirely. By generating type-safe validation code at build time, ElVal achieves **6x faster validation** with **zero memory allocations** at runtime.

Additionally, ElVal generates **OpenAPI 3.0 schemas** directly from your struct annotations, supporting complex types, external libraries (via stubs), polymorphism, and custom type rewriting.

## Table of Contents

* [ElVal — Lightning Fast Go Validator with Code Generation &amp; OpenAPI Support](#elval--lightning-fast-go-validator-with-code-generation--openapi-support)
    * [Table of Contents](#table-of-contents)
    * [Features](#features)
    * [Quick Start](#quick-start)
        * [1\. Define your struct with annotations](#1-define-your-struct-with-annotations)
        * [2\. Generate validation code](#2-generate-validation-code)
        * [3\. Use the generated Validate method](#3-use-the-generated-validate-method)
    * [Commands](#commands)
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
            * [Nested Structures &amp; Generics](#nested-structures--generics)
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
```

### 2. Generate validation code

```bash
go generate ./...
```

This creates two files:
1. user.gen.go: Contains the Validate() method.
2. user.oa.gen.go: Contains the OaSchema() method for OpenAPI.

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

generate - Generate validation code

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

lint - Validate annotations

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

Go Generate Integration

```go
//go:generate elval-gen generate -i .

type User struct {
    // @evl:validate required
    // @evl:validate min:2
    Name string
}
```

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

| Annotation      | Description        | Example                     |
|-----------------|--------------------|-----------------------------|
| @oa:title       | Schema title       | @oa:title "User Name"       |
| @oa:description | Schema description | @oa:description "Full name" |
| @oa:example     | Example value      | @oa:example "John"          |
| @oa:format      | OpenAPI format     | @oa:format email            |

#### External Types & Stubs

When using third-party types (e.g., geojson.Feature), you can define local "stub" structs to document them properly. 
This allows you to control how external types appear in your API documentation.

```go
package docs

// FeatureDocs documents geojson.Feature
// @oa:description "GeoJSON Feature with geometry and properties"
type FeatureDocs struct {
    Geometry any `json:"geometry"`
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

Support for discriminated unions using @oa:discriminator and @oa:oneOf. This is useful for modeling heterogeneous 
collections like GeoJSON geometries.

```go
package docs

// PointDocs represents a Point geometry
// @oa:description "A single geographic coordinate"
type PointDocs struct {
    // @oa:enum "Point"
    Type string `json:"type"`
    Coordinates [2]float64 `json:"coordinates"`
}

// PolygonDocs represents a Polygon geometry
// @oa:description "A closed geometric shape"
type PolygonDocs struct {
    // @oa:enum "Polygon"
    Type string `json:"type"`
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

If you want to simplify a complex type in the documentation (e.g., show a JSON blob as a simple string), 
use `@oa:rewrite.type.`

Supported rewrite types: string, integer, number, boolean, array.

```go
type UploadRequest struct {
    // Technically json.RawMessage, but documented as string
    // @oa:rewrite.type string
    // @oa:description "JSON payload as a base64 encoded string"
    Payload json.RawMessage `json:"payload"`
}
```
#### Nested Structures & Generics

ElVal automatically validates nested structures. It also supports generic wrappers like Option[T] or Result[T].

```go
type Product struct {
    // Validates the inner Review struct if present
    Reviews []model.Option[Review] 
}
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
