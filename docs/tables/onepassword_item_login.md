# Table: onepassword_item_login

Login items include fields like username, password, website, and one-time password. This category is used for saving or filling logins.

## Examples

### Basic info

```sql
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

```sql
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

### List logins of a particular vault

```sql
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

```sql
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

### List logins with password less than 8 characters

```sql
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
