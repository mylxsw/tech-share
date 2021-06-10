package bcrypt_test

import (
	"testing"

	"github.com/mylxsw/tech-share/pkg/bcrypt"
	"github.com/stretchr/testify/assert"
)

func TestHashMatch(t *testing.T) {
	password := "abcdefg"

	pass1, _ := bcrypt.Hash(password)
	pass2, _ := bcrypt.Hash(password)
	assert.NotEqual(t, pass1, pass2)

	assert.True(t, bcrypt.Match(password, pass1))
	assert.True(t, bcrypt.Match(password, pass2))
}
