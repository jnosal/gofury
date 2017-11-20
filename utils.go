package fury

import "errors"

const (
	MAX_RETRIES = 10
)

type RetriableFunc func(attempt int) (retry bool, err error)

func Retry(fn RetriableFunc) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > MAX_RETRIES {
			return errors.New("Max retries exceeded")
		}
	}
	return err
}
