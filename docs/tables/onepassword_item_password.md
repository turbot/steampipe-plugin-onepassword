---
title: "Steampipe Table: onepassword_item_password - Query 1Password Password Items using SQL"
description: "Allows users to query Password Items in 1Password, specifically the details of saved passwords, providing insights into password management and security."
---

# Table: onepassword_item_password - Query 1Password Password Items using SQL

1Password is a password manager that provides a place for users to store various passwords, software licenses, and other sensitive information in a virtual vault locked with a PBKDF2-guarded master password. By default, this encrypted vault is stored on the company's servers for a monthly fee. It provides secure, encrypted storage for sensitive information, with the ability to share items and manage permissions through vaults.

## Table Usage Guide

The `onepassword_item_password` table provides insights into Password Items within 1Password. As a security professional, explore password-specific details through this table, including password content, vault ID, and associated metadata. Utilize it to monitor and manage password security, such as identifying weak passwords, tracking password reuse, and verifying password security policies.

## Examples

### Basic info
Explore which OnePassword items have a specific tag, allowing you to better manage and categorize your passwords. This is particularly useful for identifying outdated or rarely used passwords that may need updating or removal.

```sql
select
  id,
  title,
  password,
  created_at,
  tags
from
  onepassword_item_password;
```

### List passwords along with website details
Identify password information linked to specific websites, enabling you to review and manage your digital security effectively.

```sql
select
  id,
  title,
  password,
  jsonb_pretty(u -> 'href') as website,
  created_at,
  tags
from
  onepassword_item_password,
  jsonb_array_elements(urls) as u;
```

### List passwords of a particular vault
Discover the segments that contain specific vault passwords. This is useful for managing and auditing security credentials within a particular vault.

```sql
select
  p.id,
  p.title,
  password,
  p.created_at,
  p.tags
from
  onepassword_item_password as p,
  onepassword_vault as v
where
  p.vault_id = v.id
  and v.name = 'my-creds';
```

### Show passwords that contain a specific tag
Explore which passwords are associated with a specific tag. This can be useful for identifying and managing passwords related to a particular project or service.

```sql
select
  id,
  title,
  password,
  created_at,
  tags
from
  onepassword_item_password
where
  tags @> '["amazon-use"]';
```

### List passwords that are less than 8 characters in length
Identify instances where passwords are less than 8 characters, which can be a potential security risk. This helps in improving security measures by enforcing stronger password policies.

```sql
select
  id,
  title,
  password,
  created_at,
  tags
from
  onepassword_item_password
where
  length(password) < 8;
```

### List passwords that are not unique
Explore which passwords are not unique in your system, helping to highlight potential security risks associated with password duplication. This can be useful in identifying and mitigating potential vulnerabilities in your security infrastructure.

```sql
select
  p2.id,
  p2.vault_id,
  p2.title,
  p1.password
from
  (
    select
      password,
      count(*) as count
    from
      onepassword_item_password
    group by
      password
    having
      count(*) > 1
  )
  p1
  join onepassword_item_password p2 on p1.password = p2.password;
```