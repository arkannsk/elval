package response_content

//go:generate elval-gen gen -input . -openapi

// UserResponse ответ, который может быть возвращён в JSON или XML
// @oa:response "200" "application/json,application/xml"
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse ответ с ошибкой
// @oa:response "400" "application/json"
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MultiResponse ответ с несколькими status code
// @oa:response "200" "application/json,application/xml"
// @oa:response "201" "application/json"
type CreateUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
