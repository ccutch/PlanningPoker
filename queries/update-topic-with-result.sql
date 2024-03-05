UPDATE topics

SET status = 'complete',
    result = $2

WHERE id = $1
