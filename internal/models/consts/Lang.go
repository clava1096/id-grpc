package consts

import "time"

const TokenAccessLifeTime = 12 * time.Hour
const TokenVerificationLifeTime = time.Minute * 5
const TokenResetPasswordLifeTime = 2 * time.Hour
const TokenEmailVerificationLifeTime = 2 * time.Hour
const TokenTemporaryAccessLifeTime = time.Hour

type TokenType int

const (
	TokenConfirmEmail    TokenType = iota + 1 // генерация ссылки
	TokenAccess                               // доступ к ресурсам
	TokenVerification                         // при отправке сообщения на почту(жизнь 5 минут)
	TokenResetPassword                        // для смены пароля
	TokenTemporaryAccess                      // нужен во время отправки 6-значного кода на почту, подтверждает действие пользователя
)

type UsersRole string

const (
	User      UsersRole = "user"
	Moderator UsersRole = "moderator"
	Admin     UsersRole = "admin"
)

type UUIDResponse string
