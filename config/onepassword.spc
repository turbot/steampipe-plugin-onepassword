connection "onepassword" {
  plugin = "onepassword"

  # Token is required for requests. Required.
  # See instructions at https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token
  # This can also be set via the `OP_CONNECT_TOKEN` environment variable.
  # token = "eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"

  # The host URL set to default http://localhost:8080. Optional.
  # This can also be set via the `OP_CONNECT_HOST` environment variable.
  # url = "http://localhost:8080"

  # A 1Password Events Reporting bearer token, for onepassword_signin_attempt, onepassword_audit_event
  # and onepassword_item_usage. Optional.
  # See https://support.1password.com/events-reporting/
  # events_token = "eyJhbGciOiJFUzI1NiIsImtpZCI6InR5N2tlZmE3dzZ0YnNnZ3VkcmJ5cHdmc3NlIiwidHlwIjoiSldUIn0.eyJldmVudHMiOlsiYXVkaXRldmVudHMiLCJpdGVtdXNhZ2VzIiwic2lnbmluYXR0ZW1wdHMiXX0"

  # The Events API URL for your account's region. Optional, defaults to https://events.1password.com.
  # Use https://events.ent.1password.com, https://events.1password.ca or https://events.1password.eu for those accounts.
  # events_url = "https://events.1password.com"
}
