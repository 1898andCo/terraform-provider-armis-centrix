// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"testing"
)

func TestGettingBoundaries(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	res, err := client.GetBoundaries(context.Background())
	if err != nil {
		t.Fatalf("get boundaries: %v", err)
	}
	prettyPrint(res)
}
