---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oktawave_opn Data Source - terraform-provider-oktawave"
subcategory: ""
description: |-
  
---

# oktawave_opn (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `creation_date` (String)
- `creation_user_id` (Number)
- `id` (Number) The ID of this resource.
- `last_change_date` (String)
- `name` (String)
- `private_ips` (Set of Object) (see [below for nested schema](#nestedatt--private_ips))

<a id="nestedatt--private_ips"></a>
### Nested Schema for `private_ips`

Read-Only:

- `address` (String)
- `address_v6` (String)
- `creation_date` (String)
- `instance_id` (Number)
- `interface_id` (Number)
- `mac_address` (String)

