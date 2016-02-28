package server

import "sync"

type ConnectedAgents struct {
	m   map[string]*agent
	mtx *sync.Mutex
}

func (cg *ConnectedAgents) get(id string) *agent {
	cg.mtx.Lock()
	defer cg.mtx.Unlock()

	return cg.m[id]
}

func (cg *ConnectedAgents) put(id string, agent *agent) {
	cg.mtx.Lock()
	defer cg.mtx.Unlock()

	cg.m[id] = agent
}

func (cg *ConnectedAgents) delete(id string) {
	cg.mtx.Lock()
	defer cg.mtx.Unlock()

	delete(cg.m, id)
}

var connectedAgents ConnectedAgents = ConnectedAgents{
	m:   make(map[string]*agent),
	mtx: new(sync.Mutex),
}
