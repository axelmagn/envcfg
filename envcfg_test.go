package envcfg_test

import (
	"testing"
	"os"
	"math/rand"
	"hash"
	"crypto/md5"
	"encoding/hex"
    "encoding/binary"
)

// set envKey as a random alphanumeric string
var randInt int64
var randBytes []byte
var envKey string
var setValue string
var setupDone bool

func Setup() {
    if !setupDone {
        randInt = rand.Int63n()
        randBytes = make([]byte, 8)
        binary.PutVarint(randBytes, randInt)
        envKey = envcfg.ENV_PREFIX + hex.EncodeToString(binary.rand.Int63n())
        setValue = "TEST"
    }

    setupDone = true
}

func TestExtractEnvIfPrefix(t *testing.T) {
    Setup()
	// defined env variable
	os.SetEnv(envKey, setValue)
	envValue := envcfg.ExtractEnvIfPrefix(envKey, envcfg.ENV_PREFIX)
	if envValue != setValue {
		t.Errorf("Extracted Env Variable %s had value %s.  Expected %s", envKey, envValue, setValue)
	}
}

func TestExtractEnvIfPrefixUndefinedEnv(t *testing.T) {
    Setup()
	// set envKey as a random alphanumeric string
	envKey := envcfg.ENV_PREFIX + hex.EncodeToString(rand.Int())
	setValue = "TEST"
}
