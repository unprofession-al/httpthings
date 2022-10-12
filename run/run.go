package run

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

// Run starts a web server in the mode provided. log is a function that takes only a string
// so you can bring your own logging
func Run(mode mode, listener string, handler http.Handler, log func(string)) error {
	switch mode {
	case ModeLocalServer:
		if listener == "" {
			return fmt.Errorf("no listener defined")
		}
		log(fmt.Sprintf("Running locally at 'http://%s'...\n", listener))
		return http.ListenAndServe(listener, handler)
	case ModeAzureFunc:
		port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
		if !ok {
			return fmt.Errorf("environment FUNCTIONS_CUSTOMHANDLER_PORT not defined")
		}
		listener := fmt.Sprintf(":%s", port)
		log(fmt.Sprintf("Running as Azure Function at '%s'...\n", listener))
		return http.ListenAndServe(listener, handler)
	case ModeAWSLambda:
		log("Running as AWS Lambda...\n")
		return gateway.ListenAndServe(listener, handler)
	default:
		return fmt.Errorf("unknown mode")
	}
}

// DetectRunMode tries to detect the run mode based on the environment variables present
// at launch time. The order is:
//  1. If `FUNCTIONS_CUSTOMHANDLER_PORT` is found, it is assumed that the function is
//     started in an Azure Fuctions context.
//  2. If `AWS_LAMBDA_FUNCTION_NAME` is found, it is assumed that the function is started
//     in an AWS Lambda context.
//  3. Everything else indicates that the software is requested to be started as a regular
//     web server.
func DetectRunMode() mode {
	if _, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		return ModeAzureFunc
	} else if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		return ModeAWSLambda
	} else {
		return ModeLocalServer
	}
}
