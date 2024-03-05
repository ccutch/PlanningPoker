insert into votes (topic_id, player_id, choice)
    values ($1, $2, $3)

on conflict (topic_id, player_id)
    do update set choice = $3
