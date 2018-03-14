package fury

import (
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
		Logger().Debug(cid)
		request.Header.Set(CID_HEADER, cid)
		next(rw, request)
	}
}

func RequestStatsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		url := request.URL
		path := url.Path
		qs := url.RawQuery

		Logger().Debug("---------***----------")
		Logger().Debugf("REMOTE ADDR: %s", request.RemoteAddr)
		Logger().Debugf("URL: %s, Method: %s", url, request.Method)
		Logger().Debugf("PATH: %s", path)
		Logger().Debugf("QUERY: %s", qs)
		Logger().Debugf("User Agent: %s", request.UserAgent())
		Logger().Debug("---------***----------")

		next(rw, request)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
