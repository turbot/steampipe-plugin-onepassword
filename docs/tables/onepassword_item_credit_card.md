---
title: "Steampipe Table: onepassword_item_credit_card - Query OnePassword Credit Card Items using SQL"
description: "Allows users to query Credit Card Items in OnePassword, specifically the card details, providing insights into saved credit card information and their attributes."
---

# Table: onepassword_item_credit_card - Query OnePassword Credit Card Items using SQL

OnePassword is a password manager developed by AgileBits Inc. It provides a place for users to store various passwords, software licenses, and other sensitive information in a virtual vault that is locked with a PBKDF2-guarded master password. The `onepassword_item_credit_card` resource in OnePassword represents the credit card items stored by users in their OnePassword vaults.

## Table Usage Guide

The `onepassword_item_credit_card` table provides insights into credit card items within OnePassword. As a Security Analyst, explore credit card-specific details through this table, including card numbers, cardholder names, and expiry dates. Utilize it to uncover information about saved credit cards, such as the type of cards stored, the frequency of certain card types, and the verification of card details.

## Examples

### Basic info
Discover the segments that hold your credit card information for a quick review or update. This can be particularly useful for keeping track of card expirations or identifying your most frequently used cards.

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
Explore which credit cards are stored in a specific vault to manage and review your financial information efficiently. This can be particularly useful for organizing your personal finance or auditing corporate expense accounts.

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
Explore which credit cards are associated with a specific tag, such as 'Amazon-use'. This can help in organizing and identifying credit cards that are used for specific platforms or purposes.

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
Explore which credit cards have expired to ensure secure and valid transactions. This is crucial to prevent any financial discrepancies or fraudulent activities.

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

### List credit cards marked as favourite
Discover the segments that have marked their credit cards as favourite. This can assist in understanding user preferences and habits, potentially informing business strategies or marketing efforts.

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