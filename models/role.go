package models

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
)

type Role struct {
	ID          int           `gorm:"primary_key" json:"id"`
	Name        string        `gorm:"index;size:100;not null" json:"name" binding:"required"`
	RoleModules []*RoleModule `gorm:"foreignKey:RoleId"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewRole struct {
	Name           string              `json:"name" binding:"required"`
	AllowedModules []*NewAllowedModule `json:"allowed_modules"`
}

type NewAllowedModule struct {
	ModuleID       int    `json:"moduleId"`
	AllowedActions string `json:"allowedActions"`
}

func extractModuleActions(s string) []string {
	return strings.Split(strings.ToLower(s), ";")
}

// retrieve allowed query paths for role
func GetQueryPathsFromRole(ctx context.Context, roleId int) (map[string]bool, error) {
	db := config.GetDB()
	var role Role
	if err := db.WithContext(ctx).
			Preload("RoleModules").
			Preload("RoleModules.Module").
			Where("id = ?", roleId).
			First(&role).Error; err != nil {
		return nil, errors.New("role not found")
	}

	allowedPaths := make(map[string]bool, 0)
	for _, permission := range role.RoleModules {
		validActions := extractModuleActions(permission.Module.Actions)
		allowedActions := extractModuleActions(permission.AllowedActions)
		module := permission.Module.Name

		for _, action := range allowedActions {
			// check if the action is valid

			if slices.Contains(validActions, action) {
				// changing case of action & module for older module name convention
				module = utils.UppercaseFirst(module)
				switch action {
				case "read":
					allowedPaths["get"+module] = true
					allowedPaths["get"+module+"s"] = true
					allowedPaths["paginate"+module] = true
				case "update":
					allowedPaths["update"+module] = true
					allowedPaths["toggleActive"+module] = true
				default:
					action = utils.LowercaseFirst(action)
					allowedPaths[action+module] = true
				}
			}
		}
	}
	return allowedPaths, nil
}

func mapRoleModules(ctx context.Context, input []*NewAllowedModule) ([]*RoleModule, error) {

	availabeModuleActions := make(map[int]string, 0) // moduleId:actions
	modules, err := GetResources[Module](ctx)

	if err != nil {
		return nil, err
	}
	for _, m := range modules {
		availabeModuleActions[m.ID] = m.Actions
	}

	var roleModules []*RoleModule
	for _, permission := range input {

		availableActionsString, ok := availabeModuleActions[permission.ModuleID]
		if !ok || availableActionsString == "" {
			return nil, errors.New("module_id not found")
		}
		availableActions := extractModuleActions(availableActionsString)
		inputActions := extractModuleActions(permission.AllowedActions)
		for _, action := range inputActions {
			if !slices.Contains(availableActions, action) {
				return nil, errors.New("invalid module action")
			}
		}

		roleModules = append(roleModules, &RoleModule{
			ModuleId:       permission.ModuleID,
			AllowedActions: permission.AllowedActions,
		})
	}
	return roleModules, nil
}

func CreateRole(ctx context.Context, input *NewRole) (*Role, error) {

	// check duplicate
	if err := utils.ValidateUnique[Role](ctx, "name", input.Name, 0); err != nil {
		return nil, err
	}
	roleModules, err := mapRoleModules(ctx, input.AllowedModules)
	if err != nil {
		return nil, err
	}

	role := Role{
		Name:        input.Name,
		RoleModules: roleModules,
	}
	db := config.GetDB()
	// tx := db.Begin()
	err = db.WithContext(ctx).Create(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func GetRole(ctx context.Context, id int) (*Role, error) {

	return GetResource[Role](ctx, id)

}

func GetRoles(ctx context.Context, name *string) ([]*Role, error) {

	results, err := GetResources[Role](ctx, "created_at")

	if err != nil {
		return nil, err
	}

	return results, nil
}