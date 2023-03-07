---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oktawave_oks_node Resource - terraform-provider-oktawave"
subcategory: ""
description: |-
  
---

# oktawave_oks_node (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (String) This is the same as cluster name
- `subregion_id` (Number) ID from subregions resource
- `type_id` (Number) Value from dictionary #12

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `cpu_number` (Number)
- `creation_date` (String)
- `id` (String) The ID of this resource.
- `ip_address` (String)
- `name` (String)
- `ram_mb` (Number)
- `status_id` (Number) Value from dictionary #27
- `total_disks_capacity` (Number)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)

