package sqlite

type schemaMigration struct {
	Version uint
	Queries []string
}

var schemaQueries = []schemaMigration{
	{
		Version: 1,
		Queries: []string{
			`
				CREATE TABLE IF NOT EXISTS users (
					id TEXT NOT NULL CHECK(length(id) == 36),
					email TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK(length(email) <= 255),
					name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 32),
					password_hash TEXT NOT NULL CHECK(length(password_hash) <= 60),
					created_at INTEGER NOT NULL,
					updated_at INTEGER,
					deleted_at INTEGER,
					permissions_bitmask INTEGER NOT NULL DEFAULT 0 CHECK(permissions_bitmask >= 0),
					PRIMARY KEY (id)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS roles (
					id TEXT NOT NULL CHECK(length(id) == 36),
					name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 32),
					permissions_bitmask INTEGER NOT NULL DEFAULT 0 CHECK(permissions_bitmask >= 0),
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
				CREATE TABLE IF NOT EXISTS project_priorities (
					id TEXT NOT NULL CHECK(length(id) == 36),
					name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 16),
					item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
					item_index INTEGER NOT NULL UNIQUE,
					PRIMARY KEY (id)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS project_statuses (
					id TEXT NOT NULL CHECK(length(id) == 36),
					name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 16),
					item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
					item_index INTEGER NOT NULL UNIQUE,
					flags_bitmask INTEGER NOT NULL DEFAULT 0,
					PRIMARY KEY (id)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS task_priorities (
					id TEXT NOT NULL CHECK(length(id) == 36),
					name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 16),
					item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
					item_index INTEGER NOT NULL UNIQUE,
					PRIMARY KEY (id)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS task_statuses (
					id TEXT NOT NULL CHECK(length(id) == 36),
					name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 1 AND 16),
					item_hex_color TEXT NOT NULL CHECK(length(item_hex_color) = 7),
					item_index INTEGER NOT NULL UNIQUE,
					flags_bitmask INTEGER NOT NULL DEFAULT 0,
					PRIMARY KEY (id)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS projects (
					id TEXT NOT NULL CHECK(length(id) == 36),
					slug TEXT NOT NULL UNIQUE CHECK(length(slug) BETWEEN 1 AND 8),
					summary TEXT NOT NULL CHECK(length(summary) BETWEEN 1 AND 128),
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
					id TEXT NOT NULL CHECK(length(id) == 36),
					project_id TEXT NOT NULL CHECK(length(project_id) == 36),
					user_id TEXT NOT NULL CHECK(length(user_id) == 36),
					role_id TEXT NOT NULL CHECK(length(role_id) == 36),
					PRIMARY KEY (id),
					UNIQUE(project_id,user_id),
					FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
					FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
					FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_project_user_role_project_id ON project_user_role(project_id);
				CREATE INDEX IF NOT EXISTS idx_project_user_role_user_id ON project_user_role(user_id);
				CREATE INDEX IF NOT EXISTS idx_project_user_role_role_id ON project_user_role(role_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS project_notes (
					id TEXT NOT NULL CHECK(length(id) == 36),
					project_id TEXT NOT NULL CHECK(length(project_id) == 36),
					creator_id TEXT NOT NULL CHECK(length(creator_id) == 36),
					created_at INTEGER NOT NULL,
					updated_at INTEGER,
					body TEXT NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
					FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_project_notes_project_id ON project_notes(project_id);
				CREATE INDEX IF NOT EXISTS idx_project_notes_creator_id ON project_notes(creator_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS attachments (
					id TEXT NOT NULL CHECK(length(id) == 36),
					original_name TEXT NOT NULL,
					content_type TEXT NOT NULL,
					size INTEGER NOT NULL,
					user_id TEXT NOT NULL CHECK(length(user_id) == 36),
					created_at INTEGER NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_attachments_user_id ON attachments(user_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS project_attachments (
					project_id TEXT NOT NULL CHECK(length(project_id) == 36),
					attachment_id TEXT NOT NULL CHECK(length(attachment_id) == 36),
					PRIMARY KEY (project_id, attachment_id),
					FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
					FOREIGN KEY(attachment_id) REFERENCES attachments(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_project_attachments_project_id ON project_attachments(project_id);
				CREATE INDEX IF NOT EXISTS idx_project_attachments_attachment_id ON project_attachments(attachment_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS history_operations (
					id TEXT NOT NULL CHECK(length(id) == 36),
					project_id TEXT NOT NULL CHECK(length(project_id) == 36),
					task_id TEXT CHECK(length(task_id) == 36),
					operation_type INTEGER NOT NULL,
					operation_user_id TEXT NOT NULL CHECK(length(operation_user_id) == 36),
					operation_date INTEGER NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
					FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
					FOREIGN KEY(operation_user_id) REFERENCES users(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_history_operations_project_id ON history_operations(project_id);
				CREATE INDEX IF NOT EXISTS idx_history_operations_task_id ON history_operations(task_id);
				CREATE INDEX IF NOT EXISTS idx_history_operations_operation_user_id ON history_operations(operation_user_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS project_task_counter (
					project_id TEXT NOT NULL CHECK(length(project_id) == 36),
					next_task_index INTEGER NOT NULL DEFAULT 1 CHECK(next_task_index > 0),
					PRIMARY KEY (project_id),
					FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_project_task_counter_project_id ON project_task_counter(project_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS tasks (
					id TEXT NOT NULL CHECK(length(id) == 36),
					project_id TEXT NOT NULL CHECK(length(project_id) == 36),
					task_index INTEGER NOT NULL,
					summary TEXT NOT NULL CHECK(length(summary) BETWEEN 1 AND 128),
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
					cover_attachment_id TEXT CHECK(length(cover_attachment_id) == 36),
					PRIMARY KEY (id),
					FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE,
					FOREIGN KEY(priority_id) REFERENCES task_priorities(id) ON DELETE CASCADE,
					FOREIGN KEY(status_id) REFERENCES task_statuses(id) ON DELETE CASCADE,
					UNIQUE(project_id, task_index)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS task_tags (
					task_id TEXT NOT NULL CHECK(length(task_id) == 36),
					tag TEXT NOT NULL CHECK(length(tag) <= 64),
					PRIMARY KEY (task_id, tag),
					FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS timers (
					id TEXT NOT NULL CHECK(length(id) == 36),
					user_id TEXT NOT NULL CHECK(length(user_id) == 36),
					summary TEXT NOT NULL CHECK(length(summary) BETWEEN 1 AND 32),
					started_at INTEGER NOT NULL,
					finished_at INTEGER,
					PRIMARY KEY (id),
					FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
					UNIQUE(user_id, started_at)
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS task_relations (
					task_id TEXT NOT NULL CHECK(length(task_id) == 36),
					related_task_id  TEXT NOT NULL CHECK(length(task_id) == 36),
					relation_type INTEGER NOT NULL,
					PRIMARY KEY (task_id, related_task_id),
					FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
					FOREIGN KEY(related_task_id) REFERENCES tasks(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE TABLE IF NOT EXISTS task_notes (
					id TEXT NOT NULL CHECK(length(id) == 36),
					task_id TEXT NOT NULL CHECK(length(task_id) == 36),
					creator_id TEXT NOT NULL CHECK(length(creator_id) == 36),
					created_at INTEGER NOT NULL,
					updated_at INTEGER,
					body TEXT NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
					FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_task_notes_task_id ON task_notes(task_id);
				CREATE INDEX IF NOT EXISTS idx_task_notes_creator_id ON task_notes(creator_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS task_attachments (
					task_id TEXT NOT NULL CHECK(length(task_id) == 36),
					attachment_id TEXT NOT NULL CHECK(length(attachment_id) == 36),
					PRIMARY KEY (task_id, attachment_id),
					FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
					FOREIGN KEY(attachment_id) REFERENCES attachments(id) ON DELETE CASCADE
				) STRICT;
			`,
			`
				CREATE INDEX IF NOT EXISTS idx_task_attachments_task_id ON task_attachments(task_id);
				CREATE INDEX IF NOT EXISTS idx_task_attachments_attachment_id ON task_attachments(attachment_id);
			`,
			`
				CREATE TABLE IF NOT EXISTS task_time_entries (
					id TEXT NOT NULL CHECK(length(id) == 36),
					task_id TEXT NOT NULL CHECK(length(task_id) == 36),
					user_id TEXT NOT NULL CHECK(length(user_id) == 36),
					created_at INTEGER NOT NULL,
					summary TEXT NOT NULL CHECK(length(summary) BETWEEN 1 AND 128),
					total_seconds INTEGER NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
					FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
				) STRICT;
			`,
		},
	},
}
