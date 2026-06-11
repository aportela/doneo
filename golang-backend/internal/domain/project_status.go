package domain

type ProjectStatus struct {
	ID       string
	Name     string
	HexColor string
	Index    uint
	Flags    Bitmask
}

type SearchProjectStatusesFilter struct {
	Name *string
}

// status flags
const (
	// this will be the default status on project creation
	ProjectStatusFlagDefaultOnCreate Bitmask = 1 << iota
	// project start date will be filled with current time (ONLY) when value is not set
	ProjectStatusFlagFillEmptyStartDate
	// project start date will be filled with current time (ALWAYS)
	ProjectStatusFlagSetStartDate
	// project finish date will be filled with current time (ONLY) when value is not set
	ProjectStatusFlagFillEmptyFinishDate
	// project finish date will be filled with current time (ALWAYS)
	ProjectStatusFlagSetFinishDate
	// project finish date will be clear when status with this flag is changed
	ProjectStatusFlagUnsetFinishDateOnLeave
)
