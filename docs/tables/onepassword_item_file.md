---
title: "Steampipe Table: onepassword_item_file - Query OnePassword File Items using SQL"
description: "Allows users to query File Items in OnePassword, providing detailed information about each file item stored in the OnePassword vaults."
---

# Table: onepassword_item_file - Query OnePassword File Items using SQL

OnePassword is a password management service that stores sensitive information, including File Items, in a secure and encrypted format. File Items in OnePassword include any files that users have uploaded and stored in their OnePassword vaults for safekeeping. This includes a variety of file types, such as images, documents, and more, all of which are securely encrypted and only accessible to authorized users.

## Table Usage Guide

The `onepassword_item_file` table provides insights into File Items within OnePassword. As a security analyst, explore file-specific details through this table, including the file's unique identifier, its associated vault, and other metadata. Utilize it to uncover information about stored files, such as their creation and modification dates, to assist in auditing and compliance checks.

## Examples

### Basic info
Explore the details of a specific item in a password management system. This allows you to understand the size and location of the item, which can be useful for managing storage and organization within the system.

```sql
select
  id,
  name,
  item_id,
  vault_id,
  content_path,
  size
from
  onepassword_item_file
where
  item_id = 'kvmaoszyhzbvze6g5t6vr6qg2a1';
```

### List all files of a particular vault
This example allows you to identify all the files associated with a specific vault in the 1Password service. It's particularly useful for auditing purposes, ensuring you have a comprehensive list of all files stored in a particular vault.

```sql
select
  f.id as file_id,
  f.name as file_name,
  i.title as item_title,
  content_path,
  size
from
  onepassword_item_file as f,
  onepassword_item as i,
  onepassword_vault as v
where
  f.item_id = i.id
  and f.vault_id = v.id
  and v.name = 'Venu-SteampipeTest';
```

### Show file contents of all items
Explore the contents of all items in a system, helping you gain insights into data organization and identify potential areas for cleanup or reorganization. This could be particularly useful for auditing purposes or data management initiatives.

```sql
select
  f.name as file_name,
  i.title as item_title,
  content_path,
  jsonb_pretty(content) as content
from
  onepassword_item_file as f,
  onepassword_item as i
where
  f.item_id = i.id;
```