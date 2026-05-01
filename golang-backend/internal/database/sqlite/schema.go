package sqlite

var installSchemaQueries = []string{
	// TODO: role table ?
	`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL CHECK(length(id) == 36),
			email TEXT NOT NULL UNIQUE CHECK(length(email) <= 255),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			password_hash TEXT NOT NULL CHECK(length(password_hash) <= 60),
			created_at INTEGER NOT NULL,
			updated_at INTEGER,
			is_super_user INTEGER NOT NULL DEFAULT 0 CHECK(is_super_user IN (0, 1)),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS roles (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			permission_bitmask INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 16),
			description TEXT,
			creator_id TEXT NOT NULL CHECK(length(creator_id) == 36),
			created_at INTEGER NOT NULL,
			updated_at INTEGER,
			PRIMARY KEY (id),
			FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS workspace_user_role (
			workspace_id TEXT NOT NULL CHECK(length(workspace_id) == 36),
			user_id TEXT NOT NULL CHECK(length(user_id) == 36),
			role_id TEXT NOT NULL CHECK(length(role_id) == 36),
			PRIMARY KEY (workspace_id, user_id),
			FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_priorities(
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			item_index INTEGER NOT NULL DEFAULT 0,
			item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) == 7),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_statuses (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			item_index INTEGER NOT NULL DEFAULT 0,
			item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) == 7),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_types (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT NOT NULL CHECK(length(id) == 36),
			workspace_id TEXT NOT NULL CHECK(length(workspace_id) == 36),
			key TEXT NOT NULL CHECK(length(key) <= 8),
			summary TEXT NOT NULL CHECK(length(summary) <= 128),
			description TEXT,
			creator_id TEXT NOT NULL CHECK(length(creator_id) == 36),
			created_at INTEGER NOT NULL,
			updated_at INTEGER,
			started_at INTEGER,
			finished_at INTEGER,
			due_at INTEGER,
			priority_id TEXT NOT NULL CHECK(length(priority_id) == 36),
			status_id TEXT NOT NULL CHECK(length(status_id) == 36),
			type_id TEXT NOT NULL CHECK(length(type_id) == 36),
			PRIMARY KEY (id),
			FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
			FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(priority_id) REFERENCES project_priorities(id) ON DELETE CASCADE,
			FOREIGN KEY(status_id) REFERENCES project_statuses(id) ON DELETE CASCADE,
			FOREIGN KEY(type_id) REFERENCES project_types(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_user_role (
			project_id TEXT NOT NULL CHECK(length(project_id) == 36),
			user_id TEXT NOT NULL CHECK(length(user_id) == 36),
			role_id TEXT NOT NULL CHECK(length(role_id) == 36),
			PRIMARY KEY (project_id, user_id),
			FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
		) STRICT;
	`,
	/*
		`
			CREATE TABLE IF NOT EXISTS project_task_status (
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
				creator_id TEXT NOT NULL CHECK(length(id) == 36),
				created_at INTEGER NOT NULL,
				updated_at INTEGER,
				started_at INTEGER,
				finished_at INTEGER,
				due_at INTEGER,
				PRIMARY KEY (id),
				FOREIGN KEY(project_id) REFERENCES PROJECT(id) ON DELETE CASCADE,
				FOREIGN KEY(status) REFERENCES PROJECT_TASK_STATUS(id) ON DELETE CASCADE,
				FOREIGN KEY(priority) REFERENCES PROJECT_TASK_PRIORITY(id) ON DELETE CASCADE,
				FOREIGN KEY(creator_id) REFERENCES USER(id) ON DELETE CASCADE
			) STRICT;
		`,
	*/
}
