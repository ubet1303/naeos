package rbac

import (
	"fmt"
	"sync"
)

type Permission string

const (
	PermSpecRead    Permission = "spec:read"
	PermSpecWrite   Permission = "spec:write"
	PermSpecDelete  Permission = "spec:delete"
	PermSpecExecute Permission = "spec:execute"
	PermUserRead    Permission = "user:read"
	PermUserManage  Permission = "user:manage"
	PermTeamRead    Permission = "team:read"
	PermTeamManage  Permission = "team:manage"
	PermAuditRead   Permission = "audit:read"
	PermConfigRead  Permission = "config:read"
	PermConfigWrite Permission = "config:write"
)

type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type Store struct {
	roles map[string]*Role
	users map[string]*User
	mu    sync.RWMutex
}

func NewStore() *Store {
	s := &Store{
		roles: make(map[string]*Role),
		users: make(map[string]*User),
	}
	s.registerDefaults()
	return s
}

func (s *Store) registerDefaults() {
	s.roles["admin"] = &Role{
		Name:        "admin",
		Permissions: []Permission{PermSpecRead, PermSpecWrite, PermSpecDelete, PermSpecExecute, PermUserRead, PermUserManage, PermTeamRead, PermTeamManage, PermAuditRead, PermConfigRead, PermConfigWrite},
	}
	s.roles["developer"] = &Role{
		Name:        "developer",
		Permissions: []Permission{PermSpecRead, PermSpecWrite, PermSpecExecute, PermUserRead, PermTeamRead, PermAuditRead},
	}
	s.roles["viewer"] = &Role{
		Name:        "viewer",
		Permissions: []Permission{PermSpecRead, PermUserRead, PermTeamRead},
	}
	s.roles["operator"] = &Role{
		Name:        "operator",
		Permissions: []Permission{PermSpecRead, PermSpecExecute, PermAuditRead},
	}
}

func (s *Store) AddRole(role *Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if role.Name == "" {
		return fmt.Errorf("role name must not be empty")
	}
	s.roles[role.Name] = role
	return nil
}

func (s *Store) GetRole(name string) (*Role, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.roles[name]
	return r, ok
}

func (s *Store) AddUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user.ID == "" {
		return fmt.Errorf("user ID must not be empty")
	}
	s.users[user.ID] = user
	return nil
}

func (s *Store) GetUser(id string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

func (s *Store) HasPermission(userID string, perm Permission) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return false
	}

	for _, roleName := range user.Roles {
		role, ok := s.roles[roleName]
		if !ok {
			continue
		}
		for _, p := range role.Permissions {
			if p == perm {
				return true
			}
		}
	}
	return false
}

func (s *Store) UserPermissions(userID string) []Permission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seen := make(map[Permission]bool)
	var perms []Permission

	user, ok := s.users[userID]
	if !ok {
		return nil
	}

	for _, roleName := range user.Roles {
		role, ok := s.roles[roleName]
		if !ok {
			continue
		}
		for _, p := range role.Permissions {
			if !seen[p] {
				seen[p] = true
				perms = append(perms, p)
			}
		}
	}
	return perms
}

// Authorize checks if a user has permission and returns an error if not.
func (s *Store) Authorize(userID string, perm Permission) error {
	if !s.HasPermission(userID, perm) {
		return fmt.Errorf("user %s lacks permission %s", userID, perm)
	}
	return nil
}

// Middleware is a function that checks authorization.
type Middleware func(userID string, perm Permission) error

// NewMiddleware creates an authorization middleware from a store.
func NewMiddleware(store *Store) Middleware {
	return func(userID string, perm Permission) error {
		return store.Authorize(userID, perm)
	}
}
