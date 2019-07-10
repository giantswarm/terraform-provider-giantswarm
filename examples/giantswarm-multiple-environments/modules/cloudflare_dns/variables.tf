variable "domain" {
    description = "The domain of our web service."
}

variable "dns_record" {
    description = "Value of the dns_record to add in the domain zone"
}

variable "record_type" {
    description = "DNS record type. Ex A, CNAME, TXT, AAAA"
}

variable "dns_value" {
    description = "DNS record value. Ex a public IP or CNAME"
}

variable "record_ttl" {
    description = "DNS record TTL"
    default = 3600
}
