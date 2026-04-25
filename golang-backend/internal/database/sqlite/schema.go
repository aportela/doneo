package sqlite

var installSchemaQueries = []string{
	`
		CREATE TABLE IF NOT EXISTS USER (
			id TEXT NOT NULL CHECK(length(id) == 36),
			email TEXT NOT NULL UNIQUE CHECK(length(email) <= 255),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			password_hash TEXT NOT NULL CHECK(length(password_hash) <= 60),
			ctime INTEGER NOT NULL,
			mtime INTEGER,
			flag_is_admin INTEGER NOT NULL DEFAULT 0 CHECK(flag_is_admin IN (0, 1)),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT_TYPE (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT (
			id TEXT NOT NULL CHECK(length(id) == 36),
			key TEXT NOT NULL CHECK(length(key) <= 8),
			summary TEXT NOT NULL CHECK(length(summary) <= 128),
			description TEXT,
			cuser TEXT NOT NULL CHECK(length(id) == 36),
			ctime INTEGER NOT NULL,
			mtime INTEGER,
			stime INTEGER,
			ftime INTEGER,
			dtime INTEGER,
			type TEXT NOT NULL CHECK(length(id) == 36),
			PRIMARY KEY (id),
			FOREIGN KEY(cuser) REFERENCES USER(id) ON DELETE CASCADE,
			FOREIGN KEY(type) REFERENCES PROJECT_TYPE(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT_PARTICIPANT (
			project_id TEXT NOT NULL CHECK(length(project_id) == 36),
			user_id TEXT NOT NULL CHECK(length(user_id) == 36),
			PRIMARY KEY (project_id, user_id),
			FOREIGN KEY(project_id) REFERENCES PROJECT(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES USER(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT_TASK_STATUS (
			id TEXT NOT NULL CHECK(length(id) == 36),
			project_id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			color TEXT NOT NULL CHECK(length(color) == 6),
			element_index INTEGER NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY(project_id) REFERENCES PROJECT(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT_TASK_PRIORITY (
			id TEXT NOT NULL CHECK(length(id) == 36),
			project_id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			color TEXT NOT NULL CHECK(length(color) == 6),
			element_index INTEGER NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY(project_id) REFERENCES PROJECT(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS TASK (
			id TEXT NOT NULL CHECK(length(id) == 36),
			project_id TEXT NOT NULL CHECK(length(id) == 36),
			project_task_index INTEGER NOT NULL,
			summary TEXT NOT NULL UNIQUE CHECK(length(summary) <= 128),
			description TEXT,
			status TEXT NOT NULL CHECK(length(id) == 36),
			priority TEXT NOT NULL CHECK(length(id) == 36),
			cuser TEXT NOT NULL CHECK(length(id) == 36),
			ctime INTEGER NOT NULL,
			mtime INTEGER,
			stime INTEGER,
			ftime INTEGER,
			dtime INTEGER,
			PRIMARY KEY (id),
			FOREIGN KEY(project_id) REFERENCES PROJECT(id) ON DELETE CASCADE,
			FOREIGN KEY(status) REFERENCES PROJECT_TASK_STATUS(id) ON DELETE CASCADE,
			FOREIGN KEY(priority) REFERENCES PROJECT_TASK_PRIORITY(id) ON DELETE CASCADE,
			FOREIGN KEY(cuser) REFERENCES USER(id) ON DELETE CASCADE
		) STRICT;
	`,
}
