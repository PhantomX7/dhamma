package casbin

import (
	"fmt"
	"strings"

	// Ensure this is imported
	"github.com/PhantomX7/dhamma/constants/permissions"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Client interface {
	GetEnforcer() *casbin.Enforcer
	AddRolePermissions(roleID uint64, domainID uint64, permissionsCodes []string) error
	DeleteRolePermissions(roleID uint64, domainID uint64, permissionsCodes []string) error
	GetRolePermissions(roleID uint64, domainID uint64) []string
	AddUserRole(userID uint64, roleID uint64, domainID uint64) error
	GetUserPermissions(userID uint64, domainID uint64) []string
}

type client struct {
	enforcer *casbin.Enforcer
}

func (c *client) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

// New creates a new Casbin client instance.
func New(db *gorm.DB) (Client, error) {
	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// create casbin RBAC model
	m := model.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act, type")
	m.AddDef("p", "p", "sub, dom, obj, act, type")
	m.AddDef("g", "g", "_, _, _")
	m.AddDef("g2", "g2", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `(g(r.sub, p.sub, r.dom) || g(r.sub, r2, r.dom) && g2(r2, p.sub)) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && keyMatch2(r.act, p.act) && r.type == p.type || r.sub == "root"`)

	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	// Load the policy from DB
	err = enforcer.LoadPolicy()
	if err != nil {
		panic(fmt.Sprintf("failed to load policy from DB: %v", err))
	}

	return &client{
		enforcer: enforcer,
	}, nil
}

func parsePermissionCode(code string) (object string, action string, permType string, err error) {
	permType = permissions.PermissionTypeApi // Default type
	codeToParse := code

	parts := strings.SplitN(code, ":", 2)
	if len(parts) == 2 {
		permType = permissions.PermissionTypeWeb
		codeToParse = parts[1]
	}

	actionParts := strings.SplitN(codeToParse, "/", 2)
	if len(actionParts) != 2 || actionParts[0] == "" || actionParts[1] == "" {
		err = fmt.Errorf("invalid permission code format: %s", code)
		return
	}
	object = actionParts[0]
	action = actionParts[1]
	return
}
