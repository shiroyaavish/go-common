package common

// Project is the identifier for project
type Project int8

const (
	UnknownProject Project = iota
	IntelXLabs
	FDS
	AstraLink
	SisApp
)

// ParseProjectFromString takes a string representation of a project and returns the corresponding Project type.
// Possible string values are "1", "2", "3", and "4" which map to IntelXLabs, FDS, AstraLink, and SisApp respectively.
// If the provided string does not match any of the valid values, UnknownProject is returned.
func ParseProjectFromString(project string) Project {
	switch project {
	case "1":
		return IntelXLabs
	case "2":
		return FDS
	case "3":
		return AstraLink
	case "4":
		return SisApp
	default:
		return UnknownProject
	}
}

// String returns a string representation of the Project.
// It returns "IntelXLabs" if p is equal to IntelXLabs,
// "FDS" if p is equal to FDS, "AstraLink" if p is equal to AstraLink,
// "SisApp" if p is equal to SisApp, and "Unknown" otherwise.
func (p Project) String() string {
	switch p {
	case IntelXLabs:
		return "IntelXLabs"
	case FDS:
		return "FDS"
	case AstraLink:
		return "AstraLink"
	case SisApp:
		return "SisApp"
	default:
		return "Unknown"
	}
}
