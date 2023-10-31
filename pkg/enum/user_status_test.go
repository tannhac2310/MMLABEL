package enum

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStatusParse(t *testing.T) {
	data, err := json.Marshal(UserStatusActive)
	assert.Nil(t, err)
	assert.Equal(t, "\""+UserStatusName[UserStatusActive]+"\"", string(data))

	var a UserStatus
	err = json.Unmarshal(data, &a)
	assert.Nil(t, err)
	assert.Equal(t, UserStatusActive, a)
}
