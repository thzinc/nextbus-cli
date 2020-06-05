package nextbus

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	api "github.com/dinedal/nextbus"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	watch      bool
	vehicleCmd = &cobra.Command{
		Use:   "vehicle",
		Short: "Vehicle-related information",
	}
	vehicleLocationCmd = &cobra.Command{
		Use:   "locations",
		Short: "Gets a list of vehicle locations for a transit agency",
		RunE: func(cmd *cobra.Command, args []string) error {
			signals := make(chan os.Signal)
			defer close(signals)

			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

			didWriteHeaders := false
			lastTimeString := "0"
			for {
				opts := []api.VehicleLocationParam{api.VehicleLocationTime(lastTimeString)}
				if routeTag != "" {
					opts = append(opts, api.VehicleLocationRoute(routeTag))
				}
				list, err := client.GetVehicleLocations(agencyTag, opts...)
				if err != nil {
					return errors.Wrap(err, "failed to get vehicle locations")
				}
				lastTimeString = list.LastTime.Time

				lastTimeSeconds, err := strconv.ParseInt(lastTimeString, 10, 64)
				if err != nil {
					return errors.Wrap(err, "failed to parse seconds from last time")
				}

				asOf := time.Unix(lastTimeSeconds/1000, lastTimeSeconds%1000)

				locations := []vehicleLocation{}
				for _, location := range list.VehicleList {
					locations = append(locations, vehicleLocation{
						ID:              location.ID,
						RouteTag:        location.RouteTag,
						DirTag:          location.DirTag,
						Lon:             location.Lon,
						Lat:             location.Lat,
						SecsSinceReport: location.SecsSinceReport,
						Predictable:     location.Predictable,
						Heading:         location.Heading,
						SpeedKmHr:       location.SpeedKmHr,
						AsOf:            asOf,
					})
				}

				marshal := gocsv.Marshal
				if didWriteHeaders {
					marshal = gocsv.MarshalWithoutHeaders
				}
				err = marshal(locations, os.Stdout)
				if err != nil {
					return errors.Wrap(err, "failed to write results")
				}

				didWriteHeaders = true

				if !watch {
					return nil
				}

				select {
				case <-time.After(5 * time.Second):
				case <-signals:
					return nil
				}
			}
		},
	}
)

func init() {
	vehicleCmd.AddCommand(vehicleLocationCmd)
	vehicleLocationCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch")
	vehicleLocationCmd.Flags().StringVarP(&agencyTag, "agency", "a", "", "Agency tag")
	vehicleLocationCmd.Flags().StringVarP(&routeTag, "route", "r", "", "Route tag")
	vehicleLocationCmd.MarkFlagRequired("agency")
}

type vehicleLocation struct {
	ID              string `json:"id"`
	RouteTag        string `csv:"Route tag"`
	DirTag          string `csv:"Direction tag"`
	Lon             string `csv:"Longitude"`
	Lat             string `csv:"Latitude"`
	SecsSinceReport string `csv:"Seconds since last report"`
	Predictable     string
	Heading         string
	SpeedKmHr       string    `csv:"Speed (km/hr)"`
	AsOf            time.Time `csv:"As of"`
}
