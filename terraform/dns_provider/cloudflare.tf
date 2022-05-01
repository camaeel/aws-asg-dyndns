data "cloudflare_api_token_permission_groups" "all" {

  count = var.provider_name == "cloudflare" ? 1 : 0
}

data "cloudflare_zone" "aws-dyn-dns" {
  name = var.zone_name

  count = var.provider_name == "cloudflare" ? 1 : 0
}

resource "cloudflare_api_token" "aws-dyn-dns" {
  name = "aws-dyn-dns"

  policy {
    permission_groups = [
      data.cloudflare_api_token_permission_groups.all[count.index].permissions["DNS Write"],
    ]
    resources = {
      "com.cloudflare.api.account.zone.${data.cloudflare_zone.aws-dyn-dns[count.index].id}" = "*"
    }
  }

  count = var.provider_name == "cloudflare" ? 1 : 0
}
