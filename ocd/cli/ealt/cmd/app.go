/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	appCmds "ealt/cmd/appm"

	"github.com/spf13/cobra"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "The command is used for Application Management in the EALT Edge System.",
	Long: `The command is used for Application Management in the EALT Edge System.
	It has options: create , info and delete`,
}

func init() {
	appCmd.AddCommand(appCmds.NewAppCreateCommand())
	appCmd.AddCommand(appCmds.NewAppDeleteCommand())
	appCmd.AddCommand(appCmds.NewAppInfoCommand())

	rootCmd.AddCommand(appCmd)
}
