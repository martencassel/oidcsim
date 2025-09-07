package memory

import "sync"

type inMemoryAuthorizationCodeRepo struct {
	codes map[string]string
	mu    sync.RWMutex
}

func NewInMemoryAuthorizationCodeRepo() *inMemoryAuthorizationCodeRepo {
	return &inMemoryAuthorizationCodeRepo{
		codes: make(map[string]string),
		mu:    sync.RWMutex{},
	}
}

func (r *inMemoryAuthorizationCodeRepo) Save(code string, data string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.codes[code] = data
	return nil
}

func (r *inMemoryAuthorizationCodeRepo) Get(code string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	data, exists := r.codes[code]
	return data, exists
}

func (r *inMemoryAuthorizationCodeRepo) Delete(code string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.codes, code)
}
