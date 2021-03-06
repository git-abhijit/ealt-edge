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
package init

import (
	"ealt/cmd/setup"

	"github.com/spf13/cobra"
)

// allCmd represents the all command
func NewK8SCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "k8s",
		Short: "Command to uninstall K8S cluster on MECM Node",
		Long:  `Command to uninstall K8S cluster on MECM Node`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := setup.EaltReset("k8s")
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
