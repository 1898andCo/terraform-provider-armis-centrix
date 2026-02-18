// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	armis "github.com/1898andCo/armis-sdk-go/v2/armis"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &reportResource{}
	_ resource.ResourceWithConfigure   = &reportResource{}
	_ resource.ResourceWithImportState = &reportResource{}
)

type reportResource struct {
	client *armis.Client
}

func ReportResource() resource.Resource {
	return &reportResource{}
}

// Configure adds the provider configured client to the resource.
func (r *reportResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*armis.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *armis.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *reportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_report"
}

// Schema defines the schema for the report resource.
func (r *reportResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides an Armis report resource.

The resource provisions a report in Armis with support for ASQ queries, scheduling, and export configurations.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "The unique identifier for the report.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"report_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the report.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"asq": schema.StringAttribute{
				Required:    true,
				Description: "The Armis Standard Query (ASQ) for the report. Example: 'in:devices timeFrame:\"1 Day\"'",
			},
			"email_subject": schema.StringAttribute{
				Optional:    true,
				Description: "The email subject for scheduled report notifications.",
			},
			"creation_time": schema.StringAttribute{
				Computed:      true,
				Description:   "The timestamp when the report was created.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"is_scheduled": schema.BoolAttribute{
				Computed:      true,
				Description:   "Whether the report has a schedule configured.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"schedule": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Schedule configuration for the report.",
				Attributes: map[string]schema.Attribute{
					"email": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "List of email addresses to receive the scheduled report.",
					},
					"repeat_amount": schema.StringAttribute{
						Optional:    true,
						Description: "The interval amount for report scheduling (e.g., '1', '2', '0.5').",
					},
					"repeat_unit": schema.StringAttribute{
						Optional:    true,
						Description: "The interval unit for report scheduling (e.g., 'Days', 'Weeks', 'Months').",
						Validators: []validator.String{
							stringvalidator.OneOf("Days", "Weeks", "Months"),
						},
					},
					"report_file_format": schema.StringAttribute{
						Optional:    true,
						Description: "The file format for the exported report (e.g., 'csv', 'xlsx', 'json').",
						Validators: []validator.String{
							stringvalidator.OneOf("csv", "xlsx", "json"),
						},
					},
					"time_of_day": schema.StringAttribute{
						Optional:    true,
						Description: "The time of day to run the scheduled report (e.g., '15:00').",
					},
					"timezone": schema.StringAttribute{
						Optional:    true,
						Description: "The timezone for the scheduled report (e.g., 'America/New_York', 'UTC').",
					},
					"weekdays": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "List of weekdays to run the report (e.g., 'Monday', 'Tuesday').",
					},
				},
			},
			"export_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Export configuration for the report columns.",
				Attributes: map[string]schema.Attribute{
					"columns": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Column configuration for different report types.",
						Attributes: map[string]schema.Attribute{
							"devices": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "List of device columns to include in the export.",
							},
							"vulnerabilities": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "List of vulnerability columns to include in the export.",
							},
							"activities": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "List of activity columns to include in the export.",
							},
						},
					},
				},
			},
		},
	}
}

// reportResourceModel maps the resource schema data.
type reportResourceModel struct {
	ID                  types.String                    `tfsdk:"id"`
	ReportName          types.String                    `tfsdk:"report_name"`
	ASQ                 types.String                    `tfsdk:"asq"`
	EmailSubject        types.String                    `tfsdk:"email_subject"`
	CreationTime        types.String                    `tfsdk:"creation_time"`
	IsScheduled         types.Bool                      `tfsdk:"is_scheduled"`
	Schedule            *reportScheduleModel            `tfsdk:"schedule"`
	ExportConfiguration *reportExportConfigurationModel `tfsdk:"export_configuration"`
}

type reportScheduleModel struct {
	Email            []types.String `tfsdk:"email"`
	RepeatAmount     types.String   `tfsdk:"repeat_amount"`
	RepeatUnit       types.String   `tfsdk:"repeat_unit"`
	ReportFileFormat types.String   `tfsdk:"report_file_format"`
	TimeOfDay        types.String   `tfsdk:"time_of_day"`
	Timezone         types.String   `tfsdk:"timezone"`
	Weekdays         []types.String `tfsdk:"weekdays"`
}

