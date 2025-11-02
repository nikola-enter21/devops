package policy

default allow := false

allow if {
	some allowed_route in data.roles[input.role]
	input.route == allowed_route
}
