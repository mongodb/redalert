// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package checks

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Args is a convenience type to express the "args" key in the test yaml
type Args map[string]interface{}

// ArgableFunc is a function which modifies arguments before passing them to an
// Argable
type ArgableFunc func(args Args) (Checker, error)

// Argable is any struct which can create a Checker from the YAML args we get
// back from a test block.
type Argable interface {
	FromArgs(args Args) (Checker, error)
}

// Checker is any struct that performs a system check
type Checker interface {
	Check() error
}

// RequiredArgError is an error which indicates a required arg was not given
type RequiredArgError struct {
	RequiredArg  string
	ProvidedArgs Args
}

func (rae RequiredArgError) Error() string {
	return fmt.Sprintf("%s is a required arg and was not given: provided args were: %v", rae.RequiredArg, rae.ProvidedArgs)
}

// IsRequiredArg returns a boolean indicating if the given err is a
// RequiredArgError
func IsRequiredArg(err error) bool {
	switch err.(type) {
	case RequiredArgError:
		return true
	case *RequiredArgError:
		return true
	default:
		return false
	}
}

func requiredArgs(args Args, requiredArgs ...string) error {
	for _, requiredArg := range requiredArgs {
		if _, ok := args[requiredArg]; !ok {
			return RequiredArgError{requiredArg, args}
		}
	}

	return nil
}

func decodeFromArgs(args Args, into interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      into,
	})

	if err != nil {
		return fmt.Errorf("Unable to create decoder %v: %s", args, err)
	}

	err = decoder.Decode(args)
	if err != nil {
		return fmt.Errorf("Unable to decode %v: %s", args, err)
	}

	return err
}
