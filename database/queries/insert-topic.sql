INSERT
    INTO topics (id, pod_id, prompt, status)
    VALUES ($1, $2, $3, 'upcoming')
    RETURNING id, created_at, updated_at
