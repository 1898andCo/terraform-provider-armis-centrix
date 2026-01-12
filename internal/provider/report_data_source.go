// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		Description: "Retrieves Armis report information. Optionally filter by report_id or report_name. If no filter is provided, all reports are returned.",
		Attributes: map[string]schema.Attribute{
			"report_id": schema.StringAttribute{
				Description: "An optional report ID used to filter the retrieved report information. If specified, only the report matching this ID will be returned.",
				Optional:    true,
			},
			"report_name": schema.StringAttribute{
				Description: "An optional report name used to filter the retrieved report information. If specified, only reports matching this name will be returned.",
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
							Description: "The schedule configuration for the report.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"email": schema.ListAttribute{
									Description: "A list of email addresses to receive the scheduled report.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"repeat_amount": schema.NumberAttribute{
									Description: "The repeat interval amount for the scheduled report.",
									Computed:    true,
								},
								"repeat_unit": schema.StringAttribute{
									Description: "The repeat interval unit for the scheduled report (e.g., 'days', 'weeks').",
									Computed:    true,
								},
								"report_file_format": schema.StringAttribute{
									Description: "The file format of the scheduled report (e.g., 'pdf', 'csv').",
									Computed:    true,
								},
								"time_of_day": schema.StringAttribute{
									Description: "The time of day when the scheduled report is generated.",
									Computed:    true,
								},
								"timezone": schema.StringAttribute{
									Description: "The timezone for the scheduled report.",
									Computed:    true,
								},
								"weekdays": schema.ListAttribute{
									Description: "A list of weekdays when the scheduled report is generated.",
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

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var reports []reportModel

	if !config.ReportID.IsNull() {
		// Fetch a specific report by ID
		report, err := d.client.GetReportByID(ctx, config.ReportID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Report",
				err.Error(),
			)
			return
		}

		// Map response body to model
		reportState := mapReportToModel(report)
		reports = append(reports, reportState)
	} else {
		// Fetch all reports
		allReports, err := d.client.GetReports(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Armis Reports",
				err.Error(),
			)
			return
		}

		// Filter by name if specified
		for _, report := range allReports {
			if !config.ReportName.IsNull() && report.ReportName != config.ReportName.ValueString() {
				continue
			}
			reportState := mapReportToModel(&report)
			reports = append(reports, reportState)
		}
	}

	// Save data into Terraform state
	config.Reports = reports
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// mapReportToModel converts an armis.Report to a reportModel.
func mapReportToModel(report *armis.Report) reportModel {
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

	return reportModel{
		ID:           types.NumberValue(big.NewFloat(float64(report.ID))),
		ReportName:   types.StringValue(report.ReportName),
		ReportType:   types.StringValue(report.ReportType),
		Asq:          types.StringValue(report.Asq),
		CreationTime: types.StringValue(report.CreationTime),
		IsScheduled:  types.BoolValue(report.IsScheduled),
		Schedule: &scheduleModel{
			Email:            emails,
			RepeatAmount:     types.NumberValue(big.NewFloat(report.Schedule.RepeatAmount)),
			RepeatUnit:       types.StringValue(report.Schedule.RepeatUnit),
			ReportFileFormat: types.StringValue(report.Schedule.ReportFileFormat),
			TimeOfDay:        types.StringValue(report.Schedule.TimeOfDay),
			Timezone:         types.StringValue(report.Schedule.Timezone),
			Weekdays:         weekdays,
		},
	}
}
