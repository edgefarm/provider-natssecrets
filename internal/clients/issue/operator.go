package issue

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/operator/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"
)

func OperatorPath(mount string, operator string) string {
	return mount + "/issue/operator/" + operator
}

func ReadOperator(c *vault.Client, operator string) (*v1alpha1.OperatorParameters, *natsbackend.IssueOperatorStatus, error) {
	path := OperatorPath(c.Mount, operator)

	resp, err := vault.Read[natsbackend.IssueOperatorData](c, path)
	if err != nil {
		return nil, nil, err
	}

	return &v1alpha1.OperatorParameters{
		CreateSystemAccount: resp.CreateSystemAccount,
		SyncAccountServer:   resp.SyncAccountServer,
		Claims:              resp.Claims,
	}, &resp.Status, nil
}

func WriteOperator(c *vault.Client, operator string, params *v1alpha1.OperatorParameters) error {
	path := OperatorPath(c.Mount, operator)
	request := &natsbackend.IssueOperatorParameters{
		Operator:            operator,
		CreateSystemAccount: params.CreateSystemAccount,
		SyncAccountServer:   params.SyncAccountServer,
		Claims:              params.Claims,
	}
	return vault.Write(c, path, request)
}

func DeleteOperator(c *vault.Client, operator string) error {
	path := OperatorPath(c.Mount, operator)
	return vault.Delete(c, path)
}
