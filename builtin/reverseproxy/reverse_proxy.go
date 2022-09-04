package reverseproxy

// ref links:
// https://blog.csdn.net/mengxinghuiku/article/details/65448600
// https://github.com/ilanyu/ReverseProxy
import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gookit/gcli/v3"
)

type reverseProxy struct {
	listen   string
	remote   string
	remoteIP string
}

var rp = &reverseProxy{}
var dnsServers = []string{
	"114.114.114.114",
	"114.114.115.115",
	"119.29.29.29",
	"223.5.5.5",
	"8.8.8.8",
	"208.67.222.222",
	"208.67.220.220",
}

// ReverseProxyCmd command
var ReverseProxyCmd = &gcli.Command{
	Name: "proxy",
	Func: rp.Run,
	Desc: "start a reverse proxy http server",
	Config: func(c *gcli.Command) {
		c.StrOpt(
			&rp.listen,
			"listen", "s", "127.0.0.1:1180",
			"local proxy server listen address.",
		)
		c.StrOpt(
			&rp.remote,
			"remote", "r", "",
			"the remote reverse proxy server `address`. eg http://site.com:80;true",
		)
		c.StrOpt(
			&rp.remoteIP,
			"remoteIP", "", "",
			"the remote reverse proxy server IP address.",
		)

	},
}

func (rp *reverseProxy) Run(cmd *gcli.Command, args []string) error {
	if rp.remote == "" {
		return cmd.Errorf("must be setting the remote server by -r, --remote ")
	}

	urlObj, err := url.Parse(rp.remote)
	if err != nil {
		return err
	}

	rpHandler := ReverseProxy(urlObj)

	log.Printf("Listening on %s, forwarding to %s", rp.listen, rp.remote)
	log.Fatal(http.ListenAndServe(rp.listen, rpHandler))

	return nil
}

/*************************************************************
 * Reverse proxy
 *************************************************************/

// ReverseProxy create a global reverse proxy.
// Usage:
//
//	rp := ReverseProxy(&url.URL{
//		Scheme: "http",
//		Host:   "localhost:9091",
//	}, &url.URL{
//		Scheme: "http",
//		Host:   "localhost:9092",
//	})
//	log.Fatal(http.ListenAndServe(":9090", rp))
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

		fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

		// log.Println(r.RemoteAddr + " " + r.Method + " " + r.URL.String() + " " + r.Proto + " " + r.UserAgent())
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
