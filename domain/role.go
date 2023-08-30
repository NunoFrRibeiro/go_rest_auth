package domain

import (
	"strings"
)

type RolePermisisons struct {
	rolePermissions map[string][]string
}

func (rp RolePermisisons) IsAuthorizedFor(role, routeName string) bool {
	perms := rp.rolePermissions[role]
	for _, r := range perms {
		if r == strings.TrimSpace(routeName) {
			return true
		}
	}
	return false
}

func GetRolePermissions() RolePermisisons {
	return RolePermisisons{
		map[string][]string{
			"admin": {"GetAllCustomers", "GetCustomer", "NewAccount", "NewTransaction"},
			"user":  {"GetCustomer", "NewTransaction"},
		},
	}
}
