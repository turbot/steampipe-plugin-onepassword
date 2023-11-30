---
title: "Steampipe Table: onepassword_item_medical_record - Query OnePassword Medical Records using SQL"
description: "Allows users to query Medical Records in OnePassword, specifically the detailed information about each medical record, providing insights into personal health information."
---

# Table: onepassword_item_medical_record - Query OnePassword Medical Records using SQL

OnePassword is a password manager developed by AgileBits Inc. It provides a platform where users can store various passwords, software licenses, and other sensitive information in a virtual vault locked with a PBKDF2-guarded master password. OnePassword's Medical Records feature allows users to securely store and manage sensitive health information.

## Table Usage Guide

The `onepassword_item_medical_record` table provides insights into Medical Records within OnePassword. As a security analyst, explore medical record-specific details through this table, including the record's title, type, and associated metadata. Utilize it to uncover information about medical records, such as the record's UUID, vault UUID, and the time the record was last modified.

## Examples

### Basic info
Explore the details of medical records, such as the date, dosage, and healthcare professional involved. This can be useful for gaining insights into patient care and medication management.

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
Explore the medical records stored in a specific vault to gain insights into healthcare details like medication dosage, attending professional, and the reason for treatment. This can be useful for managing and reviewing personal health information in a secure manner.

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

### Show medical records that have been tagged as `important`
Identify instances where medical records have been marked as important. This can be useful for prioritizing patient care and ensuring critical information is not overlooked.

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
Explore medical history for a particular patient, including details like medication dosages, healthcare professional interactions, and reasons for treatment. This can help in understanding the patient's health trajectory and making informed healthcare decisions.

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
Explore medical records to identify instances where a specific medication, such as Aspirin, has been prescribed. This can be useful in tracking medication usage patterns or for patient health management.

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