package utility

import (
	"github.com/PhantomX7/go-core/lib/scope"
)

func GetUserScope(websiteID uint64) scope.Scope {
	return scope.WhereIsScope("website_id", websiteID)
}
