package queries

// SignUp struct to describe register a new user.
type SignUp struct {
	Username string `json:"username" validate:"required,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}

// SignIn struct to describe login user.
type SignIn struct {
	Username string `json:"username" validate:"required,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}
