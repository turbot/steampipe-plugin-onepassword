# Table: onepassword_item_password

Password items include a field for a password. If you enter a username when creating a Password item, it will automatically convert to a Login item.

## Examples

### Basic info

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
  jsonb_pretty(u -> 'href') as website,
  p.created_at,
  p.tags
from
  onepassword_item_password as p,
  jsonb_array_elements(p.urls) as u,
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
  tags @ > '["amazon-use"]';
```

### List items with password less than 8 characters

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
  jsonb_array_elements(urls) as u where LENGTH(password) < 8;
```

### List of passwords that are not unique

```sql
SELECT
  p2.id,
  p2.vault_id,
  p2.title,
  p1.password 
FROM
  (
    SELECT
      password,
      COUNT(*) as count 
    FROM
      onepassword_item_password 
    GROUP BY
      password 
    HAVING
      COUNT(*) > 1
  )
  p1 
  JOIN
    onepassword_item_password p2
    ON p1.password = p2.password;
```
