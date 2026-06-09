package domain

type TaskStatus struct {
	ID       string
	Name     string
	HexColor string
	Index    uint
	Flags    Bitmask
}

type SearchTaskStatusesFilter struct {
	Name *string
}

// status flags
const (
	// this will be the default status on task creation
	TaskStatusFlagDefaultOnCreate Bitmask = 1 << iota
	// task start date will be filled with current time (ONLY) when value is not set
	TaskStatusFlagFillEmptyStartDate
	// task start date will be filled with current time (ALWAYS)
	TaskStatusFlagSetStartDate
	// task finish date will be filled with current time (ONLY) when value is not set
	TaskStatusFlagFillEmptyFinishDate
	// task finish date will be filled with current time (ALWAYS)
	TaskStatusFlagSetFinishDate
)
