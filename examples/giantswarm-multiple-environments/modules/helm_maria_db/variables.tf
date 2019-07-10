variable "name" {
    description = "Name of the resource."
}

variable "version" {
    description = "Version of the db selected."
    default = "6.0.0"
}

variable "username" {
    description = "Username of the db."
}

variable "password" {
    description = "Password for the username previosly created."
}