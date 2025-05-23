package point_mutation

type permission struct {
	Key string
	// Index all event attendances
	Index string
	// View event attendance details
	Show string
}

var Permissions = permission{
	Key:   "point-mutation",
	Index: "index",
	Show:  "show",
}
