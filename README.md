![image](https://hub.steampipe.io/images/plugins/turbot/onepassword-social-graphic.png)

# 1Password Plugin for Steampipe

Use SQL to query vaults, items, files and more from 1Password.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/onepassword)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/onepassword/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-onepassword/issues)

## Quick start

### Install

Download and install the latest 1Password plugin:

```bash
steampipe plugin install onepassword
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/onepassword#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/onepassword#configuration).

Add your configuration details in `~/.steampipe/config/onepassword.spc`:

```hcl
connection "onepassword" {
  plugin = "onepassword"

  # Authentication information
  token  = "eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI"
  url = "http://localhost:8080"
}
```

Or through environment variables:

```sh
export OP_CONNECT_TOKEN=eyJhbGciOiJFUzI1NiIsImtpZCI6InFuN3JwcmZhbnJqZ2V1bWU2eTNidGpjdHN5IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiVEpGVzVZTlRJSkMzSkNXRFgzQ0dWTUpCSDQiLCIxcGFzc3dvcmQuY29tL3Rva2VuIjoib2tnZGZJWHpEaDhWWkNkRHVNRjZNSUplRUlwN3ZrYUQiLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJ2YXVsdGFjY2VzcyJdLCIxcGFzc3dvcmQuY29tL3Z0cyI6W3sidSI6ImZwZDR1dW00bHJicTMycG8ybXR2ZGo0c3hpI
export OP_CONNECT_HOST=http://localhost:8080
```

Run steampipe:

```shell
steampipe query
```

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

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-onepassword.git
cd steampipe-plugin-onepassword
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/onepassword.spc
```

Try it!

```
steampipe query
> .inspect onepassword
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-onepassword/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [1Password Plugin](https://github.com/turbot/steampipe-plugin-onepassword/labels/help%20wanted)
