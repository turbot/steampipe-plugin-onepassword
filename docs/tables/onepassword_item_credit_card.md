# Table: onepassword_item_credit_card

Credit Card items include fields like card number, verification number, and expiry date for your credit and debit card information. You can use this category to fill card information in your browser.

## Examples

### Basic info

```sql
select
  id,
  title,
  card_holder,
  credit_card_number,
  expiry_date,
  created_at,
  favorite
from
  onepassword_item_credit_card;
```

### List credit cards of a particular vault

```sql
select
  c.id,
  c.title,
  card_holder,
  credit_card_number,
  expiry_date,
  c.created_at,
  favorite
from
  onepassword_item_credit_card as c,
  onepassword_vault as v
where
  c.vault_id = v.id
  and v.name = 'my-creds';
```

### Show credit cards that contain a specific tag

```sql
select
  id,
  title,
  card_holder,
  credit_card_number,
  expiry_date,
  created_at,
  favorite
from
  onepassword_item_credit_card
where
  tags @> '["amazon-use"]';
```

### List expired credit cards

```sql
select
  id,
  title,
  card_holder,
  credit_card_number,
  expiry_date,
  created_at,
  favorite
from
  onepassword_item_credit_card
where
  expiry_date < now();
```

### List credit cards that are marked as favourite

```sql
select
  id,
  title,
  card_holder,
  credit_card_number,
  expiry_date,
  created_at,
  favorite
from
  onepassword_item_credit_card
where
  favorite;
```