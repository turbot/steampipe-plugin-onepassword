# Table: onepassword_item_medical_record

Medical Record items include fields like date, location, healthcare professional, and reason for visit. You can use this category to save any of your health-related information.

## Examples

### Basic info

```sql
select
  id,
  title,
  date,
  dosage,
  healthcare_professional,
  location,
  medication,
  patient,
  reason
from
  onepassword_item_medical_record;
```

### List medical records of a particular vault

```sql
select
  r.id,
  r.title,
  date,
  dosage,
  healthcare_professional,
  location,
  medication,
  patient,
  reason
from
  onepassword_item_medical_record as r,
  onepassword_vault as v
where
  r.vault_id = v.id
  and v.name = 'my-creds';
```

### Show medical records that have been tagged as "important"

```sql
select
  id,
  title,
  date,
  dosage,
  healthcare_professional,
  location,
  medication,
  patient,
  reason
from
  onepassword_item_medical_record
where
  tags @> '["important"]';
```

### List medical records for a specific patient

```sql
select
  id,
  title,
  date,
  dosage,
  healthcare_professional,
  location,
  medication,
  patient,
  reason
from
  onepassword_item_medical_record
where
  patient like '%sid%';
```

### List medical records that contain a certain medication

```sql
select
  id,
  title,
  date,
  dosage,
  healthcare_professional,
  location,
  medication,
  patient,
  reason
from
  onepassword_item_medical_record
where
  medication like '%aspirin%';
```
