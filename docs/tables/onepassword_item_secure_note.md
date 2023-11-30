---
title: "Steampipe Table: onepassword_item_secure_note - Query 1Password Secure Notes using SQL"
description: "Allows users to query Secure Notes in 1Password, specifically the details of secure notes stored in the 1Password vaults."
---

# Table: onepassword_item_secure_note - Query 1Password Secure Notes using SQL

1Password is a password manager developed by AgileBits Inc. It provides a place for users to store various types of sensitive information, including web logins, credit card information, and secure notes. A Secure Note in 1Password is a safe and convenient way to store and keep track of important information, such as server details, emergency instructions, or plans for world domination.

## Table Usage Guide

The `onepassword_item_secure_note` table provides insights into Secure Notes within 1Password. As a Security Analyst, explore note-specific details through this table, including the note's title, vault ID, and associated metadata. Utilize it to uncover information about Secure Notes, such as those with specific tags, the notes belonging to a specific vault, and the verification of note details.

## Examples

### Basic info
Explore which secure notes in your 1Password vault have been updated recently. This allows you to keep track of changes and ensure your information is up-to-date.

```sql
select
  id,
  title,
  vault_id,
  created_at,
  notes_plain,
  tags,
  updated_at,
  version
from
  onepassword_item_secure_note;
```

### List secure notes of a particular vault
Explore secure notes within a specific vault to review and manage your sensitive information more effectively. This can be particularly useful for maintaining data integrity and ensuring the correct information is stored in the right vault.

```sql
select
  s.id,
  s.title,
  notes_plain,
  s.created_at,
  s.updated_at,
  favorite
from
  onepassword_item_secure_note as s,
  onepassword_vault as v
where
  s.vault_id = v.id
  and v.name = 'my-creds';
```

### Show secure notes that contain a specific tag
Explore secure notes that are associated with a specific tag. This could be useful for quickly identifying all notes related to a particular topic or category, such as those tagged for use with Amazon.

```sql
select
  id,
  title,
  vault_id,
  created_at,
  notes_plain,
  tags,
  updated_at,
  version
from
  onepassword_item_secure_note
where
  tags @> '["amazon-use"]';
```

### List secure notes that are marked as favourite
Discover the segments that consist of secure notes marked as favourites, allowing you to easily track and manage your most important notes. This is useful for prioritizing and accessing your key information quickly and efficiently.

```sql
select
  id,
  title,
  vault_id,
  created_at,
  notes_plain,
  tags,
  updated_at,
  version
from
  onepassword_item_secure_note
where
  favorite;
```