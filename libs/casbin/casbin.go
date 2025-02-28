package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Client interface {
	GetEnforcer() *casbin.Enforcer
}

type client struct {
	enforcer *casbin.Enforcer
}

func (c *client) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

func New(db *gorm.DB) Client {
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
		return nil
	}

	return &client{
		enforcer: enforcer,
	}
}
