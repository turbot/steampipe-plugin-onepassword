# Table: onepassword_item_password

Password items include a field for a password. If you enter a username when creating a Password item, it will automatically convert to a Login item.

## Examples

### Basic info

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
