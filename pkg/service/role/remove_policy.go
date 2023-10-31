package role

func (r *roleService) RemovePolicy(params ...string) (bool, error) {
	return r.endforcer.RemoveNamedPolicy("p", params)
}
