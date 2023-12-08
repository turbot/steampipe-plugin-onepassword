---
title: "Steampipe Table: onepassword_item_api_credential - Query 1Password API Credentials using SQL"
description: "Allows users to query 1Password API Credentials, providing detailed information about the API credentials stored in the 1Password vault."
---

# Table: onepassword_item_api_credential - Query 1Password API Credentials using SQL

1Password API Credentials is a resource within 1Password that securely stores and manages credentials for API usage. It provides an organized and secure way to store API keys, tokens, and other credentials, ensuring that they are easily accessible but protected. 1Password API Credentials contributes to maintaining a secure and efficient API management process.

## Table Usage Guide

The `onepassword_item_api_credential` table provides insights into API credentials stored within 1Password. As a developer or DevOps engineer, explore credential-specific details through this table, including the server, the key, and the password. Utilize it to manage and track API credentials, ensuring that they are securely stored and correctly used.

## Examples

### Basic info
Explore the OnePassword API credentials to gain insights into user activity and preferences. This can be useful for assessing user behavior and identifying potential security risks.

```sql+postgres
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential;
```

```sql+sqlite
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential;
```

### List API credentials stored in a particular vault
Discover the segments that have stored API credentials in a specific vault. This can be useful to audit the security of your credentials and ensure they are stored in the correct vault.

```sql+postgres
select
  c.id,
  c.title,
  username,
  credential,
  c.created_at,
  favorite
from
  onepassword_item_api_credential as c,
  onepassword_vault as v
where
  c.vault_id = v.id
  and v.name = 'my-creds';
```

```sql+sqlite
select
  c.id,
  c.title,
  username,
  credential,
  c.created_at,
  favorite
from
  onepassword_item_api_credential as c,
  onepassword_vault as v
where
  c.vault_id = v.id
  and v.name = 'my-creds';
```

### Show API credentials that contain a specific tag
Discover the segments that contain specific tags within your API credentials. This can help you quickly identify and manage credentials associated with certain projects or tasks.

```sql+postgres
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  tags @> '["nps-creds"]';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List API credentials that were updated within the last month
Determine the API credentials that have been recently updated to ensure they are current and secure. This can be beneficial for maintaining system integrity and mitigating potential security risks.

```sql+postgres
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  updated_at >= now() - interval '1 month';
```

```sql+sqlite
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  updated_at >= datetime('now', '-1 month');
```

### List API credentials marked as favourite
Explore which API credentials have been marked as favourites, to easily access them for future use. This could be particularly useful in managing and prioritizing multiple API credentials.

```sql+postgres
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  favorite;
```

```sql+sqlite
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  favorite = 1;
```