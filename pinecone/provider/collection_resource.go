// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/skyscrapr/pinecone-sdk-go/pinecone"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}

func NewCollectionResource() resource.Resource {
	return &CollectionResource{}
}

// CollectionResource defines the resource implementation.
type CollectionResource struct {
	client *pinecone.Client
}

// CollectionResourceModel describes the resource data model.
type CollectionResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Source types.String `tfsdk:"source"`
	Size   types.Int64  `tfsdk:"size"`
	Status types.String `tfsdk:"status"`
}

func (r *CollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (r *CollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Collection resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Collection identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the collection.",
				Required:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "The name of the source index to be used as the source for the collection.",
				Required:            true,
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "The size of the collection in bytes.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the collection.",
				Computed:            true,
			},
		},
	}
}

func (r *CollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pinecone.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pinecone.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *CollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CollectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := pinecone.CreateCollectionParams{
		Name:   data.Name.ValueString(),
		Source: data.Source.ValueString(),
	}

	err := r.client.Collections().CreateCollection(&payload)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create collection", err.Error())
		return
	}

	// Wait for collection to be ready
	createTimeout := 1 * time.Hour
	err = retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		collection, err := r.client.Collections().DescribeCollection(data.Name.ValueString())

		readCollectionData(collection, &data)
		// Save current status to state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

		if err != nil {
			return retry.NonRetryableError(err)
		}
		if collection.Status != "Ready" {
			return retry.RetryableError(fmt.Errorf("collection not ready. State: %s", collection.Status))
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for collection to become ready.", err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collection, err := r.client.Collections().DescribeCollection(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to describe collection", err.Error())
		return
	}

	readCollectionData(collection, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Collections currently do not support updates
}

func (r *CollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Collections().DeleteCollection(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete collection", err.Error())
		return
	}
	// Wait for collection to be deleted
	deleteTimeout := 1 * time.Hour
	err = retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		collection, err := r.client.Collections().DescribeCollection(data.Id.ValueString())
		tflog.Info(ctx, fmt.Sprintf("Deleting Collection. Status: '%s'", collection.Status))

		if err != nil {
			if pineconeErr, ok := err.(*pinecone.HTTPError); ok && pineconeErr.StatusCode == 404 {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		return retry.RetryableError(fmt.Errorf("collection not deleted. State: %s", collection.Status))
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for collection to be deleted.", err.Error())
		return
	}
}

func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func readCollectionData(collection *pinecone.Collection, model *CollectionResourceModel) {
	model.Id = types.StringValue(collection.Name)
	model.Name = types.StringValue(collection.Name)
	model.Source = types.StringValue(model.Source.ValueString())
	model.Size = types.Int64Value(int64(collection.Size))
	model.Status = types.StringValue(collection.Status)
}