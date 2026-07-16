package other

import (
	"server/model/appTypes"

	"github.com/gofrs/uuid"
)

// CurrentUser 负责将api层取得的UUID和RoleID传递给service层
type CurrentUser struct {
	UUID   uuid.UUID       `json:"uuid"`
	RoleID appTypes.RoleID `json:"role_id"`
}
