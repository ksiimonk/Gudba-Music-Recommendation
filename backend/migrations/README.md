# SQL migrations

This project keeps database migrations as plain SQL files in this folder.

Naming convention:

- `000001_description.up.sql` for applying a migration
- `000001_description.down.sql` for rolling a migration back

Apply `up` files in ascending order.
Apply `down` files in descending order.

The first migration creates the `users` table.
