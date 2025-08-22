package consts

type ErrorCode string

const (
	SQLError        ErrorCode = "SQLError"
	NotFound        ErrorCode = "NotFound"
	ErrorLoadConfig           = iota
	ErrorConnectSentry
	ErrorConnectOpenSearch
	ErrorConnectTelegram
	ErrorConnectPostgres
	ErrorConnectRedis
	ErrorConnectMinio
	ErrorConnectRabbitMQ
	ErrorConnectCASBIN
)
