---
title: "Steampipe Table: onepassword_item_software_license - Query 1Password Software Licenses using SQL"
description: "Allows users to query Software Licenses in 1Password, specifically the item details, vault details, and category details, providing insights into software license management."
---

# Table: onepassword_item_software_license - Query 1Password Software Licenses using SQL

1Password is a password manager developed by AgileBits Inc. It provides a secure and convenient way to store and manage passwords, software licenses, and other sensitive information. 1Password's Software License feature enables users to securely store and manage all their software licenses in one place.

## Table Usage Guide

The `onepassword_item_software_license` table provides insights into Software Licenses within 1Password. As an IT administrator, explore license-specific details through this table, including item details, vault details, and category details. Utilize it to manage and track all software licenses, ensuring compliance and preventing unauthorized usage.

## Examples

### Basic info
Analyze the settings to understand the status and details of your software licenses, including their creation and last update dates, and whether they are marked as favorites. This can help in managing and tracking your licenses effectively.

```sql+postgres
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

```sql+sqlite
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
Explore the software licenses stored in a specific vault to manage and keep track of your software assets effectively. This is particularly useful for auditing purposes and ensuring compliance with software licensing agreements.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have a specific tag within your software licenses. This can be particularly useful when you want to categorize and manage your licenses based on certain criteria or attributes.

```sql+postgres
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
  tags @> '["amazon-use"]';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List software licenses that are marked as favourite
Explore which software licenses are marked as favourites to better manage your software assets. This helps in quickly identifying your most important licenses and ensuring they are up-to-date.

```sql+postgres
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

```sql+sqlite
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
  favorite = 1;
```
