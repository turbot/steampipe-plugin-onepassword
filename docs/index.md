---
organization: Turbot
category: ["security"]
icon_url: "/images/plugins/turbot/1password.svg"
brand_color: "#1A8CFF"
display_name: "1Password"
short_name: "1password"
description: "Steampipe plugin to query vaults, items, files and more from 1Password."
og_description: "Query 1Password with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/1password-social-graphic.png"
---

# 1Password + Steampipe

[1Password](https://1password.com/) is a password manager, digital vault, form filler and secure digital wallet.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

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
| kvmaoszyhzbvze6g5tvr6qg2a | steampipe-test Access Token: steampipe-test | wygy6zfmgzdlzckgruraltkma | false    | API_CREDENTIAL   | 2022-10-11T20:36:34+05:30 |
| ys6wwudn2jchffycnvaruun7y | Secure Note                                 | wygy6zfmgzdlzckgturaltkma | false    | SECURE_NOTE      | 2023-04-24T15:15:56+05:30 |
| jskefwj3k5nefswdwfopxv4ca | API Credential                              | wygy6zfmgzdlzckgturaltkma | false    | API_CREDENTIAL   | 2023-04-24T14:51:08+05:30 |
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/1password/tables)**

## Quick start

### Install

Download and install the latest 1Password plugin:

```sh
steampipe plugin install onepassword
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                                      |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | 1Password requires an [Access Token](https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token) or [Access Token](https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token) and a URL for all requests. |
| Permissions | The permission scope of access tokens is limited to the vaults that the admin provides access to.                                                                                                                                                                                                                |
| Radius      | Each connection represents a single 1Password Installation.                                                                                                                                                                                                                                                      |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/onepassword.spc`)<br />2. Credentials specified in environment variables, e.g., `OP_CONNECT_TOKEN`, `OP_CONNECT_HOST`.                                                                                                            |

### Configuration

Installing the latest 1Password plugin will create a config file (`~/.steampipe/config/onepassword.spc`) with a single connection named `onepassword`:

Configure your account details in `~/.steampipe/config/onepassword.spc`:

```hcl
connection "onepassword" {
  plugin = "onepassword"

  # `token` (required) - To create an access token, refer to https://developer.1password.com/docs/connect/manage-secrets-automation#issue-revoke-or-rename-an-access-token
  # Can also be set with the OP_CONNECT_TOKEN environment variable.
  # token = "eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"

  # `url` (optional) - The host URL. Set to default http://localhost:8080
  # Can also be set with the OP_CONNECT_HOST environment variable.
  # url = "http://localhost:8080"
}
```

## Configuring 1Password Credentials

### Access Token Credentials

You may specify the Access Token to authenticate:

- `token`: Specify the access token.

```hcl
connection "onepassword" {
  plugin = "onepassword"
  token  = "eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"
}
```

### Access Token and Host URL Credentials

You may specify the Access Token and Host URL to authenticate:

- `token`: Specify the access token.
- `url` : Specify the host url.

```hcl
connection "onepassword" {
  plugin = "onepassword"
  token  = "eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"
  url = "http://localhost:8080"
}
```

### Credentials from Environment Variables

The 1Password plugin will use the 1Password environment variable to obtain credentials **only if the `token` or `url` is not specified** in the connection:

```sh
export OP_CONNECT_TOKEN="eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"
export OP_CONNECT_HOST="http://localhost:8080"
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-1password
- Community: [Slack Channel](https://steampipe.io/community/join)
