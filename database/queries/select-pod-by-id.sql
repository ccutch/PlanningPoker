SELECT
    name,
    strategy,
    private,
    status,
    created_at,
    updated_at
FROM pods
WHERE id = $1

