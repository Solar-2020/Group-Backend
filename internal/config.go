package internal

import "github.com/Solar-2020/GoUtils/common"

var (
	Config configTemplate
)

type configTemplate struct {
	common.SharedConfig
	GroupDataBaseConnectionString string `envconfig:"GROUP_DB_CONNECTION_STRING" required:"true"`
	InviteLinkPrefix              string `envconfig:"INVITE_GROUP_PREFIX_ADDRESS" default:"http://nl-mail.ru/welcome"`
	ServerSecret                  string `envconfig:"SERVER_SECRET" default:"Basic secret"`
	AccountServiceHost            string `envconfig:"ACCOUNT_SERVICE_HOST" default:"develop.pay-together.ru"`
	SendInviteLetter				bool   `envconfig:"SEND_INVITE_LETTERS" default:"false"`

	InviteLetterHost				string `envconfig:"INVITE_LETTERS_HOST" default:"smtp.mail.ru"`
	InviteLetterSender				string `envconfig:"INVITE_LETTERS_SENDER" default:"invite@pay-together.ru"`
	InviteLetterSenderPassword		string `envconfig:"INVITE_LETTERS_PASSWORD" required:"true"`
	InviteLetterBasePath			string `envconfig:"INVITE_LETTERS_BASE_PATH" default:"/templates"`
	InviteLetterTimespan			int    `envconfig:"INVITE_LETTERS_TIMESPAN" default:"20"`
}
