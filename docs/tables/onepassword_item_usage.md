---
title: "Steampipe Table: onepassword_item_usage - Query 1Password Item Usage using SQL"
description: "Allows users to query how items in shared 1Password vaults are used through the 1Password Events API, including fills, reveals, copies, shares and exports."
---

# Table: onepassword_item_usage - Query 1Password Item Usage using SQL

The 1Password Events API reports when items in shared vaults are used: filled, revealed, copied, shared, exported or fetched through an integration. Each usage records the item and vault, the user, the 1Password app and platform used, and the approximate location.

## Table Usage Guide

The `onepassword_item_usage` table provides insights into how shared credentials are used. As a security analyst, use it to see who is using which shared items, spot reveals and exports of sensitive items, and find out how widely each shared item is used.

**Important Notes**
- This table requires `events_token` in the connection configuration: a bearer token from a 1Password [Events Reporting integration](https://support.1password.com/events-reporting/) with access to item usage. Set `events_url` if your account is not on `1password.com`.
- The table returns the last 24 hours of usage unless you constrain `timestamp` in the `where` clause with `>`, `>=`, `<` or `<=`.
- Usage is reported for items in shared vaults only, not personal or private vaults.

## Examples

### Basic info
Review the item usage of the last 24 hours, with the user and the action taken.

```sql+postgres
select
  uuid,
  timestamp,
  action,
  item_uuid,
  vault_uuid,
  user_email
from
  onepassword_item_usage;
```

```sql+sqlite
select
  uuid,
  timestamp,
  action,
  item_uuid,
  vault_uuid,
  user_email
from
  onepassword_item_usage;
```

### List items revealed or exported in the last 7 days
Find the shared items whose contents were shown or exported in the past week.

```sql+postgres
select
  timestamp,
  action,
  item_uuid,
  vault_uuid,
  user_email
from
  onepassword_item_usage
where
  timestamp >= now() - interval '7 days'
  and action in ('reveal', 'export');
```

```sql+sqlite
select
  timestamp,
  action,
  item_uuid,
  vault_uuid,
  user_email
from
  onepassword_item_usage
where
  timestamp >= datetime('now', '-7 days')
  and action in ('reveal', 'export');
```

### List the most used items in the last 24 hours
Identify the shared credentials your team relies on most.

```sql+postgres
select
  item_uuid,
  count(*) as uses
from
  onepassword_item_usage
group by
  item_uuid
order by
  uses desc
limit 10;
```

```sql+sqlite
select
  item_uuid,
  count(*) as uses
from
  onepassword_item_usage
group by
  item_uuid
order by
  uses desc
limit 10;
```

### Count the distinct users of each item in the last 30 days
See how widely each shared item is used, to decide whether it should be replaced with individual credentials.

```sql+postgres
select
  item_uuid,
  vault_uuid,
  count(distinct user_uuid) as users
from
  onepassword_item_usage
where
  timestamp >= now() - interval '30 days'
group by
  item_uuid,
  vault_uuid
order by
  users desc;
```

```sql+sqlite
select
  item_uuid,
  vault_uuid,
  count(distinct user_uuid) as users
from
  onepassword_item_usage
where
  timestamp >= datetime('now', '-30 days')
group by
  item_uuid,
  vault_uuid
order by
  users desc;
```

### Show item titles and vault names for usage in the last 24 hours
Join item usage with the item and vault tables to see which items were used by name. The connection also needs credentials for vaults and items.

```sql+postgres
select
  u.timestamp,
  u.action,
  u.user_email,
  i.title as item_title,
  v.name as vault_name
from
  onepassword_item_usage as u
  join onepassword_item as i on i.id = u.item_uuid
  join onepassword_vault as v on v.id = u.vault_uuid
order by
  u.timestamp desc;
```

```sql+sqlite
select
  u.timestamp,
  u.action,
  u.user_email,
  i.title as item_title,
  v.name as vault_name
from
  onepassword_item_usage as u
  join onepassword_item as i on i.id = u.item_uuid
  join onepassword_vault as v on v.id = u.vault_uuid
order by
  u.timestamp desc;
```
