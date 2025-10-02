package test

import (
	"testing"

	"server/utils"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndParseToken(t *testing.T) {
	// 先创建再解析
	token, err := utils.CreateToken(1, "李华", "lihua@zju.edu.cn", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}
