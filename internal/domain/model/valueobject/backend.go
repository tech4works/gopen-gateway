package valueobject

type Backend struct {
	Host           []string
	Endpoint       string
	Method         string
	Group          string
	HideResponse   bool
	ForwardHeaders []string
	Authorizations []string
	Query          []string
	Headers        []Modifier
	Params         []Modifier
	Queries        []Modifier
	Body           []Modifier
}
