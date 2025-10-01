package armis

import (
	"context"
	"testing"
)

func TestGettingSearch(t *testing.T) {
	t.Parallel()
	client := integrationClient(t)

	const aql = "in:alerts status:Open timeFrame:\"7 Days\""

	res, err := client.GetSearch(context.Background(), aql, true, true)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}

	prettyPrint(res)

	if res.Total == 0 {
		t.Logf("no results returned for AQL %q", aql)
	}
}
