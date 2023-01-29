package http

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signupResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type ListResponse struct {
	Total  int64       `json:"total"`
	Count  int         `json:"count"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Data   interface{} `json:"data"`
}

type APIErrorResponse struct {
	Error string `json:"error"`
}
