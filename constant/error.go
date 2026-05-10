package constant

import (
	"github.com/gofiber/fiber/v2"
)

type InternalError int

const (
	StatusBadRequestErrorValidation            = 4001
	StatusBadRequestErrorParsingJson           = 4002
	StatusBadRequestBodyMultipartForm          = 4003
	StatusBadRequestExceededFileSize           = 4004
	StatusBadRequestNotExists                  = 4005
	StatusBadRequestExists                     = 4006
	StatusBadRequestExceededLimitSupporter     = 4007
	StatusUnauthorizedMemberIsNotRegistered    = 4011
	StatusUnauthorizedBacalegIsRegistered      = 4012
	StatusUnauthorizedTokenExpired             = 4013
	StatusForbiddenInvalidToken                = 4031
	StatusNotFound                             = 4041
	StatusMethodNotAllowed                     = 4051
	StatusInternalServerErrorDatabaseMysql     = 5001
	StatusInternalServerErrorApps              = 5002
	StatusInternalServerErrorDatabaseRedis     = 5003
	StatusInternalServerErrorServiceMemberApps = 5004
	StatusInternalServerErrorServiceAuthApps   = 5005
)

type errorInfo struct {
	HttpCode    int
	Title       string
	UserMessage UserMessage
}

type UserMessage struct {
	En string `json:"en"`
	Id string `json:"id"`
}

var constantError = map[InternalError]errorInfo{
	StatusUnauthorizedMemberIsNotRegistered: {
		HttpCode: fiber.StatusUnauthorized,
		Title:    "Member is not registered",
		UserMessage: UserMessage{
			En: "Member is not registered, please register first!",
			Id: "Member belum registrasi, silahkan register terlebih dahulu!",
		},
	},
	StatusUnauthorizedBacalegIsRegistered: {
		HttpCode: fiber.StatusUnauthorized,
		Title:    "Bacaleg is registered",
		UserMessage: UserMessage{
			En: "Bacaleg is registered, you cannot register again!",
			Id: "Bacaleg sudah teregistrasi, kamu tidak bisa registrasi lagi!",
		},
	},
	StatusUnauthorizedTokenExpired: {
		HttpCode: fiber.StatusUnauthorized,
		Title:    "Token is expired",
		UserMessage: UserMessage{
			En: "Token is expired, please do login",
			Id: "Token sudah kadaluarsa, silahkan lakukan login",
		},
	},
	StatusForbiddenInvalidToken: {
		HttpCode: fiber.StatusForbidden,
		Title:    "Invalid token",
		UserMessage: UserMessage{
			En: "Invalid token, please try again",
			Id: "Token tidak sah, silahkan dicoba lagi",
		},
	},
	StatusBadRequestErrorValidation: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The request is invalid",
		UserMessage: UserMessage{
			En: "The user request is invalid, please try again",
			Id: "Request user tidak sah, silahkan dicoba lagi",
		},
	},
	StatusBadRequestErrorParsingJson: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The json request is invalid",
		UserMessage: UserMessage{
			En: "The json request is invalid, please try again",
			Id: "Request json tidak sah, silahkan dicoba lagi",
		},
	},
	StatusBadRequestBodyMultipartForm: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The multipart form is invalid",
		UserMessage: UserMessage{
			En: "The multipart form is invalid, please try again",
			Id: "Form multipart tidak sah, silahkan dicoba lagi",
		},
	},
	StatusBadRequestExceededFileSize: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The file size is exceeded limit",
		UserMessage: UserMessage{
			En: "The file size is exceeded limit, please try again",
			Id: "Ukuran berkas melebihi batas maksimal, silahkan dicoba lagi",
		},
	},
	StatusBadRequestNotExists: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The data is not exists",
		UserMessage: UserMessage{
			En: "The data is not exists",
			Id: "Data tidak ada",
		},
	},
	StatusBadRequestExists: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The data is exists",
		UserMessage: UserMessage{
			En: "The data is exists",
			Id: "Data sudah ada",
		},
	},
	StatusBadRequestExceededLimitSupporter: {
		HttpCode: fiber.StatusBadRequest,
		Title:    "The total supporter is exceeded limit",
		UserMessage: UserMessage{
			En: "The total supporter bacaleg is exceeded limit",
			Id: "Jumlah supporter bacaleg melebihi limitasi",
		},
	},
	StatusNotFound: {
		HttpCode: fiber.StatusNotFound,
		Title:    "The endpoint is not found",
		UserMessage: UserMessage{
			En: "The endpoint is not found, please try again",
			Id: "Endpoint tidak ditemukan, silahkan dicoba lagi",
		},
	},
	StatusMethodNotAllowed: {
		HttpCode: fiber.StatusMethodNotAllowed,
		Title:    "The method is not allowed",
		UserMessage: UserMessage{
			En: "The method is not allowed on this endpoint",
			Id: "Method tidak diizinkan untuk endpoint ini",
		},
	},
	StatusInternalServerErrorDatabaseMysql: {
		HttpCode: fiber.StatusInternalServerError,
		Title:    "Internal Server Error ",
		UserMessage: UserMessage{
			En: "There is some problem with server in Database Mysql",
			Id: "Ada masalah dengan server database Mysql",
		},
	},
	StatusInternalServerErrorApps: {
		HttpCode: fiber.StatusInternalServerError,
		Title:    "Internal Server Error Apps",
		UserMessage: UserMessage{
			En: "There is some problem with Apps, Please try again",
			Id: "Ada masalah dengan apps, silahkan dicoba lagi",
		},
	},
	StatusInternalServerErrorDatabaseRedis: {
		HttpCode: fiber.StatusInternalServerError,
		Title:    "Internal Server Error",
		UserMessage: UserMessage{
			En: "There is some problem with server",
			Id: "Ada masalah dengan server",
		},
	},
	StatusInternalServerErrorServiceMemberApps: {
		HttpCode: fiber.StatusInternalServerError,
		Title:    "Internal Server Error Service",
		UserMessage: UserMessage{
			En: "There is some problem with other service, Please try again",
			Id: "Ada masalah dengan service lain, silahkan dicoba lagi",
		},
	},
	StatusInternalServerErrorServiceAuthApps: {
		HttpCode: fiber.StatusInternalServerError,
		Title:    "Internal Server Error Service",
		UserMessage: UserMessage{
			En: "There is some problem with other service, please try again",
			Id: "Ada masalah dengan service lain, silahkan dicoba lagi",
		},
	},
}

func (i InternalError) Info() errorInfo {
	return constantError[i]
}
