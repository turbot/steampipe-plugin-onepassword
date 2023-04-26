# Table: onepassword_item_software_license

Software License items include fields like version, license key, registered email, and download page.

## Examples

### Basic info

```sql
select
  id,
  title,
  license_key,
  created_at,
  updated_at,
  favorite
from
  onepassword_item_software_license;
```

### List software licenses of a particular vault

```sql
select
  s.id,
  s.title,
  license_key,
  s.created_at,
  favorite
from
  onepassword_item_software_license as s,
  onepassword_vault as v
where
  s.vault_id = v.id
  and v.name = 'my-creds';
```

### Show software licenses that contain a specific tag

```sql
select
  id,
  title,
  license_key,
  created_at,
  updated_at,
  favorite,
  tags
from
  onepassword_item_software_license
where
  tags @ > '["amazon-use"]';
```

### List favorite software licenses

```sql
select
  id,
  title,
  license_key,
  created_at,
  updated_at,
  favorite,
  tags
from
  onepassword_item_software_license
where
  favorite;
```
