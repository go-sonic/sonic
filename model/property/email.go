package property

import "reflect"

var (
	EmailHost = Property{
		KeyValue:     "email_host",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	EmailProtocol = Property{
		KeyValue:     "email_protocol",
		DefaultValue: "smtp",
		Kind:         reflect.String,
	}
	EmailSSLPort = Property{
		KeyValue:     "email_ssl_port",
		DefaultValue: 465,
		Kind:         reflect.Int,
	}
	EmailUsername = Property{
		KeyValue:     "email_username",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	EmailPassword = Property{
		KeyValue:     "email_password",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	EmailFromName = Property{
		KeyValue:     "email_from_name",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	EmailIsEnabled = Property{
		KeyValue:     "email_enabled",
		DefaultValue: false,
		Kind:         reflect.Bool,
	}
	EmailStarttls = Property{
		KeyValue:     "email_starttls",
		DefaultValue: false,
		Kind:         reflect.Bool,
	}
)
