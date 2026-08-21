package internal

/*
import (
	"encoding/json"
	"errors"
	"game/duck/logger"
	"io"

	"github.com/labstack/echo/v4"
)

func initHttp() {

		e := comm.NewEcho()

		e.Use(LogRequest)

		e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "hello") })
		e.GET("/config/token.json", GetToken)

		g := e.Group("/api", AppIDCheck, IpCheck)
		{
			g.POST("/GameList", GameList)
			g.POST("/Login", Login)
			g.POST("/GetLog", GetLog, RateLimit)
		}

		port, err := lazy.PortProvider.GetPort(lazy.ServiceName + ".http")
		if err != nil {
			panic(err)
		}

		err = e.Start(fmt.Sprintf(":%v", port))
		if err != nil {
			panic(err)
		}
}

func GetApp(c echo.Context) *MemApp {
	return c.Get("app").(*MemApp)
}

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info(c.Request().Method, c.Request().RequestURI)
		return next(c)
	}
}

func AppIDCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		h := c.Request().Header

		AppID := h.Get("AppID")
		AppSecret := h.Get("AppSecret")

		logger.Info(AppID, AppSecret)

		app := appMgr.GetApp(AppID)
		if app == nil {
			return errors.New("invalid AppID")
		}
		if app.AppSecret != AppSecret {
			return errors.New("invalid AppSecret")
		}

		c.Set("app", app)

		return next(c)
	}
}

func IpCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: impl
		return next(c)
	}
}

func RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: impl
		return next(c)
	}
}

func BindJson(c echo.Context, j any) error {
	defer c.Request().Body.Close()
	buf, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, j)
}

*/
