package fixtures

const (
	TitlePrefix       = "autotest-"
	TitleRandomLength = 10
	Description       = "description for autotest"
	Price             = 1000.0
	Quantity          = 3
	NegativePrice     = -1.0
	MaxTitleLength    = 255
	MaxPhotosCount    = 3
)

const (
	MissingAdvertisementID = "11111111-1111-1111-1111-111111111111"
	InvalidToken           = "invalid-token-xxx"
)

const ErrorCodeUnauthorized = "unauthorized"

const (
	UserEmailPrefix    = "autotest-user-"
	UserEmailDomain    = "@test123.ru"
	UserFirstName      = "Autotest"
	UserPasswordPrefix = "Aa1-"
	UserRandomLength   = 12
)

func JPEGPhoto() []byte {
	return []byte{0xFF, 0xD8, 0xFF, 0xD9}
}
