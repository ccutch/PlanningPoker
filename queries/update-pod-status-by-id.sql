UPDATE pods
    SET status = $2
    WHERE id = $1

