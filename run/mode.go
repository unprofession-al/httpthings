package run

import "strings"

type mode int

const (
	ModeLocalServer mode = mode(iota) // Indicates that the server must me started as local web server.
	ModeAzureFunc                     // Indicates that the program is running as an Azure Function.
	ModeAWSLambda                     // Indicates that the program is running in an AWS Lambda context.
	ModeUnknown                       // Indicates that the mode provided does not exist.
)

var modeText = map[mode]string{
	ModeLocalServer: "localServer",
	ModeAzureFunc:   "azureFunc",
	ModeAWSLambda:   "awsLambda",
	ModeUnknown:     "unknown",
}

// NewMode creates a mode based on the string provided. In case the string provided cannot be
// mapped to a certain mode, ModeUnknown is returned.
func NewMode(in string) mode {
	in = strings.ToLower(in)
	for m, text := range modeText {
		if text == in {
			return mode(m)
		}
	}
	return ModeUnknown
}

// String returns a string representation of the mode.
func (m mode) String() string {
	return modeText[m]
}
