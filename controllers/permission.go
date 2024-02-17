package controllers

import (
	"fmt"
	"strconv"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func AddRoleForUser(subject string, group string) {
	enforcer := getEnforcer()
	enforcer.AddRoleForUser(subject, group)
}

func AddPolicy(subject string, resource string, permission string) {
	enforcer := getEnforcer()
	enforcer.AddPolicy(subject, resource, permission)
}

func FindAccess(subject string) []int {
	var result []int
	enforcer := getEnforcer()
	resouces, _ := enforcer.GetImplicitResourcesForUser(subject)

	for _, resource := range resouces {
		for index, permission := range resource {
			if index == 1 {
				id, _ := strconv.Atoi(permission)
				result = append(result, id)
			}

		}
	}
	fmt.Println("result", result)
	return result

}
func GetAllGroups() []string {
	enforcer := getEnforcer()
	return enforcer.GetAllRoles()
}
func GetAllRolesForUser(subject string) []string {
	enforcer := getEnforcer()
	roles, _ := enforcer.GetImplicitRolesForUser(subject)

	return roles
}
func GetAllUsersForRole(subject string) []string {
	enforcer := getEnforcer()
	users, _ := enforcer.GetImplicitUsersForRole(subject)

	return users
}

func enforce(sub string, resource string, act string) (bool, error) {
	// Load model configuration file and policy store adapter
	enforcer := getEnforcer()
	// Load policies from DB dynamically
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("error in policy: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub, resource, act)
	if err != nil {
		return false, fmt.Errorf("error in policy: %w", err)
	}
	return ok, nil
}

func getEnforcer() *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		fmt.Printf("Error, details: %s\n", err)
	}
	enforcer, enforceError := casbin.NewEnforcer("./examples/group_model.conf", adapter)
	if enforceError != nil {
		fmt.Printf("Enforce Error, details: %s\n", err)
	}
	return enforcer
}
