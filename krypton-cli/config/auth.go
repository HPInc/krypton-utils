package config

func (a *Auth) getTokenType() string {
	return a.TokenType
}

func (a *Auth) getManagementServer() string {
	return a.ManagementServer
}
