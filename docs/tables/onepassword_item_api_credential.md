# Table: onepassword_item_api_credential

API Credential items include fields like username, credential, and hostname needed to access API tools.

## Examples

### Basic info

```sql
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential;
```

### List API credentials stored in a particular vault

```sql
select
  c.id,
  c.title,
  username,
  credential,
  c.created_at,
  favorite
from
  onepassword_item_api_credential as c,
  onepassword_vault as v
where
  c.vault_id = v.id
  and v.name = 'my-creds';
```

### Show API credentials that contain a specific tag

```sql
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  tags @> '["nps-creds"]';
```

### List API credentials that were updated within the last month

```sql
select
  id,
  title,
  username,
  credential,
  created_at,
  favorite
from
  onepassword_item_api_credential
where
  updated_at >= now() - interval '1 month';
```
