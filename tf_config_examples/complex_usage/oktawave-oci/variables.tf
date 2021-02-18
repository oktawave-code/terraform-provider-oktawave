variable "instance_name" {
  description = "Instance name"
}

variable "authorization_method_id" {
  description = "Ssh key (1398) or password (1399) ID"
  default     = 1399
}

variable "ssh_keys_ids" {
  description = "Ssh key id if authorization_method_id set to 1399"
  default     = null
}

variable "disk_class" {
  description = "Disk tier ID"
}

variable "init_disk_size" {
  description = "Disk size"
  default     = 5
}

variable "ip_address_ids" {
  description = "Custom ip address ID"
  default     = null
}

variable "subregion_id" {
  description = "Subregion ID"
}

variable "type_id" {
  description = "Instance type ID"
}

variable "template_id" {
  description = "Template ID"
}

variable "instances_count" {
  description = "Number of instances"
  default     = 1
}

variable "isfreemium" {
  description = "Freemium OCI"
  default     = false
}

variable "opn_ids" {
  description = "list of OPN ids"
  default     = null
}

variable "ovs_disk_name" {
  description = "Attached disk name"
}

variable "ovs_space_capacity" {
  description = "Disk size"
}

variable "ovs_tier_id" {
  description = "Disk tier"
}

variable "init_script" {
  description = "Puppet base64 script running on first boot"
}

variable "init_script_file" {
  description = "Path to Puppet manifest running on first boot"
  default     = "initscript_default.pp"
}

variable "without_public_ip" {
  description = "Create instance without default public interface"
  default     = false
}