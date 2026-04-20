package oa_discriminator

//go:generate elval-gen gen -input . -openapi

// Pet Родительская схема
// @oa:discriminator.propertyName "type"
// @oa:discriminator.mapping "cat:Cat"
// @oa:discriminator.mapping "dog:Dog"
type Pet struct {
	// @oa:example "cat"
	Type string
	Name string
}

// Cat Дочерние схемы (позже можно добавить @oa:allOf "Pet" для явного AllOf)
type Cat struct {
	Pet
	// @oa:example true
	Meows bool
}

type Dog struct {
	Pet
	// @oa:example 5
	BarkVolume int
}
