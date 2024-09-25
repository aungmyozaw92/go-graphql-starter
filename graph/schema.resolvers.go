package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.52

import (
	"context"

	"github.com/aungmyozaw92/go-graphql/models"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input models.NewUser) (*models.User, error) {
	return models.CreateUser(ctx, &input)
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*models.LoginInfo, error) {
	return models.Login(ctx, username, password)
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input models.NewUser) (*models.User, error) {
	return models.CreateUser(ctx, &input)
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id int, input models.NewUser) (*models.User, error) {
	return models.UpdateUser(ctx, id, &input)
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, userID int) (*models.User, error) {
	return models.DeleteUser(ctx, userID)
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (*models.User, error) {
	return models.ChangePassword(ctx, oldPassword, newPassword)
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id int) (*models.User, error) {
	return models.GetUser(ctx, id)
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, name *string, phone *string, mobile *string, email *string, isActive *bool) ([]*models.User, error) {
	return models.GetUsers(ctx, name, phone, mobile, email, isActive)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
