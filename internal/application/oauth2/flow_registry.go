package oauth2

import "fmt"

type FlowRegistry struct {
	flows map[string]AuthorizeFlow
}

func NewFlowRegistry() *FlowRegistry {
	return &FlowRegistry{flows: make(map[string]AuthorizeFlow)}
}

func (r *FlowRegistry) Register(responseType string, flow AuthorizeFlow) {
	r.flows[responseType] = flow
}

func (r *FlowRegistry) Resolve(responseType string) AuthorizeFlow {
	flow, ok := r.flows[responseType]
	if !ok {
		panic(fmt.Sprintf("no flow registered for response_type=%s", responseType))
	}
	return flow
}
