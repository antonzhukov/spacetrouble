CREATE TABLE bookings
(
    id bigserial primary key,
    first_name  character varying(255) DEFAULT '',
    last_name  character varying(255) DEFAULT '',
    gender  character varying(255) DEFAULT '',
    birthday timestamp default NULL,
    launchpad_id  character varying(255) DEFAULT '',
    destination_id  character varying(255) DEFAULT '',
    launch_date timestamp default NULL
);
