package domain

import "errors"

type Confidence float64

func (c Confidence) Validate() error {
	if c < 0 || c > 1 {
		return errors.New("confidence must be between 0 and 1")
	}
	return nil
}

func (c Confidence) AsPercentage() int {
	return int(c * 100)
}