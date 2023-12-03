/*
Copyright (c) 2023 Gsoc2 Inc.
*/
package cmd

import (
	"os"

	"github.com/Gsoc2/gsoc2-merge/packages/util"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:                   "reset",
	Short:                 "Used to delete all Gsoc2 related data on your machine",
	DisableFlagsInUseLine: true,
	Example:               "gsoc2 reset",
	Args:                  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// delete keyring item of current logged in user
		configFile, _ := util.GetConfigFile()

		// delete from keyring
		util.DeleteValueInKeyring(configFile.LoggedInUserEmail)

		// delete config
		_, pathToDir, err := util.GetFullConfigFilePath()
		if err != nil {
			util.HandleError(err)
		}

		os.RemoveAll(pathToDir)

		// delete secrets backup
		util.DeleteBackupSecrets()

		util.PrintSuccessMessage("Reset successful")
		Telemetry.CaptureEvent("cli-command:reset", posthog.NewProperties().Set("version", util.CLI_VERSION))
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
