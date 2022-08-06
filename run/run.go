package run

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/apex/gateway"
)

const (
	ModeLocal  = "local"
	ModeAzure  = "azurefunc"
	ModeLambda = "awslambda"
)

func Run(mode string, listener string, handler http.Handler, log bool) error {
	switch strings.ToLower(mode) {
	case ModeLocal:
		if listener == "" {
			return fmt.Errorf("No listener defined")
		}
		if log {
			fmt.Printf("Running locally at 'http://%s'...\n", listener)
		}
		return http.ListenAndServe(listener, handler)
	case ModeAzure:
		port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
		if !ok {
			return fmt.Errorf("Environment FUNCTIONS_CUSTOMHANDLER_PORT not defined")
		}
		listener := fmt.Sprintf(":%s", port)
		if log {
			fmt.Printf("Running as Azure Function at '%s'...\n", listener)
		}
		return http.ListenAndServe(listener, handler)
	case ModeLambda:
		if log {
			fmt.Printf("Running as AWS Lambda...\n")
		}
		return gateway.ListenAndServe(listener, handler)
	default:
		return fmt.Errorf("Unknown mode '%s'", mode)
	}
}

func DetectRunMode() string {
	if _, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		return ModeAzure
	} else if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		return ModeLambda
	} else {
		return ModeLocal
	}
}
