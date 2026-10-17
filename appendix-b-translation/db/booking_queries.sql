-- name: InsertOrder :exec
INSERT INTO orders (id, customer_id, status, items, placed_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetOrderByID :one
SELECT id, customer_id, status, items, placed_at FROM orders WHERE id = $1;
