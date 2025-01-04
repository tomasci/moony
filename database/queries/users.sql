-- name: UsersFindMany :many
select * from users;

-- name: UsersCreate :one
insert into users (username, password, email)
values ($1, $2, $3)
returning *;
-- UsersCreate test
insert into users (username, password, email)
values ('Hello2', 'World', 'helloworld2@milkyway.galaxy')
returning *;

-- name: UsersFindUnique :one
select * from users where username = $1;
-- UsersFindUnique test
select * from users where username = 'Hello2';