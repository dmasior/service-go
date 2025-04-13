package testdb

type UserFixture struct {
	Email    string
	Password string
}

var (
	UserJohn = UserFixture{
		Email:    "john@example.com",
		Password: "johnpassword",
	}

	UserJane = UserFixture{
		Email:    "jane@example.com",
		Password: "janepassword",
	}
)
