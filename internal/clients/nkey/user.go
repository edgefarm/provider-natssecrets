package nkey

import (
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func UserPath(mount string, operator string, account string, user string) string {
	return mount + "/nkey/operator/" + operator + "/account/" + account + "/user/" + user
}

func ReadUser(c *vault.Client, operator string, account string, user string) (*NkeyData, error) {
	path := UserPath(c.Mount, operator, account, user)
	return vault.Read[NkeyData](c, path)
}
