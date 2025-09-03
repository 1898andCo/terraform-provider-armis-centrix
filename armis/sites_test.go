// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"testing"
)

func TestGettingSites(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetSites(context.Background())
	if err != nil {
		t.Fatalf("get sites: %v", err)
	}
	prettyPrint(res)
}
