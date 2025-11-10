// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

type ListSettings struct {
	CreatedBy      string `json:"created_by"`
	CreationTime   string `json:"creation_time"`
	Description    string `json:"description"`
	LastUpdateTime string `json:"last_update_time"`
	LastUpdatedBy  string `json:"last_updated_by"`
	ListID         int    `json:"list_id"`
	ListName       string `json:"list_name"`
	ListType       string `json:"list_type"`
}

// GetLists represents API response which comes back as data and success.
type GetLists struct {
	Data    []ListSettings `json:"data"`
	Success bool           `json:"success,omitempty"`
}
