package url_validation

//go:generate elval-gen -input .

type Link struct {
	// @evl:validate required
	// @evl:validate url
	Website string

	// @evl:validate optional
	// @evl:validate url
	Blog string

	// @evl:validate required
	// @evl:validate url
	// @evl:validate contains:api
	API string
}

type Profile struct {
	// @evl:validate required
	// @evl:validate http_url
	// @evl:validate starts_with:https://
	SecureURL string
}

type Config struct {
	// @evl:validate required
	// @evl:validate url
	AnyURL string

	// @evl:validate required
	// @evl:validate http_url
	WebURL string

	// @evl:validate required
	// @evl:validate dsn
	DatabaseURL string

	// @evl:validate optional
	// @evl:validate dsn
	ClickHouseURL string
}
