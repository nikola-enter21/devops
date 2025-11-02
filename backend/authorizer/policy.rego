package authorizer

default allow := false

allow if {
	perm := data.endpoints[input.rpc]
	role_perms := data.roles[input.role].permissions
	perm in role_perms
}
