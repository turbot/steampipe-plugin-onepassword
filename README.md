![image](https://hub.steampipe.io/images/plugins/turbot/onepassword-social-graphic.png)

# 1Password Plugin for Steampipe

Use SQL to query vaults, items, files and more from 1Password.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/onepassword)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/onepassword/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
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

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs/steampipe_sqlite/overview) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/overview) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

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

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [1Password Plugin](https://github.com/turbot/steampipe-plugin-onepassword/labels/help%20wanted)
