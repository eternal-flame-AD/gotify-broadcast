package main

import (
	"sync"

	"github.com/gotify/plugin-api"
)

var usersList = new(UserPool)

// UserPool is thread-safe user pool
type UserPool struct {
	mutex sync.RWMutex
	users []plugin.UserContext
}

// AddUser adds a user context to the user pool
func (c *UserPool) AddUser(ctx plugin.UserContext) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.users = append(c.users, ctx)
}

// GetUsersList retrieves a copy of the user pool
func (c *UserPool) GetUsersList() []plugin.UserContext {
	var res []plugin.UserContext
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, u := range c.users {
		res = append(res, u)
	}
	return res
}
