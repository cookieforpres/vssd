package cmd

import (
	"os"
	"strconv"
	"strings"
	"vssd/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vssd",
	Short: "Virtual SSD command line interface",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func formatSizeString(size string) int {
	if strings.HasSuffix(size, "GB") {
		i, err := strconv.Atoi(strings.TrimSuffix(size, "GB"))
		if err != nil {
			return -1
		}

		return server.GB(i)
	} else if strings.HasSuffix(size, "KB") {
		i, err := strconv.Atoi(strings.TrimSuffix(size, "KB"))
		if err != nil {
			return -1
		}

		return server.KB(i)
	} else if strings.HasSuffix(size, "MB") {
		i, err := strconv.Atoi(strings.TrimSuffix(size, "MB"))
		if err != nil {
			return -1
		}

		return server.MB(i)
	} else if strings.HasSuffix(size, "B") {
		i, err := strconv.Atoi(strings.TrimSuffix(size, "B"))
		if err != nil {
			return -1
		}

		return server.B(i)
	} else {
		i, err := strconv.Atoi(size)
		if err != nil {
			return -1
		}

		return i
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
