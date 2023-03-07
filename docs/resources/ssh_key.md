---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oktawave_ssh_key Resource - terraform-provider-oktawave"
subcategory: ""
description: |-
  Ssh keys - used for connections security.
---

# oktawave_ssh_key (Resource)

Ssh keys - used for connections security.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of ssh key.
- `value` (String) Public ssh key

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `creation_date` (String) Creation date.
- `id` (String) The ID of this resource.
- `owner_user_id` (Number) Id of user who created this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)

