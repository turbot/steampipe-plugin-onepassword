# Table: onepassword_item

An item in 1Password is a container for securely storing a piece of sensitive information, such as a password, login credentials, credit card number, secure note, or other data. Each item in 1Password can have various attributes and fields depending on the type of information it stores, such as a website URL, username, password, expiration date, and more.

Items can be organized into various categories or folders, such as Personal, Work, or Finance, for easy management and access. They can also be tagged, favorited, or searched for quickly.

## Examples

### Basic info

```sql
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

```sql
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

## List items with production tag

```sql
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

## List the fields of all items with a specific section

```sql
select
  title,
  jsonb_pretty(f) as field
from
  onepassword_item,
  jsonb_array_elements(fields) as f
where
  f -> 'section' ->> 'label' = 'Metadata';
```
