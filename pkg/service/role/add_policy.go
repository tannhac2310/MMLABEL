package role

func (r *roleService) AddPolicy(params ...string) (bool, error) {
	return r.endforcer.AddNamedPolicy("p", params)
}
