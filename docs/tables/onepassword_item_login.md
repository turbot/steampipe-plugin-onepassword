# Table: onepassword_item_login

Login items include fields like username, password, website, and one-time password. This category is used for saving or filling logins.

## Examples

### Basic info

```sql
select
  id,
  title,
  password,
  created_at,
  updated_at,
  tags 
from
  onepassword_item_login
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
  onepassword_item_login,
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
  onepassword_item_login as p,
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
  onepassword_item_login
where
  tags @> '["amazon-use"]';
```

### List items with password less than 8 characters

```sql
select
  id,
  title,
  password,
  created_at,
  tags
from
  onepassword_item_login where LENGTH(password) < 8;
```

### List of passwords that are not unique

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
      onepassword_item_login 
    group by
      password 
    having
      count(*) > 1 
  )
  p1 
  join
    onepassword_item_login p2 
    on p1.password = p2.password;
```
