package scheduler

import (
	"context"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/tech-share/internal/auth"
	"github.com/mylxsw/tech-share/internal/service"
)

func syncUsers(authProvider auth.Auth, userSrv service.UserService) error {
	users, err := authProvider.Users()
	if err != nil {
		return err
	}

	for _, user := range users {
		if _, err := userSrv.LoadUser(context.TODO(), user.Account, service.UserInfo{
			Name:    user.Name,
			Account: user.Account,
			Uuid:    user.UUID,
			Status:  user.Status,
		}); err != nil {
			log.With(user).Errorf("load user from database failed: %v", err)
		}
	}

	return nil
}
