package nextbus

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	agenciesCmd = &cobra.Command{
		Use:   "agencies",
		Short: "Gets a list of supported transit agencies by NextBus",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := client.GetAgencyList()
			if err != nil {
				return errors.Wrap(err, "failed to get agencies")
			}

			agencies := []agency{}
			for _, item := range list {
				agencies = append(agencies, agency{
					item.Tag,
					item.Title,
					item.RegionTitle,
				})
			}

			err = gocsv.Marshal(agencies, os.Stdout)
			if err != nil {
				return errors.Wrap(err, "failed to write results")
			}

			return nil
		},
	}
)

type agency struct {
	Agency string `csv:"Agency tag" json:"AgencyTag"`
	Title  string
	Region string
}
