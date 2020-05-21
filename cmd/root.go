package cmd

import (
	"io/ioutil"

	geojson "github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:              "gjson",
	Short:            "Ship files from your local system to geojson.io for visualization and editing.",
	Args:             cobra.MaximumNArgs(1),
	PersistentPreRun: setupLogging,
	Run:              run,
	SilenceUsage:     true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func initConfig() {
	viper.AutomaticEnv()
}

func run(cmd *cobra.Command, args []string) {
	log.Debug(args)

	var filename string = args[0]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)

	for _, feature := range fc.Features {
		log.Debug(feature.Geometry)
	}
}

func setupLogging(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
