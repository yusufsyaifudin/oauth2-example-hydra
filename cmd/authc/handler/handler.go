package handler

import (
	"ysf/oauth2-example-hydra/cmd/authc/repouser"

	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
)

type Handler struct {
	HydraAdmin hydraAdmin.ClientService
	UserRepo   repouser.Repository
}
