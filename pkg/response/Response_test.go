package response

import (
	"github.com/stretchr/testify/assert"
	"golek_notifications_service/pkg/models"
	status "golek_notifications_service/pkg/operation_status"
	"testing"
)

func TestResponse(t *testing.T) {
	t.Run("create user response", func(t *testing.T) {
		response := New(true, status.Failed, "ok", models.Message{})
		assert.True(t, response.GetStatus())
		assert.Equal(t, status.Failed, response.GetStatusCode())
		assert.Equal(t, "ok", response.GetMessage())
		response.SetMessage("ko")
		assert.Equal(t, "ko", response.GetMessage())
		assert.True(t, response.ErrorIs(status.Failed))
		assert.False(t, response.ErrorIs(status.Success))
		assert.False(t, response.IsFailed())
		assert.IsType(t, models.Message{}, response.GetData())
		assert.IsType(t, Response{}, response)
		assert.Equal(t, response.ToMapStringInterface()["data"], response.GetData())
	})

}
