package response

type UserFetchEmailPasswordLogin struct {
	UserId         int64
	PasswordHash   string
	EmailConfirmed bool
}
