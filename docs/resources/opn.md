---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oktawave_opn Resource - terraform-provider-oktawave"
subcategory: ""
description: |-
  Oktawave Private Network(OPN) is a conterpart of typical VLAN network.
---

# oktawave_opn (Resource)

Oktawave Private Network(OPN) is a conterpart of typical VLAN network.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of OPN

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `creation_date` (String) Creation date
- `creation_user_id` (Number) Id of user who created this resource.
- `id` (String) The ID of this resource.
- `instance_ids` (Set of Number) List of instance ids in this OPN.
- `last_change_date` (String) Last change date

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)


