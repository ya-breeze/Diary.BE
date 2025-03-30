package webapp

import (
	"fmt"
	"html/template"
	"log/slog"
	"math"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/utils"
)

type WebAppRouter struct {
	commit  string
	logger  *slog.Logger
	cfg     *config.Config
	db      database.Storage
	cookies *sessions.CookieStore
}

func NewWebAppRouter(
	commit string, logger *slog.Logger, cfg *config.Config, db database.Storage,
) *WebAppRouter {
	return &WebAppRouter{
		commit:  commit,
		logger:  logger,
		cfg:     cfg,
		db:      db,
		cookies: sessions.NewCookieStore([]byte("SESSION_KEY")),
	}
}

func (r *WebAppRouter) Routes() goserver.Routes {
	return goserver.Routes{
		"RootPath": goserver.Route{
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: r.homeHandler,
		},
		"Login": goserver.Route{
			Method:      "POST",
			Pattern:     "/",
			HandlerFunc: r.loginHandler,
		},
		"AboutPath": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/about",
			HandlerFunc: r.aboutHandler,
		},
	}
}

func (r *WebAppRouter) loadTemplates() (*template.Template, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"formatTime": utils.FormatTime,
		"decrease": func(i int) int {
			return i - 1
		},
		"money": func(num float64) float64 {
			return math.Round(num*100) / 100
		},
		"timestamp": func(t time.Time) int64 {
			return t.Unix()
		},
		"lastMonth": func(t time.Time) time.Time {
			return time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, t.Location())
		},
		"addMonths": func(t time.Time, num int) time.Time {
			return time.Date(t.Year(), t.Month()+time.Month(num), 1, 0, 0, 0, 0, t.Location())
		},
		"addQueryParam": func(rawURL string, key string, value any) (string, error) {
			u, err := url.Parse(rawURL)
			if err != nil {
				return "", err
			}
			q := u.Query()
			q.Set(key, fmt.Sprintf("%v", value))
			u.RawQuery = q.Encode()
			return u.String(), nil
		},
	}).ParseGlob(filepath.Join("webapp", "templates", "*.tpl"))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
