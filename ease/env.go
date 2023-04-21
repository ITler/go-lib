package ease

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// GetAndUnsetEnv is following the security pattern DeleteOnFirstUse when it comes
// to reading environment variables.
//
// So the environment variable is read once and gets deleted
func GetAndUnsetEnv(key string) (result string, err error) {
	defer func() {
		if err = os.Unsetenv(key); err != nil {
			log.Fatal().Err(err)
		}
	}()
	if result, ok := os.LookupEnv(key); ok {
		return result, nil
	}
	return "", fmt.Errorf("ENV var '%s' not defined", key)
}

// LookupEnvVar returns the value of an environment variable if it exists,
// without caring about if the variable itself was defined or not.
func LookupEnvVar(key string) (value string) {
	value, defined := os.LookupEnv(key)
	if !defined {
		log.Debug().Msgf("Param '%v' not defined via ENV", key)
	}
	return value
}
