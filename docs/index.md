---
organization: Turbot
category: ["security"]
icon_url: "/images/plugins/turbot/onepassword.svg"
brand_color: "#1A8CFF"
display_name: "1Password"
short_name: "onepassword"
description: "Steampipe plugin to query vaults, items, files and more from 1Password."
og_description: "Query 1Password with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/onepassword-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# 1Password + Steampipe

[1Password](https://1password.com/) is a password manager, digital vault, form filler and secure digital wallet.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

List your 1Password items:

```sql
select
  id,
  title,
  vault_id,
  favorite,
  category,
  created_at
from
  onepassword_item;
```

```
+----------------------------+---------------------------------------------+----------------------------+----------+------------------+---------------------------+
| id                         | title                                       | vault_id                   | favorite | category         | created_at                |
+----------------------------+---------------------------------------------+----------------------------+----------+------------------+---------------------------+
| kvmaoszyhzbvze6g5tvr6qg2a  | steampipe-test Access Token: steampipe-test | wygy6zfmgzdlzckgruraltkma | false     | API_CREDENTIAL   | 2022-10-11T20:36:34+05:30 |
| ys6wwudn2jchffycnvaruun7y  | Secure Note                                 | wygy6zfmgzdlzckgturaltkma | false     | SECURE_NOTE      | 2023-04-24T15:15:56+05:30 |
| jskefwj3k5nefswdwfopxv4ca  | API Credential                              | wygy6zfmgzdlzckgturaltkma | false     | API_CREDENTIAL   | 2023-04-24T14:51:08+05:30 |
+----------------------------+---------------------------------------------+----------------------------+----------+------------------+---------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/onepassword/tables)**

## Quick start

### Install

Download and install the latest 1Password plugin:

```sh
steampipe plugin install onepassword
```

### Credentials

| Item        | Description |
| ----------- | ----------- |
| Credentials | Either a [service account token](https://developer.1password.com/docs/service-accounts/get-started/) (`ops_...`), or a [Connect server](https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token) access token and URL. |
| Permissions | Both credentials see only the vaults they have been granted access to. `onepassword_item_file` needs a Connect server. |
| Radius      | Each connection represents a single 1Password account (service account) or Connect server. |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/onepassword.spc`), where `service_account_token` takes precedence over `token`.<br />2. Credentials specified in environment variables, where `OP_SERVICE_ACCOUNT_TOKEN` takes precedence over `OP_CONNECT_TOKEN` and `OP_CONNECT_HOST`. |

### Configuration

Installing the latest 1Password plugin will create a config file (`~/.steampipe/config/onepassword.spc`) with a single connection named `onepassword`:

Configure your account details in `~/.steampipe/config/onepassword.spc`:

```hcl
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
```

Alternatively, you can also use the standard 1Password environment variables to obtain credentials **only if `service_account_token` and `token` are not specified** in the connection:

```sh
export OP_CONNECT_TOKEN=eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI
export OP_CONNECT_HOST=http://localhost:8080
```

Or, for a service account:

```sh
export OP_SERVICE_ACCOUNT_TOKEN=ops_eyJzaWduSW5BZGRyZXNzIjoibXkuMXBhc3N3b3JkLmNvbSIsInVzZXJBdXRoIjp7Im1ldGhvZCI6IlNSUGctNDA5NiIsImFsZyI6IlBCRVMyZy1IUzI1NiIsIml0ZXJhdGlvbnMiOjY1MDAwMH19
```
