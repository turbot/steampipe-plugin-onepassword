# Table: onepassword_vault

Every item you save in 1Password is stored in a vault. You can use vaults to organize your items and share with others. Items in a vault are available to everyone with access to that vault.

## Examples

### Basic info

```sql
select
  id,
  name,
  created_at,
  description,
  type,
  updated_at
from
  onepassword_vault;
```

### List vaults created in the last 30 days

```sql
select
  id,
  name,
  created_at,
  description,
  type,
  updated_at
from
  onepassword_vault
where
  created_at >= now() - interval '30 day';
```

### Show vaults with zero items

```sql
select
  id,
  name,
  created_at,
  items,
  description,
  type,
  updated_at
from
  onepassword_vault
where
  items = 0;
```

### Show personal vaults

```sql
select
  id,
  name,
  created_at,
  items,
  description,
  type,
  updated_at
from
  onepassword_vault
where
  type = 'PERSONAL';
```
