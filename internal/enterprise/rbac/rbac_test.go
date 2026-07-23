package rbac

import (
	"testing"
)

func TestStoreDefaults(t *testing.T) {
	t.Parallel()
	store := NewStore()

	admin, ok := store.GetRole("admin")
	if !ok {
		t.Fatal("admin role not found")
	}
	if len(admin.Permissions) < 5 {
		t.Errorf("admin should have many permissions, got %d", len(admin.Permissions))
	}

	viewer, ok := store.GetRole("viewer")
	if !ok {
		t.Fatal("viewer role not found")
	}
	for _, p := range viewer.Permissions {
		if p == PermSpecWrite {
			t.Error("viewer should not have spec:write")
		}
	}
}

func TestAddUserAndCheckPermission(t *testing.T) {
	t.Parallel()
	store := NewStore()

	store.AddUser(&User{ID: "alice", Name: "Alice", Roles: []string{"admin"}})

	if !store.HasPermission("alice", PermSpecWrite) {
		t.Error("admin should have spec:write")
	}
	if !store.HasPermission("alice", PermUserManage) {
		t.Error("admin should have user:manage")
	}
}

func TestUnauthorized(t *testing.T) {
	t.Parallel()
	store := NewStore()

	store.AddUser(&User{ID: "bob", Name: "Bob", Roles: []string{"viewer"}})

	if store.HasPermission("bob", PermSpecWrite) {
		t.Error("viewer should not have spec:write")
	}
	if store.HasPermission("bob", PermSpecDelete) {
		t.Error("viewer should not have spec:delete")
	}
}

func TestAuthorize(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.AddUser(&User{ID: "dev1", Name: "Dev", Roles: []string{"developer"}})

	if err := store.Authorize("dev1", PermSpecExecute); err != nil {
		t.Errorf("developer should have spec:execute: %v", err)
	}

	if err := store.Authorize("dev1", PermSpecDelete); err == nil {
		t.Error("developer should not have spec:delete")
	}

	if err := store.Authorize("unknown", PermSpecRead); err == nil {
		t.Error("unknown user should be denied")
	}
}

func TestUserPermissions(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.AddUser(&User{ID: "op1", Name: "Op", Roles: []string{"operator"}})

	perms := store.UserPermissions("op1")
	if len(perms) == 0 {
		t.Fatal("expected permissions")
	}
	found := false
	for _, p := range perms {
		if p == PermSpecExecute {
			found = true
		}
	}
	if !found {
		t.Error("operator should have spec:execute")
	}
}

func TestAddRole(t *testing.T) {
	t.Parallel()
	store := NewStore()

	store.AddRole(&Role{Name: "custom", Permissions: []Permission{PermSpecRead, PermAuditRead}})

	role, ok := store.GetRole("custom")
	if !ok {
		t.Fatal("custom role not found")
	}
	if len(role.Permissions) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(role.Permissions))
	}
}

func TestMiddleware(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.AddUser(&User{ID: "u1", Name: "U1", Roles: []string{"viewer"}})

	mw := NewMiddleware(store)

	if err := mw("u1", PermSpecRead); err != nil {
		t.Errorf("viewer should have spec:read: %v", err)
	}
	if err := mw("u1", PermSpecWrite); err == nil {
		t.Error("viewer should not have spec:write")
	}
}
