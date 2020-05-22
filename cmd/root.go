package cmd

import (
	"fmt"
	"os"

	"github.com/hectcastro/gjson/pkg/gjson"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GeoJSONIOURL contains the URL for geojson.io.
const GeoJSONIOURL = "https://geojson.io"

var rootCmd = &cobra.Command{
	Use:              "gjson",
	Short:            "Ship files from your local system to geojson.io for visualization and editing.",
	PersistentPreRun: setupLogging,
	Run:              run,
	SilenceUsage:     true,
}

// Execute is the entrypoint for the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("print", "p", false, "print the URL, rather than opening it")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func run(cmd *cobra.Command, args []string) {
	var file *os.File

	if stdinPopulated() {
		file = os.Stdin
	} else if len(args) > 0 {
		var err error

		file, err = os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No input GeoJSON data detected")
	}

	geojson := &gjson.GeoJSON{}
	err := geojson.Unmarshal(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if geojson.Size() > gjson.BigGeoJSONLength {
		log.Fatal("This file will likely display slowly on geojson.io")
	}

	printURL, _ := cmd.Flags().GetBool("print")
	if geojson.Size() <= gjson.MaxURLLength {
		displayResource("#data=data:application/json,", geojson.ToURLEncoded(), printURL)
	} else {
		displayResource("#id=gist:/", geojson.ToGist(), printURL)
	}
}

func setupLogging(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func stdinPopulated() bool {
	stat, _ := os.Stdin.Stat()

	return (stat.Mode() & os.ModeCharDevice) == 0
}

func displayResource(contentType string, content string, print bool) {
	url := fmt.Sprintf("%s/%s%s", GeoJSONIOURL, contentType, content)

	if print {
		fmt.Println(url)
	} else {
		browser.OpenURL(url)
	}
}
