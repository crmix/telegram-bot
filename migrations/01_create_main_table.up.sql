CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    ename TEXT NOT NULL,
    workday TIMESTAMP NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE groupid(
    id SERIAL PRIMARY KEY,
    groupchat_id bigint not NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)