package issue

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/issue/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func userPath(mount string, operator string, account string, user string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account + "/user/" + user
}

func ReadUser(c *vault.Client, operator string, account string, user string) (*v1alpha1.UserParameters, error) {
	path := userPath(c.Mount, operator, account, user)
	return vault.Read[v1alpha1.UserParameters](c, path)
}

func WriteUser(c *vault.Client, operator string, account string, user string, params *v1alpha1.UserParameters) error {
	path := userPath(c.Mount, operator, account, user)
	return vault.Write(c, path, params)
}

func DeleteUser(c *vault.Client, operator string, account string, user string) error {
	path := userPath(c.Mount, operator, account, user)
	return vault.Delete(c, path)
}
