# Table: onepassword_item_identity

Identity items include fields like first name, last name, address, birth date, phone number, email address, and username. You can use these items to fill address information and see email address and username suggestions when you sign up for a new account on a website.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

### List identities whose birth date is before a certain date

```sql
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

```sql
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

```sql
select
  occupation,
  count(*)
from
  onepassword_item_identity
group by
  occupation;
```
