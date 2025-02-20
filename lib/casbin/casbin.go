package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Client interface {
	GetEnforcer() *casbin.Enforcer
}

type Casbin struct {
	enforcer *casbin.Enforcer
}

func (c Casbin) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

func New(db *gorm.DB) Client {
	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// create casbin RBAC model
	m := casbinModel.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "root"`)

	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	return &Casbin{
		enforcer: enforcer,
	}
}
