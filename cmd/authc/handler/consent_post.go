package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func (h Handler) PostConsent(c echo.Context) error {
	ctx := c.Request().Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "PostConsent")
	defer func() {
		span.Finish()
		ctx.Done()
	}()

	formData := struct {
		ConsentChallenge string   `validate:"required"`
		GrantScope       []string `validate:"required"`
	}{
		ConsentChallenge: c.FormValue("consent_challenge"),
		GrantScope:       c.Request().Form["grant_scope"],
	}

	consentGetParams := admin.NewGetConsentRequestParams()
	consentGetParams.WithContext(ctx)
	consentGetParams.SetConsentChallenge(formData.ConsentChallenge)

	consentGetResp, err := h.HydraAdmin.GetConsentRequest(consentGetParams)
	if err != nil {
		// if error, redirects to ...
		str := fmt.Sprint("error GetConsentRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	// If a user has granted this application the requested scope, hydra will tell us to not show the UI.

	// Now it's time to grant the consent request. You could also deny the request if something went terribly wrong
	consentAcceptBody := &models.AcceptConsentRequest{
		GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
		GrantScope:               formData.GrantScope,
	}

	consentAcceptParams := admin.NewAcceptConsentRequestParams()
	consentAcceptParams.WithContext(ctx)
	consentAcceptParams.SetConsentChallenge(formData.ConsentChallenge)
	consentAcceptParams.WithBody(consentAcceptBody)

	consentAcceptResp, err := h.HydraAdmin.AcceptConsentRequest(consentAcceptParams)
	if err != nil {
		str := fmt.Sprint("error AcceptConsentRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	return c.Redirect(http.StatusFound, *consentAcceptResp.GetPayload().RedirectTo)
}
