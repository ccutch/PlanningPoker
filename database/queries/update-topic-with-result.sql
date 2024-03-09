UPDATE topics

SET status = $2,
    result = $3

WHERE id = $1
