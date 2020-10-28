package internal

import "github.com/Solar-2020/GoUtils/common"

var (
	Config configTemplate
)

type configTemplate struct {
	common.SharedConfig
	GroupDataBaseConnectionString string `envconfig:"GROUP_DB_CONNECTION_STRING" default:"-"`
	InviteLinkPrefix			  string `envconfig:"INVITE_GROUP_PREFIX_ADDRESS" default:"http://nl-mail.ru/welcome"`
}
