select
    player_id,
    choice,
    created_at,
    updated_at
from votes

where topic_id = $1
