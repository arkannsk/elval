package date_validation

//go:generate elval-gen -input .

type Event struct {
	// @evl:validate required
	// @evl:validate date:RFC3339
	CreatedAt string

	// @evl:validate required
	// @evl:validate date:RFC3339,RFC3339Nano
	UpdatedAt string

	// @evl:validate optional
	// @evl:validate date:2006-01-02
	DateOnly string

	// @evl:validate required
	// @evl:validate date:2006-01-02T15:04:05,RFC3339
	Timestamp string
}

type LogEntry struct {
	// @evl:validate required
	// @evl:validate date:2006-01-02T15:04:05Z07:00,RFC3339
	EventTime string

	// @evl:validate optional
	// @evl:validate date:Kitchen
	KitchenTime string
}
