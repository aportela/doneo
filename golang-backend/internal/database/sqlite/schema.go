package sqlite

var installSchemaQueries = []string{
	`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL CHECK(length(id) == 36),
			email TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK(length(email) <= 255),
			name TEXT NOT NULL CHECK(length(name) BETWEEN 1 AND 32),
			password_hash TEXT NOT NULL CHECK(length(password_hash) <= 60),
			created_at INTEGER NOT NULL,
			updated_at INTEGER,
			deleted_at INTEGER,
			is_super_user INTEGER NOT NULL DEFAULT 0 CHECK(is_super_user IN (0, 1)),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS roles (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 32),
			permission_bitmask INTEGER NOT NULL DEFAULT 0 CHECK(permission_bitmask >= 0),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_priorities(
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 32),
			item_index INTEGER NOT NULL,
			item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_statuses (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 32),
			item_index INTEGER NOT NULL,
			item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS project_types (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 32),
			item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT NOT NULL CHECK(length(id) == 36),
			key TEXT NOT NULL UNIQUE CHECK(length(key) BETWEEN 1 AND 8),
			summary TEXT NOT NULL CHECK(length(key) BETWEEN 1 AND 128),
			description TEXT,
			creator_id TEXT NOT NULL CHECK(length(creator_id) == 36),
			created_at INTEGER NOT NULL,
			updated_at INTEGER,
			deleted_at INTEGER,
			archived_at INTEGER,
			started_at INTEGER,
			finished_at INTEGER,
			due_at INTEGER,
			priority_id TEXT NOT NULL CHECK(length(priority_id) == 36),
			status_id TEXT NOT NULL CHECK(length(status_id) == 36),
			type_id TEXT NOT NULL CHECK(length(type_id) == 36),
			PRIMARY KEY (id),
			FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(priority_id) REFERENCES project_priorities(id) ON DELETE CASCADE,
			FOREIGN KEY(status_id) REFERENCES project_statuses(id) ON DELETE CASCADE,
			FOREIGN KEY(type_id) REFERENCES project_types(id) ON DELETE CASCADE
		) STRICT;
	`,
	`
		CREATE INDEX idx_projects_status_id ON projects(status_id);
		CREATE INDEX idx_projects_priority_id ON projects(priority_id);
		CREATE INDEX idx_projects_type_id ON projects(type_id);
		CREATE INDEX idx_projects_creator_id ON projects(creator_id);
	`,
	`
		CREATE TABLE IF NOT EXISTS project_user_role (
			project_id TEXT NOT NULL CHECK(length(project_id) == 36),
			user_id TEXT NOT NULL CHECK(length(user_id) == 36),
			role_id TEXT NOT NULL CHECK(length(role_id) == 36),
			PRIMARY KEY (project_id, user_id),
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
		) STRICT;
	`,
	/*
		`
		CREATE INDEX IF NOT EXISTS idx_project_user_roles_project_id ON project_user_roles(project_id);
		CREATE INDEX IF NOT EXISTS idx_project_user_roles_user_id ON project_user_roles(user_id);
		CREATE INDEX IF NOT EXISTS idx_project_user_roles_role_id ON project_user_roles(role_id);
		`,
	*/
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
