---
title: "Steampipe Table: onepassword_audit_event - Query 1Password Audit Events using SQL"
description: "Allows users to query the actions team members performed in a 1Password Business account through the 1Password Events API, such as user, vault access and security setting changes."
---

# Table: onepassword_audit_event - Query 1Password Audit Events using SQL

The 1Password Events API records the actions team members take in a 1Password Business account, such as inviting or suspending users, granting or revoking vault access, and changing account security settings like multi-factor authentication and SSO.

## Table Usage Guide

The `onepassword_audit_event` table provides insights into administrative activity in your 1Password account. As a security or compliance analyst, use it to review who changed what and when, detect security settings being turned off, and track changes to user accounts and vault access.

**Important Notes**
- This table requires `events_token` in the connection configuration: a bearer token from a 1Password [Events Reporting integration](https://support.1password.com/events-reporting/) with access to audit events. Set `events_url` if your account is not on `1password.com`.
- The table returns the last 24 hours of events unless you constrain `timestamp` in the `where` clause with `>`, `>=`, `<` or `<=`.

## Examples

### Basic info
Review the actions performed in the last 24 hours and who performed them.

```sql+postgres
select
  uuid,
  timestamp,
  action,
  object_type,
  object_uuid,
  actor_email
from
  onepassword_audit_event;
```

```sql+sqlite
select
  uuid,
  timestamp,
  action,
  object_type,
  object_uuid,
  actor_email
from
  onepassword_audit_event;
```

### List security settings turned off in the last 30 days
Detect multi-factor authentication, Duo or SSO being disabled, which weakens how your team signs in.

```sql+postgres
select
  timestamp,
  action,
  object_type,
  actor_email
from
  onepassword_audit_event
where
  timestamp >= now() - interval '30 days'
  and action in ('disblmfa', 'disblduo', 'disblsso')
order by
  timestamp desc;
```

```sql+sqlite
select
  timestamp,
  action,
  object_type,
  actor_email
from
  onepassword_audit_event
where
  timestamp >= datetime('now', '-30 days')
  and action in ('disblmfa', 'disblduo', 'disblsso')
order by
  timestamp desc;
```

### List users suspended in the last 24 hours
Confirm that offboarded team members were suspended, and by whom.

```sql+postgres
select
  timestamp,
  actor_email,
  object_details ->> 'email' as suspended_user
from
  onepassword_audit_event
where
  action = 'suspend'
  and object_type = 'user';
```

```sql+sqlite
select
  timestamp,
  actor_email,
  json_extract(object_details, '$.email') as suspended_user
from
  onepassword_audit_event
where
  action = 'suspend'
  and object_type = 'user';
```

### List vault access granted or revoked in the last 7 days
Track changes to who can open which vaults.

```sql+postgres
select
  timestamp,
  action,
  object_uuid as vault_uuid,
  aux_details ->> 'email' as user_email,
  actor_email
from
  onepassword_audit_event
where
  timestamp >= now() - interval '7 days'
  and object_type = 'vault'
  and action in ('grant', 'revoke')
order by
  timestamp desc;
```

```sql+sqlite
select
  timestamp,
  action,
  object_uuid as vault_uuid,
  json_extract(aux_details, '$.email') as user_email,
  actor_email
from
  onepassword_audit_event
where
  timestamp >= datetime('now', '-7 days')
  and object_type = 'vault'
  and action in ('grant', 'revoke')
order by
  timestamp desc;
```

### Count actions by team member in the last 30 days
See which team members perform the most administrative actions, to review whether their access is still appropriate.

```sql+postgres
select
  actor_email,
  count(*) as actions
from
  onepassword_audit_event
where
  timestamp >= now() - interval '30 days'
group by
  actor_email
order by
  actions desc;
```

```sql+sqlite
select
  actor_email,
  count(*) as actions
from
  onepassword_audit_event
where
  timestamp >= datetime('now', '-30 days')
group by
  actor_email
order by
  actions desc;
```
