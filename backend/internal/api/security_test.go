package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMaskKeep(t *testing.T) {
	cases := map[string]string{
		"postgres://u:secret@db:5432/x":   "postgres://u:********@db:5432/x",
		"telegram://123:abc@telegram?c=1": "telegram://123:********@telegram?c=1",
		"discord://tok@chan":              "discord://********@chan",
		"generic+http://h/hook":           "generic+http://h/hook",
		"":                                "",
	}
	for in, want := range cases {
		if got := maskURL(in); got != want {
			t.Errorf("maskURL(%q) = %q, want %q", in, got, want)
		}
		if got := keep(maskURL(in), in, maskURL); got != in {
			t.Errorf("keep lost %q, got %q", in, got)
		}
	}
	if keep("new", "old", maskAll) != "new" || keep(hidden, "old", maskAll) != "old" {
		t.Error("keep with maskAll")
	}
}

func TestSameSiteOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(sameSiteOnly)
	r.GET("/x", func(c *gin.Context) { c.Status(200) })

	for _, tc := range []struct {
		site, origin string
		want         int
	}{
		{"", "", 200},                    // curl / Home Assistant
		{"same-origin", "", 200},         // the app itself
		{"cross-site", "", 403},          // <img> from another site
		{"", "http://evil.example", 403}, // old browser, foreign Origin
		{"", "http://wyl.lan:8840", 200}, // old browser, own Origin
	} {
		req := httptest.NewRequest("GET", "http://wyl.lan:8840/x", nil)
		if tc.site != "" {
			req.Header.Set("Sec-Fetch-Site", tc.site)
		}
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Errorf("%+v: got %d", tc, w.Code)
		}
	}
}
