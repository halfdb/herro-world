package authorization

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
)

const (
	keyCid      = "cid"
	keyChat     = "chat"
	keyUserChat = "user_chat"
)

func AuthorizeChatMember(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := GetUid(c)
		cid := 0
		err := echo.PathParamsBinder(c).Int("cid", &cid).BindError()
		if err != nil || cid == 0 {
			return echo.ErrBadRequest
		}
		c.Set(keyCid, cid)

		userChatCh := make(chan *models.UserChat, 1)
		errCh := make(chan error, 2) // shared error channel
		go func() {
			userChat, err := dao.FetchUserChat(uid, cid, true)
			if err == sql.ErrNoRows {
				userChatCh <- nil // use nil to imply no result
			} else if err != nil {
				errCh <- err
				close(userChatCh)
			} else {
				userChatCh <- userChat
				c.Set(keyUserChat, userChat)
			}
		}()

		chatCh := make(chan *models.Chat, 1)
		go func() {
			chat, err := dao.FetchChat(cid, true)
			if err != nil {
				errCh <- err
				close(chatCh)
			} else {
				chatCh <- chat
				c.Set(keyChat, chat)
			}
		}()

		userChat, ok := <-userChatCh
		if !ok {
			return <-errCh
		}
		if userChat == nil || userChat.DeletedAt.Valid {
			return echo.ErrForbidden
		}

		chat, ok := <-chatCh
		if !ok {
			return <-errCh
		}
		if chat.DeletedAt.Valid {
			return echo.ErrForbidden
		}
		return next(c)
	}
}

func GetCid(c echo.Context) int {
	return c.Get(keyCid).(int)
}

func GetChat(c echo.Context) *models.Chat {
	return c.Get(keyChat).(*models.Chat)
}

func GetUserChat(c echo.Context) *models.UserChat {
	return c.Get(keyUserChat).(*models.UserChat)
}
