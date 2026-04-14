package main

//go:generate elval-gen -input .

import "time"

type User struct {
	// @evl:decor ctx-get:user_id
	// @evl:validate required
	ID string

	// @evl:decor httpctx-get:X-User-Role
	// @evl:validate required
	// @evl:validate enum:admin,user,guest
	Role string

	// @evl:decor env-get:APP_ENV
	Environment string

	// @evl:decor time-now
	CreatedAt time.Time

	// @evl:decor uuid-gen
	RequestID string
}
