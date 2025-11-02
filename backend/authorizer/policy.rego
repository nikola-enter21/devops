package authorizer

default allow := false

allow if input.rpc == "user.v1.UserService/Healthz"
allow if input.rpc == "user.v1.UserService/CheckDatabase"

allow if {
	perm := data.endpoints[input.rpc]
	role_perms := data.roles[input.role].permissions
	perm in role_perms
}
