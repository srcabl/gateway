package model

import (
	"fmt"

	"github.com/gofrs/uuid"
	postspb "github.com/srcabl/protos/posts"
	sharedpb "github.com/srcabl/protos/shared"
	sourcespb "github.com/srcabl/protos/sources"
)

func GetLinkByURLRequest(input CreatePostRequest) *postspb.GetLinkRequest {
	return &postspb.GetLinkRequest{Url: input.URL, GetBy: postspb.GetLinkRequest_URL}
}

// CreatePostRequestToPBDetermineSourceRequest converts a graphql create post request to a grpc determine post source request
func CreatePostRequestToPBDetermineSourceRequest(input CreatePostRequest) *sourcespb.DetermineLinkSourceRequest {
	return &sourcespb.DetermineLinkSourceRequest{Url: input.URL}
}

// CreatePostRequestToPBCreateLinkRequest converts a graphql create link request to a grpc create link request
func CreatePostRequestToPBCreateLinkRequest(input CreatePostRequest, sources []*sharedpb.SourceNode) *postspb.CreateLinkRequest {
	var sourcesUUIDs [][]byte
	for _, s := range sources {
		sourcesUUIDs = append(sourcesUUIDs, s.Source.Uuid)
	}
	return &postspb.CreateLinkRequest{
		Url:             input.URL,
		SourceHeadUuids: sourcesUUIDs,
	}
}

// CreatePostRequestToPBCreatePostRequest converts a graphql create post request to a grpc create post request
func CreatePostRequestToPBCreatePostRequest(input CreatePostRequest, userID, linkID []byte) *postspb.CreatePostRequest {
	return &postspb.CreatePostRequest{
		UserUuid: userID,
		LinkUuid: linkID,
		Title:    input.Title,
		Comment:  input.Comment,
	}
}

// PBCreatePostLinkResponseToCommonPostResponse converts a grpc create post response to a common post response
func PBCreatePostLinkResponseToCommonPostResponse(postRes *postspb.CreatePostResponse, link *sharedpb.Link, resErr error) *CommonPostResponse {
	fmt.Println("transforming")
	var errors []*Error
	var post *PartialPost
	if resErr != nil {
		fmt.Println("making the error")
		errors = append(errors, PBResponseErrorToError(resErr))
	}

	if postRes.Post != nil {
		fmt.Println("transoforming user")
		partuser, userErr := PBPostToPartialPost(postRes.Post, link)
		if userErr != nil {
			fmt.Println("error while transoforming user")
			errors = append(errors, userErr)
		} else {
			fmt.Println("setting transoforming user")
			post = partuser
		}
	}

	return &CommonPostResponse{
		Errors: errors,
		Post:   post,
	}
}

// PBPostToPartialPost converts a grpc post and link to a grapql partail post
func PBPostToPartialPost(post *sharedpb.Post, link *sharedpb.Link) (*PartialPost, *Error) {
	postUUID, err := uuid.FromBytes(post.Uuid)
	if err != nil {
		field := "ID"
		message := err.Error()
		return nil, &Error{
			Field:   &field,
			Message: &message,
		}
	}
	userUUID, err := uuid.FromBytes(post.UserUuid)
	if err != nil {
		field := "userID"
		message := err.Error()
		return nil, &Error{
			Field:   &field,
			Message: &message,
		}
	}
	return &PartialPost{
		ID:      postUUID.String(),
		UserID:  userUUID.String(),
		LinkURL: link.Url,
		Comment: post.Comment.PrimaryContent,
	}, nil
}
