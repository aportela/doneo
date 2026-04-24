package cli

import (
	"flag"
)

type CommandLineParams struct {
	Verbose        bool
	InsertBulkData bool
}

func HandleFlags() (CommandLineParams, error) {
	params := CommandLineParams{}
	flag.BoolVar(&params.Verbose, "verbose", false, "Displays verbose output")

	flag.Parse()

	return params, nil
}
