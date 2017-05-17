package fury

import (
	"github.com/jnosal/gofury"
	"math/rand"
	"net/http"
	"time"
)

const (
	CID_HEADER = "X-Request-CID"
)

func generateCID() string {
	alphabet := []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alphabetLength := len(alphabet)

	cid := make([]rune, alphabetLength)

	for i := range cid {
		cid[i] = alphabet[rand.Intn(alphabetLength)]
	}
	return string(cid)
}

func RequestCIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		cid := generateCID()
		fury.Logger().Debug(cid)
		request.Header.Set(CID_HEADER, cid)
		next(rw, request)
	}
}

func RequestStatsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		url := request.URL
		path := url.Path
		qs := url.RawQuery

		fury.Logger().Debug("---------***----------")
		fury.Logger().Debug("REMOTE ADDR: %s", request.RemoteAddr)
		fury.Logger().Debugf("URL: %s, Method: %s", url, request.Method)
		fury.Logger().Debugf("PATH: %s", path)
		fury.Logger().Debugf("QUERY: %s", qs)
		fury.Logger().Debugf("User Agent: %s", request.UserAgent())
		fury.Logger().Debug("---------***----------")

		next(rw, request)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
