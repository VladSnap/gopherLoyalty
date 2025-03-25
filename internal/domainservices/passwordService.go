package domainservices

type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(storedHash, password string) (bool, error)
}
