---
title: "Steampipe Table: onepassword_user - Query 1Password Users using SQL"
description: "Allows users to query the people in a 1Password Business account through the 1Password Users API, including their email address, display name, state and creation time."
---

# Table: onepassword_user - Query 1Password Users using SQL

The 1Password Users API lists the people in a 1Password Business account. Each user has an email address, a display name, a state (active or suspended) and the time the user was created.

## Table Usage Guide

The `onepassword_user` table provides insights into who has a 1Password account. As a security or compliance analyst, use it to reconcile 1Password users against your HR or identity records, find suspended users, and review recently added accounts.

**Important Notes**
- This table requires `users_client_id`, `users_client_secret` and `users_account_id` in the connection configuration: the client ID and secret of an OAuth application created in 1Password with the `users.view` scope, and the 26-character ID of the account. Set `users_api_url` if your account is not on `1password.com`.
- The Users API is not available for accounts that use [automated provisioning](https://support.1password.com/provisioning/).
- `state` is null for users who are neither active nor suspended, such as users pending confirmation.

## Examples

### Basic info
List every user in the account with their state and creation time.

```sql+postgres
select
  id,
  email,
  display_name,
  state,
  create_time
from
  onepassword_user;
```

```sql+sqlite
select
  id,
  email,
  display_name,
  state,
  create_time
from
  onepassword_user;
```

### List suspended users
Find the users who can no longer sign in, for example to confirm offboarding or to clean up accounts.

```sql+postgres
select
  email,
  display_name,
  create_time
from
  onepassword_user
where
  state = 'SUSPENDED';
```

```sql+sqlite
select
  email,
  display_name,
  create_time
from
  onepassword_user
where
  state = 'SUSPENDED';
```

### List users created in the last 30 days
Review the accounts added recently, to check that each was expected.

```sql+postgres
select
  email,
  display_name,
  create_time
from
  onepassword_user
where
  create_time >= now() - interval '30 days'
order by
  create_time desc;
```

```sql+sqlite
select
  email,
  display_name,
  create_time
from
  onepassword_user
where
  create_time >= datetime('now', '-30 days')
order by
  create_time desc;
```

### List active users outside the company email domain
Find active users whose email address is not on your company domain, such as personal or contractor addresses.

```sql+postgres
select
  email,
  display_name,
  create_time
from
  onepassword_user
where
  state = 'ACTIVE'
  and email not like '%@example.com';
```

```sql+sqlite
select
  email,
  display_name,
  create_time
from
  onepassword_user
where
  state = 'ACTIVE'
  and email not like '%@example.com';
```

### Count users by state
See how many users are active, suspended or in another state such as pending confirmation.

```sql+postgres
select
  coalesce(state, 'OTHER') as state,
  count(*) as users
from
  onepassword_user
group by
  state;
```

```sql+sqlite
select
  coalesce(state, 'OTHER') as state,
  count(*) as users
from
  onepassword_user
group by
  state;
```
