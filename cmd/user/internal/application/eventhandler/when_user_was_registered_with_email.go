package eventhandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/vardius/go-api-boilerplate/cmd/user/internal/application/config"
	"github.com/vardius/go-api-boilerplate/cmd/user/internal/application/mailer"
	"github.com/vardius/go-api-boilerplate/cmd/user/internal/domain/user"
	"github.com/vardius/go-api-boilerplate/cmd/user/internal/infrastructure/persistence"
	"github.com/vardius/go-api-boilerplate/pkg/auth/oauth2"
	"github.com/vardius/go-api-boilerplate/pkg/domain"
	"github.com/vardius/go-api-boilerplate/pkg/errors"
	"github.com/vardius/go-api-boilerplate/pkg/eventbus"
	"github.com/vardius/go-api-boilerplate/pkg/log"
)

// WhenUserWasRegisteredWithEmail handles event
func WhenUserWasRegisteredWithEmail(db *sql.DB, repository persistence.UserRepository, tokenProvider oauth2.TokenProvider) eventbus.EventHandler {
	fn := func(parentCtx context.Context, event domain.Event) {
		ctx, cancel := context.WithTimeout(parentCtx, time.Second*120)
		defer cancel()

		logger := log.New(config.Env.App.Environment)
		logger.Info(ctx, "[EventHandler] %s\n", event.Payload)

		e := user.WasRegisteredWithEmail{}

		err := json.Unmarshal(event.Payload, &e)
		if err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", errors.Wrap(err))
			return
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", errors.Wrap(err))
			return
		}
		defer tx.Rollback()

		if err := repository.Add(ctx, userWasRegisteredWithEmailModel{e}); err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", errors.Wrap(err))
			return
		}

		if err := tx.Commit(); err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", errors.Wrap(err))
			return
		}

		token, err := tokenProvider.RetrieveToken(ctx, string(e.Email))
		if err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", errors.Wrap(err))
			return
		}

		if err := mailer.SendLoginEmail(ctx, "WhenUserWasRegisteredWithEmail", string(e.Email), token.AccessToken); err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", errors.Wrap(err))
			return
		}
	}

	return fn
}

type userWasRegisteredWithEmailModel struct {
	e user.WasRegisteredWithEmail
}

// GetID the id
func (u userWasRegisteredWithEmailModel) GetID() string {
	return u.e.ID.String()
}

// GetEmail the email
func (u userWasRegisteredWithEmailModel) GetEmail() string {
	return string(u.e.Email)
}

// GetFacebookID facebook id
func (u userWasRegisteredWithEmailModel) GetFacebookID() string {
	return ""
}

// GetGoogleID google id
func (u userWasRegisteredWithEmailModel) GetGoogleID() string {
	return ""
}
