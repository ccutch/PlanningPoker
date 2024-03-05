SELECT
    pod_id,
    name,
    owner,
    created_at,
    updated_at
FROM players
WHERE id = $1

