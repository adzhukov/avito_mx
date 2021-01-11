CREATE TABLE IF NOT EXISTS tasks (
    task_id bigserial,
    seller_id bigint NOT NULL,
    file_url text NOT NULL,
    status varchar(20) NOT NULL,
    error text,
    created int,
    updated int,
    deleted int,
    invalid int
);
