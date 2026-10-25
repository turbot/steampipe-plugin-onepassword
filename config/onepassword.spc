connection "onepassword" {
  plugin = "onepassword"

  # Token is required for requests. Required.
  # See instructions at https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token
  # This can also be set via the `OP_CONNECT_TOKEN` environment variable.
  # token = "eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"

  # The host URL set to default http://localhost:8080. Optional.
  # This can also be set via the `OP_CONNECT_HOST` environment variable.
  # url = "http://localhost:8080"

  # A service account token, used instead of a Connect server. Optional.
  # Takes precedence over `token` when both are set.
  # See https://developer.1password.com/docs/service-accounts/get-started/
  # This can also be set via the `OP_SERVICE_ACCOUNT_TOKEN` environment variable.
  # service_account_token = "ops_eyJzaWduSW5BZGRyZXNzIjoibXkuMXBhc3N3b3JkLmNvbSIsInVzZXJBdXRoIjp7Im1ldGhvZCI6IlNSUGctNDA5NiIsImFsZyI6IlBCRVMyZy1IUzI1NiIsIml0ZXJhdGlvbnMiOjY1MDAwMH19"
}
