create table pods(
    id          varchar(10) primary key,
    name        text,
    private     boolean,
    strategy    text,
    status      text default 'waiting',
    created_at  timestamp default now(),
    updated_at  timestamp default now()
);

create table players(
    id          varchar(10) primary key,
    pod_id      varchar(10) references pods(id),
    name        text,
    owner       boolean default false,
    active      boolean default true,
    created_at  timestamp default now(),
    updated_at  timestamp default now()
);

create table topics(
    id          varchar(10) primary key,
    pod_id      varchar(10) references pods(id),
    prompt      text,
    status      text default 'upcoming',
    result      integer default 0,
    created_at  timestamp default now(),
    updated_at  timestamp default now()
);



