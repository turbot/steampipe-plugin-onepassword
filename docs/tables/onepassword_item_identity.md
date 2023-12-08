---
title: "Steampipe Table: onepassword_item_identity - Query 1Password Identity Items using SQL"
description: "Allows users to query Identity Items in 1Password, providing insights into the identity details stored in the 1Password vault."
---

# Table: onepassword_item_identity - Query 1Password Identity Items using SQL

1Password is a password manager developed by AgileBits Inc. It provides a secure and convenient way to store and manage various types of sensitive information including passwords, identities, and secure notes. An Identity Item in 1Password represents personal information about a user, such as their name, address, phone number, and more.

## Table Usage Guide

The `onepassword_item_identity` table provides insights into Identity Items within 1Password. As a security analyst, explore Identity Item-specific details through this table, including the associated vault, category, and the personal details stored in the item. Utilize it to uncover information about the stored identities, such as the associated tags, the time of creation and modification, and the verification of personal details.

## Examples

### Basic info
Explore which employees belong to a specific department within your company, gaining insights into their job titles and occupations. This can help in understanding the distribution of roles and responsibilities within each department.

```sql+postgres
select
  id,
  title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity;
```

```sql+sqlite
select
  id,
  title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity;
```

### List identities of a particular vault
Explore the profiles stored in a specific secure vault to gain insights into the identities it contains. This can be useful for auditing purposes, ensuring that only the appropriate identities are stored in a particular vault.

```sql+postgres
select
  i.id,
  i.title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity as i,
  onepassword_vault as v
where
  i.vault_id = v.id
  and v.name = 'my-creds';
```

```sql+sqlite
select
  i.id,
  i.title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity as i,
  onepassword_vault as v
where
  i.vault_id = v.id
  and v.name = 'my-creds';
```

### Show identities that contain a specific tag
Explore which identities are associated with a specific tag in order to manage and organize your data more effectively. This can be beneficial in scenarios such as identifying all employees associated with a specific project or department.

```sql+postgres
select
  id,
  title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity
where
  tags @> '["chat-company"]';
```

```sql+sqlite
Error: The corresponding SQLite 
```

### List identities whose birth date is before a certain date
Discover the identities that were born before a specific date. This can be useful for demographic analysis or to identify potential age-related factors in your data.

```sql+postgres
select
  id,
  title,
  first_name,
  last_name,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity
where
  birth_date < '1990-01-01';
```

```sql+sqlite
select
  id,
  title,
  first_name,
  last_name,
  company,
  department,
  gender,
  job_title,
  occupation
from
  onepassword_item_identity
where
  birth_date < '1990-01-01';
```

### Show identities with the job title software engineer
Explore which identities in your database are associated with the job title 'software engineer'. This can be useful for identifying potential candidates for internal job postings or for understanding the distribution of roles within your organization.

```sql+postgres
select
  id,
  title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  occupation
from
  onepassword_item_identity
where
  job_title = 'software engineer';
```

```sql+sqlite
select
  id,
  title,
  first_name,
  last_name,
  birth_date,
  company,
  department,
  gender,
  occupation
from
  onepassword_item_identity
where
  job_title = 'software engineer';
```

### Get the number of identities for each occupation
Explore the distribution of identities across different occupations. This can provide insights into the diversity of roles within your organization.

```sql+postgres
select
  occupation,
  count(*)
from
  onepassword_item_identity
group by
  occupation;
```

```sql+sqlite
select
  occupation,
  count(*)
from
  onepassword_item_identity
group by
  occupation;
```
