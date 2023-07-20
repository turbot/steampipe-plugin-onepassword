# Table: onepassword_item_file

1Password lets you securely store your most important files, so they’re always available when you need them.

- You **_must_** specify `item_id` in a `where` clause in order to use this table.

## Examples

### Basic info

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
