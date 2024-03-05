SELECT
    id,
    prompt,
    result,
    created_at,
    updated_at
FROM topics
WHERE pod_id = $1
  AND status = 'complete'
ORDER BY updated_at DESC