type reportExportConfigurationModel struct {
	Columns *reportExportColumnsModel `tfsdk:"columns"`
}

type reportExportColumnsModel struct {
	Devices         []types.String `tfsdk:"devices"`
	Vulnerabilities []types.String `tfsdk:"vulnerabilities"`
	Activities      []types.String `tfsdk:"activities"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *reportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan reportResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	report := buildArmisReport(plan)
	tflog.Info(ctx, "Creating report in Armis", map[string]any{"report_name": plan.ReportName.ValueString()})

	newReport, err := r.client.CreateReport(ctx, report)
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error creating report %q", plan.ReportName.ValueString()), err)
		return
	}

	tflog.Info(ctx, "Report created successfully", map[string]any{
		"report_id":   newReport.ID,
		"report_name": newReport.ReportName,
	})

	// Map API response to state
	plan.ID = types.StringValue(strconv.Itoa(newReport.ID))
	plan.CreationTime = types.StringValue(newReport.CreationTime)
	plan.IsScheduled = types.BoolValue(newReport.IsScheduled)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read reads the resource state from the API.
func (r *reportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state reportResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading report from Armis", map[string]any{"report_id": state.ID.ValueString()})

	report, err := r.client.GetReportByID(ctx, state.ID.ValueString())
	if err != nil {
		var ae *armis.APIError
		if errors.As(err, &ae) && ae.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, "Report not found, removing from state", map[string]any{
				"report_id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error reading report %s", state.ID.ValueString()), err)
		return
	}

	if report == nil {
		tflog.Warn(ctx, "Report not found, removing from state", map[string]any{
			"report_id": state.ID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state with API response
	state.ID = types.StringValue(strconv.Itoa(report.ID))
	state.ReportName = types.StringValue(report.ReportName)
	state.ASQ = types.StringValue(report.Asq)
	state.CreationTime = types.StringValue(report.CreationTime)
	state.IsScheduled = types.BoolValue(report.IsScheduled)

	// email_subject is write-only; the API does not return it, so we
	// preserve the value already in state (populated from req.State.Get).

	// export_configuration is not returned by the GetReportByID API,
	// so we preserve the value from the existing Terraform state.

	// Map schedule from API response if present
	if report.IsScheduled {
		state.Schedule = &reportScheduleModel{
			RepeatAmount:     types.StringValue(fmt.Sprintf("%g", report.Schedule.RepeatAmount)),
			RepeatUnit:       types.StringValue(report.Schedule.RepeatUnit),
			ReportFileFormat: types.StringValue(report.Schedule.ReportFileFormat),
			TimeOfDay:        types.StringValue(report.Schedule.TimeOfDay),
			Timezone:         types.StringValue(report.Schedule.Timezone),
		}

		// Map email list
		if len(report.Schedule.Email) > 0 {
			state.Schedule.Email = make([]types.String, len(report.Schedule.Email))
			for i, email := range report.Schedule.Email {
				state.Schedule.Email[i] = types.StringValue(email)
			}
		}

		// Map weekdays list
		if len(report.Schedule.Weekdays) > 0 {
			state.Schedule.Weekdays = make([]types.String, len(report.Schedule.Weekdays))
			for i, weekday := range report.Schedule.Weekdays {
				state.Schedule.Weekdays[i] = types.StringValue(weekday)
			}
		}
	} else {
		state.Schedule = nil
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates an existing report in Armis.
func (r *reportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan reportResourceModel
	var state reportResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reportID := state.ID.ValueString()
	updateReq := buildUpdateReportRequest(plan)

	tflog.Info(ctx, "Updating report in Armis", map[string]any{
		"report_id":   reportID,
		"report_name": plan.ReportName.ValueString(),
	})

	updatedReport, err := r.client.UpdateReport(ctx, reportID, updateReq)
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error updating report %s", reportID), err)
		return
	}

	tflog.Info(ctx, "Report updated successfully", map[string]any{
		"report_id":   updatedReport.ID,
		"report_name": updatedReport.ReportName,
	})

	// Preserve immutable computed fields from state; update mutable computed fields from API response
	plan.ID = state.ID
	plan.CreationTime = state.CreationTime
	plan.IsScheduled = types.BoolValue(updatedReport.IsScheduled)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource.
func (r *reportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state reportResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting report from Armis", map[string]any{"report_id": state.ID.ValueString()})

	success, err := r.client.DeleteReport(ctx, state.ID.ValueString())
	if err != nil {
		appendAPIError(&resp.Diagnostics, fmt.Sprintf("Error deleting report %s", state.ID.ValueString()), err)
		return
	}

	if !success {
		resp.Diagnostics.AddError(
			"Error Deleting Armis Report",
			"Could not delete report: operation returned unsuccessful status",
		)
		return
	}

	tflog.Info(ctx, "Report deleted successfully", map[string]any{"report_id": state.ID.ValueString()})
}

// ImportState imports an existing report into Terraform state.
func (r *reportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// typesStringSliceToStrings converts a slice of types.String to a slice of strings.
func typesStringSliceToStrings(ts []types.String) []string {
	if len(ts) == 0 {
		return nil
	}
	result := make([]string, len(ts))
	for i, t := range ts {
		result[i] = t.ValueString()
	}
	return result
}

// buildScheduleFromPlan converts the Terraform schedule model to an SDK CreateSchedule.
func buildScheduleFromPlan(schedule *reportScheduleModel) armis.CreateSchedule {
	return armis.CreateSchedule{
		RepeatAmount:     schedule.RepeatAmount.ValueString(),
		RepeatUnit:       schedule.RepeatUnit.ValueString(),
		ReportFileFormat: schedule.ReportFileFormat.ValueString(),
		TimeOfDay:        schedule.TimeOfDay.ValueString(),
		Timezone:         schedule.Timezone.ValueString(),
		Email:            typesStringSliceToStrings(schedule.Email),
		Weekdays:         typesStringSliceToStrings(schedule.Weekdays),
	}
}

// buildExportColumnsFromPlan converts the Terraform export columns model to SDK ExportColumns.
func buildExportColumnsFromPlan(columns *reportExportColumnsModel) armis.ExportColumns {
	return armis.ExportColumns{
		Devices:         typesStringSliceToStrings(columns.Devices),
		Vulnerabilities: typesStringSliceToStrings(columns.Vulnerabilities),
		Activities:      typesStringSliceToStrings(columns.Activities),
	}
}

// buildArmisReport converts the Terraform model to an SDK CreateReportRequest.
func buildArmisReport(plan reportResourceModel) armis.CreateReportRequest {
	report := armis.CreateReportRequest{
		ReportName:   plan.ReportName.ValueString(),
		ASQ:          plan.ASQ.ValueString(),
		EmailSubject: plan.EmailSubject.ValueString(),
	}

	if plan.Schedule != nil {
		report.Schedule = buildScheduleFromPlan(plan.Schedule)
	}

	if plan.ExportConfiguration != nil && plan.ExportConfiguration.Columns != nil {
		report.ExportConfiguration = armis.ExportConfiguration{
			Columns: buildExportColumnsFromPlan(plan.ExportConfiguration.Columns),
		}
	}

	return report
}

// buildUpdateReportRequest converts the Terraform model to an SDK UpdateReportRequest.
func buildUpdateReportRequest(plan reportResourceModel) armis.UpdateReportRequest {
	report := armis.UpdateReportRequest{
		ReportName:   plan.ReportName.ValueString(),
		ASQ:          plan.ASQ.ValueString(),
		EmailSubject: plan.EmailSubject.ValueString(),
	}

	if plan.Schedule != nil {
		schedule := buildScheduleFromPlan(plan.Schedule)
		report.Schedule = &schedule
	}

	if plan.ExportConfiguration != nil && plan.ExportConfiguration.Columns != nil {
		exportConfig := armis.ExportConfiguration{
			Columns: buildExportColumnsFromPlan(plan.ExportConfiguration.Columns),
		}
		report.ExportConfiguration = &exportConfig
	}

	return report
}
