package multiline_description

//go:generate elval-gen gen -input . -openapi

type User struct {
	// @oa:title "User Name"
	// @oa:description "The full name of the user."
	// @oa:description "It should include first and last name."
	// @oa:description "Max length is 50 characters."
	Name string `json:"name"`
}

// @oa:ignore
type Ignore struct {
	Name string `json:"name"`
}

type IgnoreField struct {
	Name string `json:"name"`
	// @oa:ignore
	Age int `json:"age"`
}
