package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func (h Handler) PostLogin(c echo.Context) error {
	ctx := c.Request().Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "PostLogin")
	defer func() {
		span.Finish()
		ctx.Done()
	}()

	formData := struct {
		LoginChallenge string `validate:"required"`
		Email          string `validate:"required"`
		Password       string `validate:"required"`
		RememberMe     string `validate:"required"`
	}{
		LoginChallenge: c.FormValue("login_challenge"),
		Email:          c.FormValue("email"),
		Password:       c.FormValue("password"),
		RememberMe:     c.FormValue("remember_me"),
	}

	// TODO validation

	var rememberMe = formData.RememberMe == "true"

	user, err := h.UserRepo.GetUserByEmail(c.Request().Context(), formData.Email)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}

	if user.Password != formData.Password {
		return c.String(http.StatusNotFound, "Wrong username and password")
	}

	// Using Hydra Admin to accept login request!
	loginGetParam := admin.NewGetLoginRequestParams()
	loginGetParam.SetLoginChallenge(formData.LoginChallenge)

	_, err = h.HydraAdmin.GetLoginRequest(loginGetParam)
	if err != nil {
		// if error, redirects to ...
		str := fmt.Sprint("error GetLoginRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	subject := fmt.Sprint(user.ID)

	loginAcceptParam := admin.NewAcceptLoginRequestParams()
	loginAcceptParam.WithContext(ctx)
	loginAcceptParam.SetLoginChallenge(formData.LoginChallenge)
	loginAcceptParam.SetBody(&models.AcceptLoginRequest{
		Subject:  &subject,
		Remember: rememberMe,
	})

	respLoginAccept, err := h.HydraAdmin.AcceptLoginRequest(loginAcceptParam)
	if err != nil {
		// if error, redirects to ...
		str := fmt.Sprint("error AcceptLoginRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	// If success, it will redirect to consent page using handler GetConsent
	// It then show the consent form
	return c.Redirect(http.StatusFound, *respLoginAccept.GetPayload().RedirectTo)
}
