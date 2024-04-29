---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pinecone_projects Data Source - terraform-provider-pinecone"
subcategory: ""
description: |-
  Projects data source
---

# pinecone_projects (Data Source)

Projects data source



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) Projects identifier
- `projects` (Attributes List) List of the indexes in your project (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Read-Only:

- `id` (String) Project identifier
- `name` (String) Index name