// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/1898andCo/armis-sdk-go/v2/armis"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func appendAPIError(diags *diag.Diagnostics, title string, err error) {
	if err == nil {
		diags.AddError(title, "(no additional error details returned)")
		return
	}

	var apiErr *armis.APIError
	if errors.As(err, &apiErr) {
		body := strings.TrimSpace(string(apiErr.Body))
		if body == "" {
			body = "(empty response body)"
		} else if json.Valid(apiErr.Body) {
			var pretty bytes.Buffer
			if indentErr := json.Indent(&pretty, apiErr.Body, "", "  "); indentErr == nil {
				body = pretty.String()
			}
		}

		diags.AddError(
			title,
			fmt.Sprintf("API error %d %s\nResponse body:\n%s", apiErr.StatusCode, http.StatusText(apiErr.StatusCode), body),
		)
		return
	}

	diags.AddError(title, fmt.Sprintf("API error: %v", err))
}
