// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/1898andCo/armis-sdk-go/v2/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &reportsDataSource{}
	_ datasource.DataSourceWithConfigure = &reportsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *reportsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*armis.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *armis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// ReportsDataSource is a helper function to simplify the provider implementation.
func ReportsDataSource() datasource.DataSource {
	return &reportsDataSource{}
}

// reportsDataSource is the data source implementation.
type reportsDataSource struct {
	client *armis.Client
}

// Metadata returns the data source type name.
func (d *reportsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reports"
}

// Schema defines the schema for the reports data source.
func (d *reportsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Armis report information. Supports filtering by report_id or report_name. If report_id is provided, it takes precedence and fetches a single report directly from the API. If only report_name is provided, all reports are fetched and filtered client-side by name (case-sensitive exact match). If no filter is provided, all reports are returned.",
		Attributes: map[string]schema.Attribute{
			"report_id": schema.StringAttribute{
				Description: "Optional report ID to fetch a specific report. Takes precedence over report_name if both are provided. Note: The report ID is stored as a number in the results, but provided as a string for filtering.",
				Optional:    true,
			},
			"report_name": schema.StringAttribute{
				Description: "Optional report name to filter reports (case-sensitive exact match). Uses client-side filtering after fetching all reports. Ignored if report_id is provided.",
				Optional:    true,
			},
			"reports": schema.ListNestedAttribute{
				Description: "A computed list of reports. Each object in the list contains detailed information about a report.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.NumberAttribute{
							Description: "A unique identifier for the report.",
							Computed:    true,
						},
						"report_name": schema.StringAttribute{
							Description: "The name of the report.",
							Computed:    true,
						},
						"report_type": schema.StringAttribute{
							Description: "The type of the report.",
							Computed:    true,
						},
						"asq": schema.StringAttribute{
							Description: "The ASQ (Armis Standard Query) used by the report.",
							Computed:    true,
						},
						"creation_time": schema.StringAttribute{
							Description: "The timestamp when the report was created.",
							Computed:    true,
						},
						"is_scheduled": schema.BoolAttribute{
							Description: "Indicates whether the report is scheduled.",
							Computed:    true,
						},
						"schedule": schema.SingleNestedAttribute{
							Description: "The schedule configuration for the report. Only present when is_scheduled is true.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"email": schema.ListAttribute{
									Description: "List of email addresses to receive the scheduled report.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"repeat_amount": schema.NumberAttribute{
									Description: "The repeat interval amount for the scheduled report. Can be a decimal value.",
									Computed:    true,
								},
								"repeat_unit": schema.StringAttribute{
									Description: "The repeat interval unit for the scheduled report (e.g., 'days', 'weeks', 'months').",
									Computed:    true,
								},
								"report_file_format": schema.StringAttribute{
									Description: "The file format of the scheduled report (e.g., 'pdf', 'csv', 'xlsx').",
									Computed:    true,
								},
								"time_of_day": schema.StringAttribute{
									Description: "The time of day when the scheduled report is generated. Format: HH:MM (24-hour format).",
									Computed:    true,
								},
								"timezone": schema.StringAttribute{
									Description: "The timezone for the scheduled report (e.g., 'America/New_York', 'UTC').",
									Computed:    true,
								},
								"weekdays": schema.ListAttribute{
									Description: "List of weekdays when the scheduled report is generated (e.g., 'Monday', 'Tuesday'). Only applicable for weekly schedules.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

// reportsDataSourceModel maps the data source schema data.
type reportsDataSourceModel struct {
	ReportID   types.String  `tfsdk:"report_id"`
	ReportName types.String  `tfsdk:"report_name"`
	Reports    []reportModel `tfsdk:"reports"`
}

// reportModel maps the report schema data.
type reportModel struct {
	ID           types.Number   `tfsdk:"id"`
	ReportName   types.String   `tfsdk:"report_name"`
	ReportType   types.String   `tfsdk:"report_type"`
	Asq          types.String   `tfsdk:"asq"`
	CreationTime types.String   `tfsdk:"creation_time"`
	IsScheduled  types.Bool     `tfsdk:"is_scheduled"`
	Schedule     *scheduleModel `tfsdk:"schedule"`
}

// scheduleModel maps the schedule schema data.
type scheduleModel struct {
	Email            []types.String `tfsdk:"email"`
	RepeatAmount     types.Number   `tfsdk:"repeat_amount"`
	RepeatUnit       types.String   `tfsdk:"repeat_unit"`
	ReportFileFormat types.String   `tfsdk:"report_file_format"`
	TimeOfDay        types.String   `tfsdk:"time_of_day"`
	Timezone         types.String   `tfsdk:"timezone"`
	Weekdays         []types.String `tfsdk:"weekdays"`
}

// Read refreshes the Terraform state with the latest data.
func (d *reportsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config reportsDataSourceModel

	tflog.Info(ctx, "Reading reports data source")

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var reports []reportModel

	// Fetch report by ID if specified, otherwise fetch all reports
	if !config.ReportID.IsNull() {
		reports = d.fetchReportByID(ctx, config.ReportID.ValueString(), resp)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		reports = d.fetchAndFilterReports(ctx, config.ReportName, resp)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Info(ctx, "Setting reports state", map[string]any{"report_count": len(reports)})

	// Save data into Terraform state
	config.Reports = reports
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fetchReportByID fetches a single report by ID.
func (d *reportsDataSource) fetchReportByID(ctx context.Context, reportID string, resp *datasource.ReadResponse) []reportModel {
	tflog.Debug(ctx, "Fetching report by ID", map[string]any{"report_id": reportID})

	report, err := d.client.GetReportByID(ctx, reportID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Report",
			fmt.Sprintf("Failed to fetch report with ID '%s': %s", reportID, err.Error()),
		)
		return nil
	}

	// Check for nil report response
	if report == nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Report",
			fmt.Sprintf("Report with ID '%s' was not found or returned empty response.", reportID),
		)
		return nil
	}

	tflog.Debug(ctx, "Successfully fetched report by ID", map[string]any{
		"report_id":   reportID,
		"report_name": report.ReportName,
	})

	// Map response body to model
	reportState := mapReportToModel(report)
	return []reportModel{reportState}
}

// fetchAndFilterReports fetches all reports and optionally filters by name.
func (d *reportsDataSource) fetchAndFilterReports(ctx context.Context, reportName types.String, resp *datasource.ReadResponse) []reportModel {
	tflog.Debug(ctx, "Fetching all reports")

	allReports, err := d.client.GetReports(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Armis Reports",
			fmt.Sprintf("Failed to fetch reports from Armis API: %s", err.Error()),
		)
		return nil
	}

	tflog.Debug(ctx, "Successfully fetched reports", map[string]any{"total_count": len(allReports)})

	// Filter by name if specified
	filterByName := !reportName.IsNull()
	nameFilter := reportName.ValueString()

	if filterByName {
		tflog.Debug(ctx, "Filtering reports by name", map[string]any{"report_name": nameFilter})
	}

	var reports []reportModel
	for i, report := range allReports {
		// Validate report data before processing
		if report.ReportName == "" {
			tflog.Warn(ctx, "Skipping report with empty name", map[string]any{"index": i, "report_id": report.ID})
			continue
		}

		// Filter by name if specified
		if filterByName && report.ReportName != nameFilter {
			continue
		}

		reportState := mapReportToModel(&report)
		reports = append(reports, reportState)
	}

	if filterByName {
		tflog.Debug(ctx, "Filtered reports by name", map[string]any{
			"report_name":    nameFilter,
			"matched_count":  len(reports),
			"original_count": len(allReports),
		})
	}

	return reports
}

// mapReportToModel converts an armis.Report to a reportModel.
func mapReportToModel(report *armis.Report) reportModel {
	var schedulePtr *scheduleModel

	// Only map schedule if the report is scheduled and has schedule data
	// This prevents nil pointer dereference when accessing report.Schedule fields
	if report.IsScheduled {
		// Map email addresses
		var emails []types.String
		for _, email := range report.Schedule.Email {
			emails = append(emails, types.StringValue(email))
		}

		// Map weekdays
		var weekdays []types.String
		for _, weekday := range report.Schedule.Weekdays {
			weekdays = append(weekdays, types.StringValue(weekday))
		}

		schedulePtr = &scheduleModel{
			Email:            emails,
			RepeatAmount:     types.NumberValue(big.NewFloat(report.Schedule.RepeatAmount)),
			RepeatUnit:       types.StringValue(report.Schedule.RepeatUnit),
			ReportFileFormat: types.StringValue(report.Schedule.ReportFileFormat),
			TimeOfDay:        types.StringValue(report.Schedule.TimeOfDay),
			Timezone:         types.StringValue(report.Schedule.Timezone),
			Weekdays:         weekdays,
		}
	}

	// Handle reportType - use null if empty since some Armis reports don't have a type
	reportType := types.StringValue(report.ReportType)
	if report.ReportType == "" {
		reportType = types.StringNull()
	}

	return reportModel{
		ID:           types.NumberValue(big.NewFloat(float64(report.ID))),
		ReportName:   types.StringValue(report.ReportName),
		ReportType:   reportType,
		Asq:          types.StringValue(report.Asq),
		CreationTime: types.StringValue(report.CreationTime),
		IsScheduled:  types.BoolValue(report.IsScheduled),
		Schedule:     schedulePtr,
	}
}
