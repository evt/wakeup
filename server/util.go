package server

import (
	"time"

	"github.com/pkg/errors"
)

// CallTimeFormat is a call time format
const CallTimeFormat = "15:04"

// validateCallTime returns error if provided time string is not in format xx:yy
func validateCallTime(t string) error {
	_, err := time.Parse(CallTimeFormat, t)
	return err
}

func addRetryPeriod(callTime, retryPeriod string) (string, error) {
	if callTime == "" {
		return "", errors.New("No call time provided")
	}
	if retryPeriod == "" {
		return "", errors.New("No call time provided")
	}
	ct, err := time.Parse(CallTimeFormat, callTime)
	if err != nil {
		return "", errors.Wrap(err, "addRetryPeriod->time.Parse")
	}
	retryPeriodDuration, err := time.ParseDuration(retryPeriod)
	if err != nil {
		return "", errors.Wrap(err, "addRetryPeriod->time.ParseDuration")
	}
	return ct.Add(retryPeriodDuration).Format(CallTimeFormat), nil
}
