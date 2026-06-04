package taskrepository

// insert task query
/*
	UPDATE project_task_counter SET
		next_task_index = next_task_index + 1
	WHERE project_id = ?
	RETURNING next_task_index - 1;

	INSERT INTO tasks (id, project_id, task_index, summary, description, creator_id, created_at, updated_at, deleted_at, archived_at, started_at, finished_at, due_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
*/
