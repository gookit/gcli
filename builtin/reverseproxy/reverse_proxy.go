package reverseproxy

// ref links:
// https://blog.csdn.net/mengxinghuiku/article/details/65448600
// https://github.com/ilanyu/ReverseProxy
import (
	"github.com/gookit/cliapp"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxyCommand create command
func ReverseProxyCommand() *cliapp.Command {
	c := &cliapp.Command{
		Name: "proxy",
		Func: func(cmd *cliapp.Command, args []string) int {
			return 0
		},

		Description: "start a reverse proxy http server",
	}

	return c
}

/*************************************************************
 * Reverse proxy
 *************************************************************/

// ReverseProxy create a global reverse proxy.
// usage:
// 	rp := ReverseProxy(&url.URL{
// 		Scheme: "http",
// 		Host:   "localhost:9091",
// 	}, &url.URL{
// 		Scheme: "http",
// 		Host:   "localhost:9092",
// 	})
// 	log.Fatal(http.ListenAndServe(":9090", rp))
func ReverseProxy(targets ...*url.URL) *httputil.ReverseProxy {
	if len(targets) == 0 {
		panic("Please add at least one remote target server")
	}

	var target *url.URL

	// if only one target
	if len(targets) == 1 {
		target = targets[0]
	}

	director := func(req *http.Request) {
		if len(targets) > 1 {
			target = targets[rand.Int()%len(targets)]
		}

		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		// req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{Director: director}
}
