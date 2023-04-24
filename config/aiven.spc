connection "aiven" {
  plugin = "aiven"

  # You can connect to Aiven using one of the options below:

  # Using API Key authentication
  # `api_key` (required) - Create an authentication token in the Aiven Console for use with the Aiven CLI or API.
  # To create an authentication token, refer to https://docs.aiven.io/docs/platform/howto/create_authentication_token
  # Can also be set with the AIVEN_TOKEN environment variable.
  # api_key = "oGAxUvrjAdL3QBhWnaJI67Pc9P0rPDzDfhykzVfBYPlmvVH8WdJMKaeVKzcrl4CnyXpjiaKJCCNT+OkbpxfWdDNqwZPngS"

  # Using User authentication (without 2FA)
  # email = "test@turbot.com"
  # password = "test@123"

  # If no credentials are specified, the plugin will use Aiven CLI authentication.
}
