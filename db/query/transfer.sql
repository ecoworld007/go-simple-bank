-- name: CreateTranfer :one
INSERT INTO transfers (
  amount,
  from_account_id,
  to_account_id
)
VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTranfer :one
SELECT * FROM transfers 
where id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTranfer :one
UPDATE transfers 
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTranfer :exec
DELETE FROM transfers
where id = $1;