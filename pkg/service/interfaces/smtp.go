package interfaces

type SMTPService interface {
	// SendEmail(e)

	SendUserVerificationEmail(email, firstName, token string) error
}
