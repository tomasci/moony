-- name: UsersFindMany :many
select *
from users;

-- name: UsersCreate :one
insert into users (username, password, email)
values ($1, $2, $3)
returning *;
-- test
insert into users (username, password, email)
values ('Hello', 'World', 'helloworld@milkyway.galaxy')
returning *;

