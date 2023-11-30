---
title: "Steampipe Table: onepassword_vault - Query 1Password Vaults using SQL"
description: "Allows users to query Vaults in 1Password, specifically the details and metadata of each vault, providing insights into the organization and management of secure information."
---

# Table: onepassword_vault - Query 1Password Vaults using SQL

A Vault in 1Password is a secure location where users can store sensitive information, such as passwords, credit card details, and secure notes. Vaults help in organizing information, making it easier to find and share. They provide an additional layer of security as each vault can have separate permissions.

## Table Usage Guide

The `onepassword_vault` table provides insights into Vaults within 1Password. As a security administrator, explore vault-specific details through this table, including vault UUID, name, description, and associated metadata. Utilize it to uncover information about vaults, such as their access controls, the type of information stored, and the organization of secure data.

## Examples

### Basic info
Explore the creation and modification details of your 1Password vaults. This allows you to keep track of vault updates and understand their types and descriptions for better management.

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
Discover the segments that have recently added secure vaults in the past month. This is useful for tracking the creation of new storage spaces for sensitive information.

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
Discover the segments that have zero items stored in them to understand where storage is currently unused or potentially unutilized. This can be useful in identifying areas for better resource allocation or for detecting any anomalies.

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
Discover the segments that comprise your personal vaults. This query is useful for gaining insights into your individual vaults, including when they were created and last updated, and what items they contain.

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