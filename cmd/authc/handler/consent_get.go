package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func (h Handler) GetConsent(c echo.Context) error {
	ctx := c.Request().Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetConsent")
	defer func() {
		span.Finish()
		ctx.Done()
	}()

	consentChallenge := strings.TrimSpace(c.QueryParam("consent_challenge"))
	if consentChallenge == "" {
		return c.Render(http.StatusOK, "consent.html", map[string]interface{}{
			"ErrorTitle":   "Cannot Accept Consent Request",
			"ErrorContent": "Consent challenge is empty",
		})
	}

	consentGetParams := admin.NewGetConsentRequestParams()
	consentGetParams.WithContext(ctx)
	consentGetParams.SetConsentChallenge(consentChallenge)

	consentGetResp, err := h.HydraAdmin.GetConsentRequest(consentGetParams)
	if err != nil {
		return c.Render(http.StatusOK, "consent.html", map[string]interface{}{
			"ErrorTitle":   "Cannot Accept Consent Request",
			"ErrorContent": err.Error(),
		})
	}

	// If a user has granted this application the requested scope, hydra will tell us to not show the UI.
	if consentGetResp.GetPayload().Skip {
		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		// Now it's time to grant the consent request.
		// You could also deny the request if something went terribly wrong
		consentAcceptBody := &models.AcceptConsentRequest{
			GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
			GrantScope:               consentGetResp.GetPayload().RequestedScope,
		}

		consentAcceptParams := admin.NewAcceptConsentRequestParams()
		consentAcceptParams.WithContext(ctx)
		consentAcceptParams.SetConsentChallenge(consentChallenge)
		consentAcceptParams.WithBody(consentAcceptBody)

		consentAcceptResp, err := h.HydraAdmin.AcceptConsentRequest(consentAcceptParams)
		if err != nil {
			str := fmt.Sprint("error AcceptConsentRequest", err.Error())
			return c.String(http.StatusUnprocessableEntity, str)
		}

		return c.Redirect(http.StatusFound, *consentAcceptResp.GetPayload().RedirectTo)
	}

	consentMessage := fmt.Sprintf("Application %s wants access resources on your behalf and to:",
		consentGetResp.GetPayload().Client.ClientName,
	)

	return c.Render(http.StatusOK, "consent.html", map[string]interface{}{
		"ConsentChallenge": consentChallenge,
		"ConsentMessage":   consentMessage,
		"RequestedScopes":  consentGetResp.GetPayload().RequestedScope,
	})
}
