package main

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/labstack/echo/v4"
)

func Authenticate(adapter *gormadapter.Adapter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {

			// ctx := e.Request().Context()

			// // Get this info from jwt token
			// user, _, _ := e.Request().BasicAuth()
			// method := e.Request().Method
			// path := e.Request().URL.Path

			// // Casbin enforces policy
			// ok, err := enforce(ctx, user, path, method, adapter)
			// if err != nil || !ok {

			// 	return &echo.HTTPError{
			// 		Code:    http.StatusForbidden,
			// 		Message: "not allowed",
			// 	}
			// }
			// if !ok {
			// 	return err
			// }
			return next(e)
		}
	}
}

func enforce(ctx context.Context, sub string, obj string, act string, adapter *gormadapter.Adapter) (bool, error) {
	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("./examples/group_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("error in policy: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("error in policy: %w", err)
	}
	return ok, nil
}

func listAllUsers() []string {
	// enforcer, err := casbin.NewEnforcer("./examples/model.conf", "./examples/policy.csv")
	adapter, _ := gormadapter.NewAdapterByDB(DB)
	enforcer, err := casbin.NewEnforcer("./examples/group_model.conf", adapter)

	userId := "1"

	if err != nil {
		fmt.Printf("Error, details: %s\n", err)
	}
	allSubjects, _ := enforcer.GetRolesForUser(userId)
	fmt.Println(allSubjects)
	return allSubjects
}

func AuthenticateUser(adapter *gormadapter.Adapter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {

			// ctx := e.Request().Context()

			// // Get this info from jwt token
			// // user, _, _ := e.Request().BasicAuth()
			// method := e.Request().Method
			// // path := e.Request().URL.Path

			// // Casbin enforces policy
			// userId := e.QueryParam("userId")
			// documentId := e.QueryParam("documentId")
			// documentCreatedBy := e.QueryParam("docCreatedBy") // => get created by from documentId

			// fmt.Println("userId", userId)
			// fmt.Println("docCreatedBy", documentCreatedBy)

			// ok, err := enforceUserDocumentRule(ctx, userId, documentCreatedBy, method, adapter)

			// if err != nil || !ok {
			// 	ok1, _ := enforceUserDocumentRule(ctx, userId, documentId, method, adapter)
			// 	if !ok1 {
			// 		return &echo.HTTPError{
			// 			Code:    http.StatusForbidden,
			// 			Message: "not allowed",
			// 		}
			// 	}
			// }
			// if !ok {
			// 	return err
			// }
			// fmt.Println("Got access")
			return next(e)
		}
	}
}

func enforceUserDocumentRule(ctx context.Context, sub string, obj string, act string, adapter *gormadapter.Adapter) (bool, error) {
	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("./examples/group_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("error in policy: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("error in policy: %w", err)
	}
	return ok, nil
}
