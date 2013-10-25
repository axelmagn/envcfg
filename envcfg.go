// Copyright 2012 Axel Magnuson.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Library for reading .ecfg environment config files.  These files fulfill the
// function of configurations that use both file-coded values and environment
// variables.  This way sensitive information can be delegated to environment
// variables, and satisfy configuration best practices laid out by 12factor.net
//
// Each non-blank line should take the form of one of the following directives:
//
//
// // Settings["KEY_FLAG"] == envcfg.TRUE
//
// KEY_FLAG
//
// // Settings["KEY"] == LITERAL
//
// KEY			LITERAL
//
// // Settings["KEY"] == os.GetEnv(ENV_KEY). Throws error if env variable isn't set.
//
// KEY			"ENV:"ENV_KEY
//
// // Settings["KEY"] == os.GetEnv(ENV_KEY). Defaults to DEFAULT if env variable isn't set.
//
// KEY			"ENV:"ENV_KEY	DEFAULT
package envcfg

import (
	"errors"
	"io"
	"os"
	"regexp"
)

// regex for splitting a line by one or more spaces
var re_spaces *regexp.Regexp = regexp.MustCompile("\\s+")
var max_tokens int = 3

var ENV_PREFIX string = "ENV:"

// value literals
var TRUE string = "1"

type Settings map[string]string

// Read settings and store them in a Settings map
func ReadSettings(reader io.Reader) (*Settings, error) {
	settings := make(map[string]string)
	lineScanner := NewScanner(reader)

	// each line is a settings value
	for lineScanner.Scan() {
		line := lineScanner.Text()

		// split line by spaces
		tokens = re_spaces.Split(line, max_tokens)

		switch len(tokens) {

		case 1:
			// a variable flag.  Set this value to TRUE
			settings[tokens[0]] = TRUE

		case 2:
			// a key-value pair
			key := tokens[0]
			value := tokens[1]

			// get env if prefix
			envValue := ExtractEnvIfPrefix(value, ENV_PREFIX)
			// empty string indicates undefined env variable.  Config specifies
			// no default,  so it's required and we should throw an error
			if envValue == "" {
				err := errors.New("Environment variable is undefined:\t" + value)
				return nil, err
			}

			if envValue == nil {
				settings[key] = value
			} else {
				settings[key] = envValue
			}

		case 3:
			// a key-envKey-default triple
			key := tokens[0]
			envKey := tokens[1]
			valueDefault := tokens[2]

			// get env if prefix
			envValue := ExtractEnvIfPrefix(envKey, ENV_PREFIX)
			// if nil, then there wasn't a prefix, meaning it was an illegal
			// triple
			if envValue == nil {
				err := errors.New("Default provided for a literal string:\t" + line)
			} else if envValue == "" {
				envValue = valueDefault
			}

			settings[key] = envValue
		}
	}
}

// Get an env variable if the key starts with a prefix. It returns a nil if
// the prefix is not present, otherwise it returns the environment variable.
// If the prefix is present but the variable is undefined, it returns an empty
// string.
func ExtractEnvIfPrefix(envKey string, envPrefix string) string {

	// check for environment prefix
	prefixLen := len(envPrefix)
	if len(envKey) >= prefix_len {
		prefix := envKey[0:prefix_len]
	} else {
		prefix := nil
	}

	// get env variable if prefix present
	if prefix == envPrefix {
		return os.GetEnv(envKey[prefixLen:])
	} else {
		return nil
	}

}
