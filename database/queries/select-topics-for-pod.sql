SELECT
    id,
    prompt,
    result,
    created_at,
    updated_at
FROM topics
WHERE pod_id = $1
  AND status = $2
ORDER BY created_at ASC

