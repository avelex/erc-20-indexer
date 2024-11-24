-- name: SaveTransfer :exec
INSERT INTO erc20_transfers (timestamp, sender, recipient, amount, tx_hash, token_address) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetLastEvents :many
SELECT * FROM erc20_transfers ORDER BY timestamp DESC LIMIT $1;
