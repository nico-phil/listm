package orchestrator

type Orchestrator struct{}

func New() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Start() {
	processAllWorkspace()
}

func processAllWorkspace() error {
	// get capaigns or each worksapce
	// get active list for each compaign
	// get leads for each list
	return nil
}
