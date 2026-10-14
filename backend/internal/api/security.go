package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sarat03/lanlookout/internal/models"
)

// sameSiteOnly - refuse API calls a browser makes on behalf of another site.
// The API has no login and changes data on GET, so without this any web page
// could load <img src="http://lanlookout:8840/api/host/del/1">.
// Scripts, curl and Home Assistant send no Sec-Fetch-Site and are not affected.
// ponytail: browsers without Sec-Fetch-Site (pre-2023 Safari) fall back to the Origin check.
func sameSiteOnly(c *gin.Context) {

	site := c.GetHeader("Sec-Fetch-Site")
	origin := c.GetHeader("Origin")

	crossOrigin := false
	if origin != "" {
		u, err := url.Parse(origin)
		crossOrigin = err != nil || u.Host != c.Request.Host
	}

	if site == "cross-site" || (site == "" && crossOrigin) {
		c.AbortWithStatusJSON(http.StatusForbidden, "Cross-site request blocked")
		return
	}
	c.Next()
}

const hidden = "********"

// masked - config without secrets, for /api/config
func masked(cfg models.Conf) models.Conf {
	cfg.ShoutURL = maskURL(cfg.ShoutURL)
	cfg.PGConnect = maskURL(cfg.PGConnect)
	cfg.InfluxToken = maskAll(cfg.InfluxToken)
	cfg.MqttPass = maskAll(cfg.MqttPass)
	cfg.DNSPass = maskAll(cfg.DNSPass)
	return cfg
}

func maskAll(s string) string {
	if s == "" {
		return ""
	}
	return hidden
}

// maskURL hides the user-info part: the password, or the token Shoutrrr keeps there.
// ponytail: secrets in query params (e.g. ?token=) are not hidden.
func maskURL(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.User == nil {
		return s
	}
	if _, ok := u.User.Password(); ok {
		u.User = url.UserPassword(u.User.Username(), "MASK")
	} else {
		u.User = url.User("MASK")
	}
	return strings.ReplaceAll(u.String(), "MASK", hidden)
}

// keep - a form sent back the masked value unchanged: keep the real one
func keep(sent, current string, mask func(string) string) string {
	if current != "" && sent == mask(current) {
		return current
	}
	return sent
}
