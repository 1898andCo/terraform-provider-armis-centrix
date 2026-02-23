// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package sweep

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// SweepArmisReports will delete all Armis reports with names starting with "tfacc".
func SweepArmisReports(name string) *resource.Sweeper {
	return &resource.Sweeper{
		Name: name,
		F: func(_ string) error {
			client, err := ConfigureSweeperClient(name)
			if err != nil {
				return fmt.Errorf("error configuring Armis client: %w", err)
			}
			if client == nil {
				return nil
			}

			ctx := context.Background()
			reports, err := client.GetReports(ctx)
			if err != nil {
				return fmt.Errorf("error listing Armis reports: %w", err)
			}

			prefix := "tfacc"
			for _, report := range reports {
				if strings.HasPrefix(report.ReportName, prefix) {
					log.Printf("[INFO] Deleting Armis report: %s", report.ReportName)
					_, err := client.DeleteReport(ctx, strconv.Itoa(report.ID))
					if err != nil {
						log.Printf("[ERROR] Failed to delete Armis report %s: %s", report.ReportName, err)
					}
				}
			}

			return nil
		},
	}
}
