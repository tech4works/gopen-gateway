package vo

type BackendDependenciesConfig struct {
	ids  []string
	idxs []int
}

func NewBackendDependenciesConfig(ids []string, idxs []int) *BackendDependenciesConfig {
	return &BackendDependenciesConfig{
		ids:  ids,
		idxs: idxs,
	}
}

func (d *BackendDependenciesConfig) IDs() []string {
	return d.ids
}

func (d *BackendDependenciesConfig) Indexes() []int {
	return d.idxs
}
