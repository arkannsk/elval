package string_checks

//go:generate elval-gen -input .

type Document struct {
	// @evl:validate required
	// @evl:validate contains:README
	Name string

	// @evl:validate required
	// @evl:validate starts_with:https://
	// @evl:validate ends_with:.com
	URL string

	// @evl:validate required
	// @evl:validate contains:world
	Content string

	// @evl:validate optional
	// @evl:validate starts_with:img_
	ImageName string
}

type File struct {
	// @evl:validate required
	// @evl:validate ends_with:.go
	Path string

	// @evl:validate required
	// @evl:validate contains:test
	Name string
}
