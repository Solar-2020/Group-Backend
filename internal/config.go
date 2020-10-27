package internal

var (
	Config configTemplate
)

type configTemplate struct {
	Port                          string `envconfig:"PORT" default:"8099"`
	GroupDataBaseConnectionString string `envconfig:"GROUP_DB_CONNECTION_STRING" default:"-"`
	InviteLinkPrefix			  string `envconfig:"INVITE_GROUP_PREFIX_ADDRESS" default:"http://nl-mail.ru/welcome"`
	AuthServiceAddress			  string  `envconfig:"AUTH_SERVICE_ADDRESS" default:""`
}
