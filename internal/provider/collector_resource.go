// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	armis "github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &collectorResource{}
	_ resource.ResourceWithConfigure = &collectorResource{}
)

type collectorResource struct {
	client *armis.Client
}

func CollectorResource() resource.Resource {
	return &collectorResource{}
}

// Configure adds the provider configured client to the resource.
func (r *collectorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *collectorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collector"
}

// Schema defines the schema for the user resource.
func (r *collectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
		Provides an Armis collector

		The resource provisions a collector to perform asset discovery, monitoring, and threat mitigation across various environments.
		For more information, see: https://media.armis.com/sb-armis-collectors-en.pdf
		`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The full name of the collector.",
			},
			"deployment_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of deployment. Valid options include 'VHDX', 'AMI', 'QCOW2', 'OVA', and 'VHD'",
				Validators: []validator.String{
					stringvalidator.OneOf("VHDX", "AMI", "QCOW2", "OVA", "VHD"),
				},
			},
			"license_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The license key associated with the collector.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password associated with the collector.",
			},
			"user": schema.StringAttribute{
				Computed:    true,
				Description: "The unique username of the user.",
			},
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "A unique identifier for the user resource.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

// collectorResourceModel maps the resource schema data.
type collectorResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	DeploymentType types.String `tfsdk:"deployment_type"`
	LicenseKey     types.String `tfsdk:"license_key"`
	Password       types.String `tfsdk:"password"`
	User           types.String `tfsdk:"user"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *collectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan collectorResourceModel
	tflog.Info(ctx, "Creating collector")

	// Parse the plan from Terraform
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	collector := armis.CreateCollectorSettings{
		Name:           plan.Name.ValueString(),
		DeploymentType: plan.DeploymentType.ValueString(),
	}

	// Create the collector via the client
	newCollector, err := r.client.CreateCollector(ctx, collector)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating collector",
			"Could not create collector, unexpected error: "+err.Error(),
		)
		return
	}

	// Map the response to Terraform state
	plan.ID = types.StringValue(strconv.Itoa(newCollector.CollectorID))
	plan.LicenseKey = types.StringValue(newCollector.LicenseKey)
	plan.Password = types.StringValue(newCollector.Password)
	plan.User = types.StringValue(newCollector.User)

	// Save the state
	tflog.Info(ctx, "Setting state for collector")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read collector resource information.
func (r *collectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state collectorResourceModel
	tflog.Info(ctx, "Retrieving current state")
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed collector value from Armis
	collector, err := r.client.GetCollectorByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Armis Collector",
			"Could not read Armis collector ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if collector == nil {
		resp.Diagnostics.AddError(
			"Error Reading Armis Collector",
			"Could not read Armis collector ID "+state.ID.ValueString()+": Collector not found",
		)
		return
	}

	//	Overwrite collector with refreshed state
	state.Name = types.StringValue(collector.Name)
	state.ID = types.StringValue(strconv.Itoa(collector.CollectorNumber))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *collectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan collectorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state to get the collector ID
	var state collectorResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that the collector ID is available
	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error Updating Collector",
			"The collector ID is missing from the state. This is required to update the collector.",
		)
		return
	}

	// Map the Terraform model to the API's collector struct
	collector := armis.UpdateCollectorSettings{
		Name:           plan.Name.ValueString(),
		DeploymentType: plan.DeploymentType.ValueString(),
	}

	// Update existing collector
	// and then fetch the updated collector from the API.
	_, err := r.client.UpdateCollector(ctx, state.ID.ValueString(), collector)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Armis collector",
			"Could not update collector, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated collector from the API to ensure the state is fully
	// populated.
	updatedCollector, err := r.client.GetCollectorByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Armis Collector",
			"Could not read Armis collector ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated collector options and timestamp
	// Map the response to Terraform state
	plan.ID = types.StringValue(strconv.Itoa(updatedCollector.CollectorNumber))
	plan.Name = types.StringValue(updatedCollector.Name)
	plan.DeploymentType = types.StringValue(plan.DeploymentType.ValueString())
	plan.LicenseKey = types.StringValue(state.LicenseKey.ValueString())
	plan.Password = types.StringValue(state.Password.ValueString())
	plan.User = types.StringValue(state.User.ValueString())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *collectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state collectorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	success, err := r.client.DeleteCollector(ctx, state.ID.ValueString())
	if err != nil || !success {
		resp.Diagnostics.AddError(
			"Error Deleting Armis collector",
			"Could not delete collector, unexpected error: "+err.Error(),
		)
		return
	}
}
