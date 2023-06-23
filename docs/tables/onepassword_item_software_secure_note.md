# Table: onepassword_item_secure_note

Secure Note items contain a text field that can be formatted with Markdown.

## Examples

### Basic info

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
