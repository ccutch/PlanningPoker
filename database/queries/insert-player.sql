INSERT
    INTO players (id, pod_id, name, owner)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at
