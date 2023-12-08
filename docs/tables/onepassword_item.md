---
title: "Steampipe Table: onepassword_item - Query 1Password Items using SQL"
description: "Allows users to query 1Password Items, specifically retrieving details about individual items stored in the 1Password vaults."
---

# Table: onepassword_item - Query 1Password Items using SQL

1Password is a password management service that stores sensitive information, including passwords, software licenses, and notes, in a secure and encrypted format. It allows users to create and manage multiple vaults, each containing different sets of items. These items can range from login credentials, secure notes, credit card information, to identities.

## Table Usage Guide

The `onepassword_item` table provides insights into 1Password Items within 1Password. As a security analyst, explore item-specific details through this table, including categories, vault IDs, and associated metadata. Utilize it to uncover information about items, such as their categories, the vaults they belong to, and the intricacies of their details.

## Examples

### Basic info
Explore which items in your OnePassword vault have been recently updated or edited. This can help you keep track of changes and ensure your data remains secure.

```sql+postgres
select
  id,
  vault_id,
  title,
  created_at,
  updated_at,
  last_edited_by,
  tags
from
  onepassword_item;
```

```sql+sqlite
select
  id,
  vault_id,
  title,
  created_at,
  updated_at,
  last_edited_by,
  tags
from
  onepassword_item;
```

### List items that have been updated in the last 30 days
Explore which items have been modified in the past month to keep track of recent changes.

```sql+postgres
select
  id,
  vault_id,
  title,
  created_at,
  updated_at,
  last_edited_by,
  tags
from
  onepassword_item
where
  updated_at > now() - interval '30 day';
```

```sql+sqlite
select
  id,
  vault_id,
  title,
  created_at,
  updated_at,
  last_edited_by,
  tags
from
  onepassword_item
where
  updated_at > datetime('now', '-30 day');
```

### List items with production tag
Identify instances where items are tagged as 'production' to focus on operational elements.

```sql+postgres
select
  id,
  title,
  category,
  version,
  tags
from
  onepassword_item
where
  tags @> '["production"]';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List the fields of all items with a specific section
Analyze the fields of items within a specific section, such as 'Metadata', to gain insights into categorized data.

```sql+postgres
select
  title,
  jsonb_pretty(f) as field
from
  onepassword_item,
  jsonb_array_elements(fields) as f
where
  f -> 'section' ->> 'label' = 'Metadata';
```

```sql+sqlite
select
  title,
  f.value as field
from
  onepassword_item,
  json_each(fields) as f
where
  json_extract(json_extract(f.value, '$.section'), '$.label') = 'Metadata';
```