SELECT
    pod_id,
    prompt,
    status,
    result,
    created_at,
    updated_at
FROM topics
WHERE id = $1

