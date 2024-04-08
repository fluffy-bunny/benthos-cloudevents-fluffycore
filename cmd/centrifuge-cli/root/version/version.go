package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version global
var Version string

func init() {
	SetVersion("0.0.0")
}

// SetVersion ...
func SetVersion(version string) {
	Version = version
}

// Init command
func Init(rootCmd *cobra.Command) {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the app",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
	rootCmd.AddCommand(command)
}
