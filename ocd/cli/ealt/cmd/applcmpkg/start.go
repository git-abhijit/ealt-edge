/*
Copyright 2020 Huawei Technologies Co., Ltd.

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
package applcmpkg

import (
	"ealt/cmd/adapter"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
func NewApplcmStartCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "start",
		Short: "To start the application instance on MEP Node",
		Long:  `To start the application instance on MEP Node`,
		RunE: func(cmd *cobra.Command, args []string) error {
			theFlags := []string{cmd.Flag("appid").Value.String(),
				cmd.Flag("hostname").Value.String(),
				cmd.Flag("hostip").Value.String()}
			err := adapter.BuilderRequest(theFlags, "NewApplcmStartCommand")
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringP("appid", "i", "", "Application Instance ID to be started")
	cmd.Flags().StringP("hostname", "n", "", "Host Name where the App has to be instantiated")
	cmd.Flags().StringP("hostip", "o", "", "Host IP / MEP Node IP")
	cmd.MarkFlagRequired("appid")
	cmd.MarkFlagRequired("hostname")
	cmd.MarkFlagRequired("hostip")
	return cmd
}
