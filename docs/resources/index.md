---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pinecone_index Resource - terraform-provider-pinecone"
subcategory: ""
description: |-
  Index resource
---

# pinecone_index (Resource)

Index resource

## Example Usage

```terraform
terraform {
  required_providers {
    pinecone = {
      source = "skyscrapr/pinecone"
    }
  }
}

provider "pinecone" {
  environment = "us-west4-gcp"
  # api_key = set via PINECONE_API_KEY env variable
}

resource "pinecone_index" "test" {
  name      = "tftestindex"
  dimension = 512
  metric    = "cosine"
  pod_type  = "s1.x1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dimension` (Number) The dimensions of the vectors to be inserted in the index
- `name` (String) The name of the index to be created. The maximum length is 45 characters.

### Optional

- `metadata_config` (Attributes) Configuration for the behavior of Pinecone's internal metadata index. By default, all metadata is indexed; when metadata_config is present, only specified metadata fields are indexed. To specify metadata fields to index, provide an array of the following form: [example_metadata_field] (see [below for nested schema](#nestedatt--metadata_config))
- `metric` (String) The distance metric to be used for similarity search. You can use 'euclidean', 'cosine', or 'dotproduct'.
- `pod_type` (String) The type of pod to use. One of s1, p1, or p2 appended with . and one of x1, x2, x4, or x8.
- `pods` (Number) The number of pods for the index to use,including replicas.
- `replicas` (Number) The number of replicas. Replicas duplicate your index. They provide higher availability and throughput.
- `source_collection` (String) The name of the collection to create an index from.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) Index identifier

<a id="nestedatt--metadata_config"></a>
### Nested Schema for `metadata_config`

Optional:

- `indexed` (List of String) The indexed fields.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) Timeout defaults to 5 mins. Accepts a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) Timeout defaults to 5 mins. Accepts a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
