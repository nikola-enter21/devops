package authorizer

default allow := false

allow if {
	input.rpc == "grpc.health.v1.Health/Check"
}

allow if {
	input.rpc == "grpc.health.v1.Health/Watch"
}

allow if {
	perm := data.endpoints[input.rpc]
	role_perms := data.roles[input.role].permissions
	perm in role_perms
}
