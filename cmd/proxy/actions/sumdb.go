package actions

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gomods/athens/pkg/paths"
)

func sumdbProxy(url *url.URL, sumdbProxyTo, nosumPatterns []string) http.Handler {
	rp := httputil.NewSingleHostReverseProxy(url)
	url = sumdbProxyWrapper(url, sumdbProxyTo)
	rp.Director = func(req *http.Request) {
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}
	if len(nosumPatterns) > 0 {
		return noSumWrapper(rp, url.Host, nosumPatterns)
	}
	return rp
}

func noSumWrapper(h http.Handler, host string, patterns []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/lookup/") {
			for _, p := range patterns {
				if paths.MatchesPattern(p, r.URL.Path[len("/lookup/"):]) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

func sumdbProxyWrapper(url *url.URL, sumdbProxyTo []string) *url.URL {
	for _, p := range sumdbProxyTo {
		strs := strings.Split(p, ":")
		if len(strs) != 2 {
			log.Printf("can't format sumdbProxyTo '%s'\n", p)
			continue
		}
		if url.Host == strs[0] {
			url.Host = strs[1]
		}
	}
	return url
}
