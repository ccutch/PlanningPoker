SELECT
    id,
    name,
    owner,
    created_at,
    updated_at
FROM players
WHERE pod_id = $1

