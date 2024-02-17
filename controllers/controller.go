package controllers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-casbin/database"
	"github.com/go-casbin/model"
	"github.com/labstack/echo/v4"
)

const publicGroup = "public_group"

func CreateUser(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	database.Instance.Create(&u)
	return c.JSON(http.StatusOK, u)
}

func ListUsersGroup(c echo.Context) error {
	userId := c.QueryParam("userId")
	var user model.User
	database.Instance.First(&user, userId)
	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, GetAllRolesForUser(userId))
}
func ListGroupUsers(c echo.Context) error {
	group := c.QueryParam("group")
	allgroups := GetAllGroups()
	if !slices.Contains(allgroups, group) {
		return c.JSON(http.StatusNotFound, "Group not found")
	}

	return c.JSON(http.StatusOK, GetAllUsersForRole(group))
}

func CreateDocuent(c echo.Context) error {
	doc := new(model.Document)
	if err := c.Bind(doc); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	database.Instance.Create(&doc)
	return c.JSON(http.StatusOK, doc)
}

func CreateGroup(c echo.Context) error {
	var payload struct {
		UserId    string
		GroupName string
	}
	fmt.Println("payload", payload)
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
	handleError(err)

	AddRoleForUser(payload.UserId, payload.GroupName)
	return c.JSON(http.StatusCreated, payload)
}

func AddPermission(c echo.Context) error {
	var payload struct {
		DocumentId string
		GroupName  string
		UserId     string
		Permission string
	}

	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
	handleError(err)

	if len(payload.GroupName) > 0 {
		AddPolicy(payload.GroupName, payload.DocumentId, payload.Permission)
	} else if len(payload.UserId) > 0 {
		AddPolicy(payload.UserId, payload.DocumentId, payload.Permission)
	} else {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	return c.JSON(http.StatusCreated, "Permission added")
}

func GetAllDocumentsByUser(c echo.Context) error {
	var documents []model.Document
	userId := c.QueryParam("userId")
	documentids := FindAccess(userId)
	publicDocuments := FindAccess(publicGroup)
	documentids = append(documentids, publicDocuments...)
	database.Instance.Find(&documents, documentids)
	return c.JSON(http.StatusOK, documents)
}

func GetDocumentById(c echo.Context) error {
	var document model.Document
	documentId := c.QueryParam("documentId")
	userId := c.QueryParam("userId")
	access := c.QueryParam("access")

	ok, _ := enforce(userId, documentId, access)
	public_access, _ := enforce(publicGroup, documentId, "*")

	fmt.Println("public_access", public_access)
	// fmt.Println("public_access", publicDoc)
	if !ok && !public_access {
		fmt.Println("No access")
		return c.JSON(http.StatusForbidden, "Access Denied")
	}
	fmt.Println("Got access")
	database.Instance.First(&document, documentId)
	if document.ID != 0 {
		return c.JSON(http.StatusOK, document)
	}
	return c.JSON(http.StatusNotFound, "Document not found")
}

func MakeDocumentsPublic(c echo.Context) error {
	var documents []model.Document
	var payload struct {
		DocIds []int
		Groups []string
	}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
	handleError(err)
	database.Instance.Find(&documents, payload.DocIds)
	for _, doc := range documents {
		docId := strconv.FormatUint(uint64(doc.ID), 10)
		AddPolicy(publicGroup, docId, "*")
	}
	if len(payload.Groups) > 0 {
		groups := GetAllGroups()
		for _, group := range payload.Groups {
			if slices.Contains(groups, group) {
				AddPolicy(publicGroup, group, "*")
			}

		}
	}

	return c.JSON(http.StatusOK, documents)
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error has occoured", err)
		panic("Error has occoured")
	}
}
