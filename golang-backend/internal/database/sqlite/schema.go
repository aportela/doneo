package sqlite

var installSchemaQueries = []string{
	`
		CREATE TABLE IF NOT EXISTS USER (
			id TEXT NOT NULL CHECK(length(id) == 36),
			email TEXT NOT NULL UNIQUE CHECK(length(email) <= 255),
			name TEXT NOT NULL CHECK(length(id) == 36),
			password_hash TEXT NOT NULL CHECK(length(password_hash) <= 60),
			ctime INTEGER NOT NULL,
			mtime INTEGER,
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		REPLACE INTO USER (id, email, name, password_hash, ctime, mtime) VALUES("019dba5d-83a4-7f97-bdf1-97a5fb3d5869", "foo@ba.r", "John Doe", "$2y$10$vskPu8bji.zmwDI9Jvi3C.uXG8IDytaeHLakb4nOzlFyoPTlwBklW", 1776948241, NULL);
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT_TYPE (
			id TEXT NOT NULL CHECK(length(id) == 36),
			name TEXT NOT NULL CHECK(length(name) <= 32),
			PRIMARY KEY (id)
		) STRICT;
	`,
	`
		REPLACE INTO PROJECT_TYPE (id, name) VALUES("019dba84-c5c4-7654-8400-d84b02164065", "Personal");
	`,
	`
		REPLACE INTO PROJECT_TYPE (id, name) VALUES("019dba85-0669-7fd4-86ed-dbe36df285af", "Work");
	`,
	`
		CREATE TABLE IF NOT EXISTS PROJECT (
			id TEXT NOT NULL CHECK(length(id) == 36),
			key TEXT NOT NULL CHECK(length(key) <= 8),
			summary TEXT NOT NULL UNIQUE CHECK(length(summary) <= 64),
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
		REPLACE INTO PROJECT (id, key, summary, description, cuser, ctime, mtime, stime, ftime, dtime, type) VALUES("019dba70-8fe5-769d-baae-6d075c5e47ee", "NEW", "New project", "Dummy description", "019dba5d-83a4-7f97-bdf1-97a5fb3d5869", "1776949459", NULL, 1776951094, NULL, NULL, "019dba84-c5c4-7654-8400-d84b02164065");
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
		REPLACE INTO PROJECT_TASK_STATUS (id, project_id, name, element_index, color) VALUES("019dba98-76a5-709d-9e3e-d56dfc5ff7e6", "019dba70-8fe5-769d-baae-6d075c5e47ee", "TODO", 0, "fc0390");
	`,
	`
		REPLACE INTO PROJECT_TASK_STATUS (id, project_id, name, element_index, color) VALUES("019dba9b-b690-7e90-9d10-6a220cea4dda", "019dba70-8fe5-769d-baae-6d075c5e47ee", "In Progress", 1, "3bcded");
	`,
	`
		REPLACE INTO PROJECT_TASK_STATUS (id, project_id, name, element_index, color) VALUES("019dba9b-ce39-7566-a849-a06b69aaaeff", "019dba70-8fe5-769d-baae-6d075c5e47ee", "Done", 2, "73f571");
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
		REPLACE INTO PROJECT_TASK_PRIORITY (id, project_id, name, element_index, color) VALUES("019dba9e-759f-7d1c-b5dc-3e3c07752d67", "019dba70-8fe5-769d-baae-6d075c5e47ee", "Low", 0, "a1c8e6");
	`,
	`
		REPLACE INTO PROJECT_TASK_PRIORITY (id, project_id, name, element_index, color) VALUES("019dba9e-759f-75ab-8404-04a58d34e16e", "019dba70-8fe5-769d-baae-6d075c5e47ee", "Medium", 1, "d6812b");
	`,
	`
		REPLACE INTO PROJECT_TASK_PRIORITY (id, project_id, name, element_index, color) VALUES("019dba9e-759f-7542-a3d8-89382ca7e71e", "019dba70-8fe5-769d-baae-6d075c5e47ee", "High", 2, "ed1a2c");
	`,
	`
		CREATE TABLE IF NOT EXISTS TASK (
			id TEXT NOT NULL CHECK(length(id) == 36),
			project_id TEXT NOT NULL CHECK(length(id) == 36),
			project_task_index INTEGER NOT NULL,
			summary TEXT NOT NULL UNIQUE CHECK(length(summary) <= 64),
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
	`
		REPLACE INTO TASK (id, project_id, project_task_index, summary, description, status, priority, cuser, ctime, mtime, stime, ftime, dtime)
		VALUES("019dbaac-207d-71d1-a2af-aea033c54931", "019dba70-8fe5-769d-baae-6d075c5e47ee", 0, "Review task insertion", "This is a test of schema insertions", "019dba9b-b690-7e90-9d10-6a220cea4dda", "019dba9e-759f-7542-a3d8-89382ca7e71e", "019dba5d-83a4-7f97-bdf1-97a5fb3d5869", 1776953444, NULL, NULL, NULL, NULL);
	`,
	`
		UPDATE TASK SET project_task_index = (SELECT COALESCE(MAX(project_task_index), 0) + 1 FROM TASK WHERE project_id = "019dba70-8fe5-769d-baae-6d075c5e47ee" AND project_task_index = 0) WHERE id = "019dbaac-207d-71d1-a2af-aea033c54931";
	`,
}
