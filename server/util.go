package server

import "time"

// CallTimeFormat is a call time format
const CallTimeFormat = "15:04"

// validateCallTime returns error if provided time string is not in format xx:yy
func validateCallTime(t string) error {
	_, err := time.Parse(CallTimeFormat, t)
	return err
}
