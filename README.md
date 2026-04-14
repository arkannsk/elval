# ElVal - Lightning Fast Go Validator with Code Generation

[![Go Version](https://img.shields.io/github/go-mod/go-version/arkannsk/elval)](https://golang.org/)
[![License](https://img.shields.io/github/license/arkannsk/elval)](LICENSE)
[![Benchmarks](https://img.shields.io/badge/benchmarks-463ns%2Fop-brightgreen)](BENCHMARKS.md)

**ElVal** is a code generation-based validator for Go that eliminates reflection overhead.
It's **6x faster** than go-playground/validator with **zero memory allocations**.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Annotations](#annotations)
    - [Required & Optional](#required--optional)
    - [String Validators](#string-validators)
    - [Numeric Validators](#numeric-validators)
    - [Comparison Validators](#comparison-validators)
    - [Enum Validators](#enum-validators)
    - [Date Validators](#date-validators)
    - [Slice Validators](#slice-validators)
    - [URL Validators](#url-validators)
- [Nested Structures](#nested-structures)
- [Custom Validators](#custom-validators)
- [Decorators](#decorators)
- [Performance](#performance)

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

//go:generate elval-gen -input .

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

### 3. Use the generated Validate method

```go
func main() {
user := User{
Name:  "John Doe",
Email: "john@example.com",
Age:   25,
}

if err := user.Validate(); err != nil {
fmt.Printf("Validation failed: %v\n", err)
}
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
func init() {
validator.RegisterCustom("x-color", func (value any, params string) error {
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
