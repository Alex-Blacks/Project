CREATE TABLE users(
	id integer generated always as identity PRIMARY KEY,
	name TEXT not null,
	email text unique not null
);

CREATE TABLE tasks(
	id integer generated always as identity PRIMARY KEY,
	user_id integer references users(id) on delete cascade,
	title text not null
);

