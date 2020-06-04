package nextbus

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	api "github.com/thzinc/nextbus"
)

type outputType string

var (
	outputTypeCSV  outputType = "csv"
	outputTypeJSON outputType = "json"
)

var (
	client    *api.Client = api.DefaultClient
	agencyTag string
	routeTag  string
	output    string
	// RootCmd is the command line interface to interact with the NextBus API
	RootCmd = &cobra.Command{
		Use:   "nextbus",
		Short: "Command line interface to interact with the NextBus API",
		RunE: func(cmd *cobra.Command, args []string) error {
			client.GetVehicleLocations(agencyTag)
			return errors.New("not implemented")
		},
	}
)

func init() {
	RootCmd.AddCommand(agenciesCmd)
	RootCmd.AddCommand(vehicleCmd)
}

// Execute the root command line interface
func Execute() error {
	return RootCmd.Execute()
}
