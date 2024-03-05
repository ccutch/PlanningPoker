select
    choice,
    created_at,
    updated_at
from votes

where topic_id = $1
  and player_id = $2
