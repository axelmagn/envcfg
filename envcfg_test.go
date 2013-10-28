// Copyright 2012 Axel Magnuson.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Tests for envcfg package
package envcfg_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/axelmagn/envcfg"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

var setValue string = "TEST"

var settings map[string]string

func randString() string {
	randInt := rand.Int63()
	randBytes := make([]byte, 10)
	binary.PutVarint(randBytes, randInt)
	return hex.EncodeToString(randBytes)
}

func TestSetenvControl(t *testing.T) {
	envKey := envcfg.ENV_PREFIX + randString()
	err := os.Setenv(envKey, setValue)
	if err != nil {
		t.Errorf("Error while setting Env Variable %s to %s: %s", envKey, setValue, err.Error())
	}
	envValue := os.Getenv(envKey)
	if envValue != setValue {
		t.Errorf("Extracted Env Variable %s had value %s.  Expected %s.", envKey, envValue, setValue)
	}

}

func TestExtractEnvIfPrefix(t *testing.T) {
	// defined env variable
	envKey := envcfg.ENV_PREFIX + randString()
	err := os.Setenv(envKey[len(envcfg.ENV_PREFIX):], setValue)
	if err != nil {
		t.Errorf("Error while setting Env Variable %s to %s: %s", envKey, setValue, err.Error())
	}
	envValue, prefixPresent := envcfg.ExtractEnvIfPrefix(envKey, envcfg.ENV_PREFIX)
	if envValue != setValue || prefixPresent != true {
		t.Errorf("Extracted Env Variable %s had value %s, %b.  Expected %s, %b.", envKey, envValue, prefixPresent, setValue, true)
	}
}

func TestExtractEnvIfPrefix_UndefinedEnv(t *testing.T) {
	// set envKey as a random alphanumeric string
	envKey := envcfg.ENV_PREFIX + randString()
	envValue, prefixPresent := envcfg.ExtractEnvIfPrefix(envKey, envcfg.ENV_PREFIX)
	if envValue != "" || prefixPresent != true {
		t.Errorf("Extracted Env Variable %s had value %s.  Expected \"\" for undefined env variable.", envKey, envValue)
	}
}

func TestExtractEnvIfPrefix_NoPrefix(t *testing.T) {
	// set envKey as a random alphanumeric string
	envKey := randString()
	envValue, prefixPresent := envcfg.ExtractEnvIfPrefix(envKey, envcfg.ENV_PREFIX)
	if envValue != "" || prefixPresent != false {
		t.Errorf("Extracted Env Variable %s had value %s.  Expected nil for absent prefix.", envKey, envValue)
	}
}

// we use this as the initializer for the settings variable
// I don't know if that's good practice or not.
func TestReadSettings(t *testing.T) {
	var err error
	var cfgFile *os.File
	// Create sample settings data
	rawSettings := `
	# Key flag
	KEY_FLAG

	# Value Literal
	VLKEY	VLVALUE

	# Env Key
	EKKEY	ENV:ENVCFG_TEST_ENV_KEY_VALUE

	# Env Key Default Where Env Variable is defined
	EKDKEY	ENV:ENVCFG_TEST_ENV_DEFINED		ekd_default

	# Env Key Default Where Env Variable is not defined
	EKUKEY	ENV:ENVCFG_TEST_ENV_UNDEFINED	eku_default
	`

	// write sample to temp file
	cfgFile, err = ioutil.TempFile(os.TempDir(), "envcfg_test_config.ecfg")
	if err != nil {
		t.Errorf("Error while creating temporary settings file: %s", err.Error())
	}
	_, err = cfgFile.WriteString(rawSettings)
	if err != nil {
		t.Errorf("Error while writing settings to temp file: %s", err.Error())
	}
	err = cfgFile.Close()
	if err != nil {
		t.Errorf("Error while closing temp file: %s", err.Error())
	}
	cfgFileName := cfgFile.Name()

	// open sample file
	cfgFile, err = os.Open(cfgFileName)
	if err != nil {
		t.Errorf("Error while opening temp settings file for reading: %s", err.Error())
	}

	// configure ENV for different settings
	os.Setenv("ENVCFG_TEST_ENV_KEY_VALUE", "ek_value")
	os.Setenv("ENVCFG_TEST_ENV_DEFINED", "ekd_value")

	// read settings
	settings, err = envcfg.ReadSettings(cfgFile)
	if err != nil {
		t.Errorf("ReadSettings failed with error: %s", err.Error())
	}
}

func ExampleReadSettings_KeyFlag() {
	fmt.Println(settings["KEY_FLAG"])
	// Output: 1
}

func ExampleReadSettings_ValueLiteral() {
	fmt.Println(settings["VLKEY"])
	// Output: VLVALUE
}

func ExampleReadSettings_EnvKey() {
	fmt.Println(settings["EKKEY"])
	// Output: ek_value
}

func ExampleReadSettings_EnvKeyDefault() {
	fmt.Println(settings["EKDKEY"])
	// Output: ekd_value
}

func ExampleReadSettings_EnvKeyDefaultUndefined() {
	fmt.Println(settings["EKUKEY"])
	// Output: eku_default
}
