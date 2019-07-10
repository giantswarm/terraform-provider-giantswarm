resource "cloudflare_record" "root" {
  domain = "${var.domain}"
  name   = "@"
  value  = "${var.dns_value}"
  type   = "${var.record_type}"
  ttl    = "${var.record_ttl}"
}

resource "cloudflare_record" "www" {
  domain = "${var.domain}"
  name   = "${var.dns_record}"
  value  = "${var.dns_value}"
  type   = "${var.record_type}"
  ttl    = "${var.record_ttl}"
}
