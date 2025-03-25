package config

import (
	"errors"
)

type OptionsValidator struct {
}

func (vld *OptionsValidator) Validate(opts *AppConfig) error {
	runesBaseURL := []rune(opts.RunAddress)

	if string(runesBaseURL[len(runesBaseURL)-1:]) == "/" {
		return errors.New("incorrect -a argument, don't put a slash at the end")
	}

	return nil
}
