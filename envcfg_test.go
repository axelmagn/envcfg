package envcfg_test

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"os"
	"testing"
	"github.com/axelmagn/envcfg"
)

var setValue string = "TEST"

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

func TestExtractEnvIfPrefixUndefinedEnv(t *testing.T) {
	// set envKey as a random alphanumeric string
	envKey := envcfg.ENV_PREFIX + randString()
	envValue, prefixPresent := envcfg.ExtractEnvIfPrefix(envKey, envcfg.ENV_PREFIX)
	if envValue != "" || prefixPresent != true {
		t.Errorf("Extracted Env Variable %s had value %s.  Expected \"\" for undefined env variable.", envKey, envValue)
	}
}

func TestExtractEnvIfPrefixNoPrefix(t *testing.T)	{
	// set envKey as a random alphanumeric string
	envKey := randString()
	envValue, prefixPresent := envcfg.ExtractEnvIfPrefix(envKey, envcfg.ENV_PREFIX)
	if envValue != "" || prefixPresent != false {
		t.Errorf("Extracted Env Variable %s had value %s.  Expected nil for absent prefix.", envKey, envValue)
	}
}
