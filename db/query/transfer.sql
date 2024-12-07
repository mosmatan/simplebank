-- name: CreateTranfer :one
INSERT INTO transfers (
  from_acount_id,
  to_acount_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTranfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTranfers :many
SELECT * FROM transfers
WHERE from_acount_id = $1 OR to_acount_id = $2
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTranfer :one
UPDATE transfers
SET from_acount_id = $2,
to_acount_id = $3,
amount = $4
WHERE id = $1
RETURNING * ;

-- name: DeleteTranfer :exec
DELETE FROM transfers 
WHERE id = $1;