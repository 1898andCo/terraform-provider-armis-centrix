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

// SweepArmisCollectors returns a sweeper function to clean up Armis collectors.
func SweepArmisCollectors(name string) *resource.Sweeper {
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

			// Get all collectors
			ctx := context.Background()
			collectors, err := client.GetCollectors(ctx)
			if err != nil {
				return fmt.Errorf("error getting Armis collectors: %w", err)
			}

			prefix := "test."
			for _, collector := range collectors {
				if strings.HasPrefix(collector.Name, prefix) {
					log.Printf("[INFO] Deleting Armis collector: %s", collector.Name)
					_, err := client.DeleteCollector(ctx, strconv.Itoa(collector.CollectorNumber))
					if err != nil {
						log.Printf("[ERROR] Failed to delete Armis collector %s: %s", collector.Name, err)
					}
				}
			}

			return nil
		},
	}
}
