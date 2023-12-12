---
title: "Steampipe Table: onepassword_item_login - Query 1Password Login Items using SQL"
description: "Allows users to query 1Password Login Items, providing detailed information about each login item stored within the 1Password vaults."
---

# Table: onepassword_item_login - Query 1Password Login Items using SQL

1Password is a password management service that stores sensitive information, such as passwords, secure notes, and software licenses, in secure, encrypted vaults. This service is widely used by individuals and businesses to manage and protect their sensitive information. A Login Item in 1Password is a type of item that contains the login details for a specific website or service, including the username and password.

## Table Usage Guide

The `onepassword_item_login` table provides insights into Login Items within 1Password. As a security analyst or system administrator, explore each Login Item's details through this table, including the associated username, password, website, and other related information. Utilize it to manage and monitor the login information for various services, ensuring the security and integrity of sensitive data.

## Examples

### Basic info
Gain insights into the creation and modification dates of your login items, along with any tags associated with them. This allows for easy tracking and management of your login credentials over time.

```sql+postgres
select
  id,
  title,
  username,
  password,
  created_at,
  updated_at,
  tags
from
  onepassword_item_login;
```

```sql+sqlite
select
  id,
  title,
  username,
  password,
  created_at,
  updated_at,
  tags
from
  onepassword_item_login;
```

### List logins along with website details
Explore the details of your saved logins, including the associated websites, to better manage your online accounts. This can help in tracking account creation dates and ensuring password security.

```sql+postgres
select
  id,
  title,
  username,
  password,
  jsonb_pretty(u -> 'href') as website,
  created_at
from
  onepassword_item_login,
  jsonb_array_elements(urls) as u;
```

```sql+sqlite
select
  id,
  title,
  username,
  password,
  u.value as website,
  created_at
from
  onepassword_item_login,
  json_each(urls) as u;
```

### List logins of a particular vault
Explore which logins are associated with a specified secure vault. This is useful to assess the elements within a specific vault for better management and security.

```sql+postgres
select
  p.id,
  p.title,
  username,
  password,
  p.created_at,
  p.tags
from
  onepassword_item_login as p,
  onepassword_vault as v
where
  p.vault_id = v.id
  and v.name = 'my-creds';
```

```sql+sqlite
select
  p.id,
  p.title,
  username,
  password,
  p.created_at,
  p.tags
from
  onepassword_item_login as p,
  onepassword_vault as v
where
  p.vault_id = v.id
  and v.name = 'my-creds';
```

### Show logins that contain a specific tag
Explore which login items are associated with a specific tag to better manage and categorize your credentials. This can be particularly useful for identifying and organizing logins related to a certain project or platform, such as Amazon.

```sql+postgres
select
  id,
  title,
  username,
  password,
  created_at,
  tags
from
  onepassword_item_login
where
  tags @> '["amazon-use"]';
```

```sql+sqlite
Error: SQLite does not support array operations.
```

### List logins with password length less than 8 characters
Identify instances where user passwords may be less secure due to their short length. This is useful for auditing account security and identifying potential vulnerabilities.

```sql+postgres
select
  id,
  title,
  username,
  password,
  created_at
from
  onepassword_item_login
where
  length(password) < 8;
```

```sql+sqlite
select
  id,
  title,
  username,
  password,
  created_at
from
  onepassword_item_login
where
  length(password) < 8;
```