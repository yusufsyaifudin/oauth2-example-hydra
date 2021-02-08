package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"ysf/oauth2-example-hydra/cmd/authc/handler"
	"ysf/oauth2-example-hydra/cmd/authc/repouser"
	"ysf/oauth2-example-hydra/pkg/template"
	"ysf/oauth2-example-hydra/pkg/tracer"
	"ysf/oauth2-example-hydra/views"

	"github.com/opentracing/opentracing-go"

	"github.com/labstack/echo-contrib/jaegertracing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	hydra "github.com/ory/hydra-client-go/client"
)

var (
	adminURL, _ = url.Parse("http://localhost:4445")
	hydraClient = hydra.NewHTTPClientWithConfig(nil,
		&hydra.TransportConfig{
			Schemes:  []string{adminURL.Scheme},
			Host:     adminURL.Host,
			BasePath: adminURL.Path,
		},
	)
)

var userInfo = []repouser.UserInfo{
	{
		ID:       1,
		Email:    "user@example.com",
		Password: "password",
	},
	{
		ID:       2,
		Email:    "user2@example.com",
		Password: "password",
	},
}

func main() {

	// Prepare Opentracing
	var (
		tracerServiceName     = "AuthC"
		tracerURL             = "localhost:6831"
		tracerService, closer = tracer.New(true, tracerServiceName, tracerURL, 1)
	)

	defer func() {
		if closer == nil {
			_, _ = fmt.Fprintf(os.Stderr, "tracer closer is nil\n")
			return
		}

		if err := closer.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "closing tracer error: %s\n", err.Error())
			return
		}
	}()

	// set global tracer of this application
	opentracing.SetGlobalTracer(tracerService)

	controller := handler.Handler{
		HydraAdmin: hydraClient.Admin,
		UserRepo:   repouser.NewMemory(userInfo),
	}

	t := &Template{}

	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Logger())
	e.Use(jaegertracing.TraceWithConfig(jaegertracing.TraceConfig{
		Tracer: tracerService,
	}))

	e.GET("/authentication/login", controller.GetLogin)
	e.POST("/authentication/login", controller.PostLogin)
	e.GET("/authentication/consent", controller.GetConsent)
	e.POST("/authentication/consent", controller.PostConsent)

	if err := e.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}

type Template struct{}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	tmpl := template.New("", &template.BinData{
		Asset:      views.Asset,
		AssetDir:   views.AssetDir,
		AssetNames: views.AssetNames,
	})

	tpl, err := tmpl.Parse(fmt.Sprintf("views/authc/%s", name))
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}
