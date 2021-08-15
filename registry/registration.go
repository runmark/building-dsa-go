package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceUrl       string
	RequiredServices []ServiceName
	UpdateServiceURL string
	HearbeatURL      string
}

type ServiceName string

const (
	LogService           = ServiceName("LogService")
	GradingService       = ServiceName("GradingService")
	TeacherPortalService = ServiceName("TeacherPortalService")
)

type patchEntry struct {
	ServiceName ServiceName
	URL         string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
