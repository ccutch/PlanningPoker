create table votes(
    topic_id    varchar(10) references topics(id),
    player_id   varchar(10) references players(id),
    choice      integer,
    created_at  timestamp default now(),
    updated_at  timestamp default now(),
    primary key (topic_id, player_id)
);
