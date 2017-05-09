package fury

import (
	"github.com/op/go-logging"
	"os"
)

func Logger() *logging.Logger {
	logger := logging.MustGetLogger("fury")
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.7s} %{color:reset} %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	return logger
}
