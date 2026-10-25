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

| Item        | Description                                                                                                                                                                                                                                                                                                      |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | 1Password requires an [Access Token](https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token) or [Access Token](https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token) and a URL for all requests. The `onepassword_user` table instead needs the client ID and secret of a [Users API](https://www.1password.dev/users-api/get-started) OAuth application and your account ID. |
| Permissions | The permission scope of access tokens is limited to the vaults that the admin provides access to. The Users API OAuth application needs the `users.view` scope. |
| Radius      | Each connection represents a single 1Password Installation.                                                                                                                                                                                                                                                      |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/onepassword.spc`)<br />2. Credentials specified in environment variables, e.g., `OP_CONNECT_TOKEN`, `OP_CONNECT_HOST`.                                                                                                            |

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

  # The client ID and client secret of a 1Password OAuth application with the `users.view` scope,
  # and the 26-character ID of your account, for the onepassword_user table. Optional; all three
  # must be set together.
  # See https://www.1password.dev/users-api/get-started
  # users_client_id     = "H3WVOKGN4B7JXQYCQ6DJFAW5ZU"
  # users_client_secret = "Q2s4ZmVnaGprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OWFi"
  # users_account_id    = "WX52YXL2MVEPBK43P3R4XMWKXI"

  # The Users API URL for your account's region. Optional, defaults to https://api.1password.com.
  # Use https://api.ent.1password.com, https://api.1password.ca or https://api.1password.eu for those accounts.
  # users_api_url = "https://api.1password.com"
}
```

Alternatively, you can also use the standard 1Password environment variables to obtain credentials **only if `token` and `url` are not specified** in the connection:

```sh
export OP_CONNECT_TOKEN=eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI
export OP_CONNECT_HOST=http://localhost:8080
```


