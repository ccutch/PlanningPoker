INSERT
    INTO pods (id, name, strategy, private)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at
