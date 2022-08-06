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

func Run(mode string, listener string, handler http.Handler) error {
	switch strings.ToLower(mode) {
	case ModeLocal:
		if listener == "" {
			return fmt.Errorf("No listener defined")
		}
		fmt.Sprintf("Running locally at '%s'...\n", listener)
		return http.ListenAndServe(listener, handler)
	case ModeAzure:
		port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
		if !ok {
			return fmt.Errorf("Environment FUNCTIONS_CUSTOMHANDLER_PORT not defined")
		}
		listener := fmt.Sprintf(":%s", port)
		fmt.Sprintf("Running as Azure Function at '%s'...\n", listener)
		return http.ListenAndServe(listener, handler)
	case ModeLambda:
		fmt.Sprintf("Running as AWS Lambda...\n")
		return gateway.ListenAndServe(listener, handler)
	default:
		return fmt.Errorf("Unknown mode '%s'", mode)
	}
}
