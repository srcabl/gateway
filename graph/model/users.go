package model

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/srcabl/gateway/internal/util"
	sharedpb "github.com/srcabl/protos/shared"
	"github.com/srcabl/protos/users"
	userspb "github.com/srcabl/protos/users"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserRequestToPBCreateUserRequest converts a graphql register user request to a grpc create user request
func ValidateRegisterAndRegisterUserRequestToPBCreateUserRequest(input RegisterUserRequest) (*userspb.CreateUserRequest, []*Error) {
	if isvalid, message := util.ValidateUsernameRequirements(input.Username); !isvalid {
		field := "username"
		return nil, []*Error{{Field: &field, Message: &message}}
	}
	if isvalid, message := util.ValidateMinPasswordRequirements(input.Password); !isvalid {
		field := "password"
		return nil, []*Error{{Field: &field, Message: &message}}
	}
	if isvalid, message := util.ValidateEmailRequirements(input.Email); !isvalid {
		field := "email"
		return nil, []*Error{{Field: &field, Message: &message}}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		field := "password"
		message := "password is not valid"
		return nil, []*Error{{Field: &field, Message: &message}}
	}
	return &userspb.CreateUserRequest{
		Username:        input.Username,
		Email:           input.Email,
		HashedPasssword: string(hash),
		//TODO display and description
	}, nil
}

// CurrentUserRequestToPBGetUserRequest converts a graphql current user request to a grpc get user request
func CurrentUserRequestToPBGetUserRequest(userUUID []byte) *userspb.GetUserRequest {
	return &users.GetUserRequest{
		Uuid: userUUID,
	}
}

// LoginUserRequestToPBValidateUserCredentials converts a graphql login request to a grpc validate user request
func LoginUserRequestToPBValidateUserCredentials(input LoginUserRequest) (*userspb.ValidateUserCredentialsRequest, []*Error) {
	userName, email, validateBy := determineValidateByFields(input.UsernameOrEmail)
	return &userspb.ValidateUserCredentialsRequest{
		Username:       userName,
		Email:          email,
		ValidateUserBy: validateBy,
		Password:       input.Password,
	}, nil
}

func ToPBFollowRequest(followerUUID, followedUUID []byte, followType userspb.FollowRequest_FollowedType) *userspb.FollowRequest {
	return &userspb.FollowRequest{
		FollowerUuid: followerUUID,
		FollowedUuid: followedUUID,
		Type:         followType,
	}
}

func determineValidateByFields(usernameOrEmail string) (string, string, userspb.ValidateUserCredentialsRequest_ValidateUserBy) {
	if isvalid, _ := util.ValidateEmailRequirements(usernameOrEmail); isvalid {
		return "", usernameOrEmail, users.ValidateUserCredentialsRequest_EMAIL
	}
	return usernameOrEmail, "", users.ValidateUserCredentialsRequest_USERNAME
}

// PBCreateUserResponseToCommonUserResponse converts a grpc get user response to a graphql common useer response
func PBCreateUserResponseToCommonUserResponse(res *userspb.CreateUserResponse, resErr error) *CommonUserResponse {
	return userResponseToCommonUserResponse(res, resErr)
}

// PBGetUserResponseToCommonUserResponse converts a grpc get user response to a graphql common useer response
func PBGetUserResponseToCommonUserResponse(res *userspb.GetUserResponse, resErr error) *CommonUserResponse {
	return userResponseToCommonUserResponse(res, resErr)
}

// PBValidateUserResponseToCommonUserResponse converts a grpc get user response to a graphql common useer response
func PBValidateUserResponseToCommonUserResponse(res *userspb.ValidateUserCredentialsResponse, resErr error) *CommonUserResponse {
	return userResponseToCommonUserResponse(res, resErr)
}

type userGetter interface {
	GetUser() *sharedpb.User
}

func userResponseToCommonUserResponse(ug userGetter, resErr error) *CommonUserResponse {
	fmt.Println("transforming")
	var errors []*Error
	var user *PartialUser
	if resErr != nil {
		fmt.Println("making the error")
		errors = append(errors, PBResponseErrorToError(resErr))
	}

	if ug.GetUser() != nil {
		fmt.Println("transoforming user")
		partuser, userErr := PBUserToPartialUser(ug.GetUser())
		if userErr != nil {
			fmt.Println("error while transoforming user")
			errors = append(errors, userErr)
		} else {
			fmt.Println("setting transoforming user")
			user = partuser
		}
	}

	return &CommonUserResponse{
		Errors: errors,
		User:   user,
	}
}

// PBUserToPartialUser converts a grpc user to a graphql partial user
func PBUserToPartialUser(user *sharedpb.User) (*PartialUser, *Error) {
	uuid, err := uuid.FromBytes(user.Uuid)
	if err != nil {
		field := "ID"
		message := err.Error()
		return nil, &Error{
			Field:   &field,
			Message: &message,
		}
	}
	return &PartialUser{
		ID:       uuid.String(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// PBResponseErrorToError converts a grpc error to a graphql error
func PBResponseErrorToError(err error) *Error {
	// TODO figure out field
	field := ""
	message := err.Error()
	return &Error{
		Field:   &field,
		Message: &message,
	}
}
