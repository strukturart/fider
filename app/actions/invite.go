package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

// InviteUsers is used to invite new users into Fider
type InviteUsers struct {
	IsSampleInvite bool
	Input          *models.InviteUsers
	Invitations    []*models.UserInvitation
}

// Returns the struct to bind the request to
func (action *InviteUsers) BindTarget() interface{} {
	action.Input = new(models.InviteUsers)
	return action.Input
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *InviteUsers) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (action *InviteUsers) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Input.Subject == "" {
		result.AddFieldFailure("subject", "Subject is required.")
	} else if len(action.Input.Subject) > 70 {
		result.AddFieldFailure("subject", "Subject must have less than 70 characters.")
	}

	if action.Input.Message == "" {
		result.AddFieldFailure("message", "Message is required.")
	} else if !strings.Contains(action.Input.Message, app.InvitePlaceholder) {
		msg := fmt.Sprintf("Your message is missing the invitation link placeholder. Please add '%s' to your message.", app.InvitePlaceholder)
		result.AddFieldFailure("message", msg)
	}

	//When it's a sample invite, we skip recipients validation
	if !action.IsSampleInvite {

		if len(action.Input.Recipients) == 0 {
			result.AddFieldFailure("recipients", "At least one recipient is required.")
		} else if len(action.Input.Recipients) > 30 {
			result.AddFieldFailure("recipients", "Too many recipients. We limit at 30 recipients per invite.")
		}

		for _, email := range action.Input.Recipients {
			if email != "" {
				messages := validate.Email(email)
				result.AddFieldFailure("recipients", messages...)
			}
		}

		if result.Ok {
			action.Invitations = make([]*models.UserInvitation, 0)
			for _, email := range action.Input.Recipients {
				if email != "" {
					err := bus.Dispatch(ctx, &query.GetUserByEmail{Email: email})
					if errors.Cause(err) == app.ErrNotFound {
						action.Invitations = append(action.Invitations, &models.UserInvitation{
							Email:           email,
							VerificationKey: models.GenerateSecretKey(),
						})
					}
				}
			}

			if len(action.Invitations) == 0 {
				result.AddFieldFailure("recipients", "All these addresses have already been registered on this site.")
			}
		}

	}

	return result
}
