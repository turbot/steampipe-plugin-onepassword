---
title: "Steampipe Table: onepassword_signin_attempt - Query 1Password Sign-in Attempts using SQL"
description: "Allows users to query sign-in attempts to a 1Password Business account through the 1Password Events API, including failed attempts, the user, the client and the approximate location."
---

# Table: onepassword_signin_attempt - Query 1Password Sign-in Attempts using SQL

The 1Password Events API reports every successful and failed attempt to sign in to a 1Password Business account. Each attempt records its outcome, the user it was for, the 1Password app and platform used, and the approximate location it came from.

## Table Usage Guide

The `onepassword_signin_attempt` table provides insights into sign-ins to your 1Password account. As a security analyst, use it to review failed sign-ins, spot repeated MFA failures or attempts blocked by firewall rules, and see where and from which clients your team signs in.

**Important Notes**
- This table requires `events_token` in the connection configuration: a bearer token from a 1Password [Events Reporting integration](https://support.1password.com/events-reporting/) with access to sign-in attempts. Set `events_url` if your account is not on `1password.com`.
- The table returns the last 24 hours of attempts unless you constrain `timestamp` in the `where` clause with `>`, `>=`, `<` or `<=`.

## Examples

### Basic info
Review the sign-in attempts of the last 24 hours, with their outcome and the user they were for.

```sql+postgres
select
  uuid,
  timestamp,
  category,
  type,
  target_user_email,
  client_ip_address,
  country
from
  onepassword_signin_attempt;
```

```sql+sqlite
select
  uuid,
  timestamp,
  category,
  type,
  target_user_email,
  client_ip_address,
  country
from
  onepassword_signin_attempt;
```

### List failed sign-ins in the last 7 days
Find every sign-in attempt in the past week that did not succeed, to investigate credential stuffing or misconfigured clients.

```sql+postgres
select
  timestamp,
  category,
  type,
  target_user_email,
  client_ip_address,
  country
from
  onepassword_signin_attempt
where
  timestamp >= now() - interval '7 days'
  and category <> 'success'
order by
  timestamp desc;
```

```sql+sqlite
select
  timestamp,
  category,
  type,
  target_user_email,
  client_ip_address,
  country
from
  onepassword_signin_attempt
where
  timestamp >= datetime('now', '-7 days')
  and category <> 'success'
order by
  timestamp desc;
```

### List users with ten or more failed sign-ins in the last 24 hours
Identify accounts that may be under attack, or users who are locked out and need help.

```sql+postgres
select
  target_user_email,
  count(*) as failures
from
  onepassword_signin_attempt
where
  category in ('credentials_failed', 'mfa_failed')
group by
  target_user_email
having
  count(*) >= 10
order by
  failures desc;
```

```sql+sqlite
select
  target_user_email,
  count(*) as failures
from
  onepassword_signin_attempt
where
  category in ('credentials_failed', 'mfa_failed')
group by
  target_user_email
having
  count(*) >= 10
order by
  failures desc;
```

### Count sign-in attempts by country
See where sign-in attempts came from over the past 30 days, to spot activity from countries your team does not work in.

```sql+postgres
select
  country,
  count(*) filter (where category = 'success') as successful,
  count(*) filter (where category <> 'success') as failed
from
  onepassword_signin_attempt
where
  timestamp >= now() - interval '30 days'
group by
  country
order by
  failed desc;
```

```sql+sqlite
select
  country,
  sum(case when category = 'success' then 1 else 0 end) as successful,
  sum(case when category <> 'success' then 1 else 0 end) as failed
from
  onepassword_signin_attempt
where
  timestamp >= datetime('now', '-30 days')
group by
  country
order by
  failed desc;
```

### List the apps and platforms used to sign in
Review which 1Password apps and operating systems your team signs in with, for example to find outdated app versions.

```sql+postgres
select
  client ->> 'app_name' as app_name,
  client ->> 'app_version' as app_version,
  client ->> 'os_name' as os_name,
  count(*) as sign_ins
from
  onepassword_signin_attempt
where
  category = 'success'
group by
  app_name,
  app_version,
  os_name
order by
  sign_ins desc;
```

```sql+sqlite
select
  json_extract(client, '$.app_name') as app_name,
  json_extract(client, '$.app_version') as app_version,
  json_extract(client, '$.os_name') as os_name,
  count(*) as sign_ins
from
  onepassword_signin_attempt
where
  category = 'success'
group by
  app_name,
  app_version,
  os_name
order by
  sign_ins desc;
```
