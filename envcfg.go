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
//
// // Lines beginning with octothorpe are comments
// "#" COMMENT
package envcfg

import (
	"bufio"
	"errors"
	"io"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// regex for splitting a line by one or more spaces
var re_spaces *regexp.Regexp = regexp.MustCompile("\\s+")
var max_tokens int = 3

var ENV_PREFIX string = "ENV:"
var COMMENT_PREFIX string = "#"

// value literals
var TRUE string = "1"

// Read settings and store them in a Settings map
func ReadSettings(reader io.Reader) (map[string]string, error) {
	settings := make(map[string]string)
	lineScanner := bufio.NewScanner(reader)

	// each line is a settings value
	for lineScanner.Scan() {
		line := strings.Trim(lineScanner.Text(), " \t")

		// skip comments
		if strings.HasPrefix(line, COMMENT_PREFIX) {
			continue
		}

		// split line by spaces
		tokens := re_spaces.Split(line, max_tokens)

		switch len(tokens) {

		case 1:
			// a variable flag.  Set this value to TRUE
			settings[tokens[0]] = TRUE

		case 2:
			// a key-value pair
			key := tokens[0]
			value := tokens[1]

			// get env if prefix
			envValue, prefixPresent := ExtractEnvIfPrefix(value, ENV_PREFIX)
			// empty string indicates undefined env variable.  Config specifies
			// no default,  so it's required and we should throw an error
			if envValue == "" && prefixPresent {
				err := errors.New("Environment variable is undefined: " + value)
				return nil, err
			}

			if !prefixPresent {
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
			envValue, prefixPresent := ExtractEnvIfPrefix(envKey, ENV_PREFIX)
			// if nil, then there wasn't a prefix, meaning it was an illegal
			// triple
			if !prefixPresent && envValue == "" {
				err := errors.New(fmt.Sprintf("Default provided for a literal string: %s (%v)", line, tokens))
				return nil, err
			} else if envValue == "" {
				envValue = valueDefault
			}

			settings[key] = envValue
		}
	}

	return settings, nil
}

// Get an env variable if the key starts with a prefix.  Returns the extracted 
// variable, as well as a boolean indicating whether the prefix was present.
func ExtractEnvIfPrefix(envKey string, envPrefix string) (string, bool) {
	var prefix string
	// check for environment prefix
	prefixLen := len(envPrefix)
	if len(envKey) >= prefixLen {
		prefix = envKey[0:prefixLen]
	} else {
		prefix = ""
	}

	// get env variable if prefix present
	if prefix == envPrefix {
		return os.Getenv(envKey[prefixLen:]), true
	} else {
		return "", false
	}

}
