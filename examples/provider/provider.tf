provider "bastion" {
  host      = "bastion.example.com"
  port      = 22
  username  = "bastionadmin"
  use_agent = true # use ssh agent for authentication
}
