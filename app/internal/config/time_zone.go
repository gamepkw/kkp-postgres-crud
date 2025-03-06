package config

import "time"

func SetTimeZone(locationName string) error {
	timeZone, err := time.LoadLocation(locationName)
	if err != nil {
		return err
	}
	time.Local = timeZone

	return nil
}
