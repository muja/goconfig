package goconfig

import "testing"
import "io/ioutil"
import "github.com/stretchr/testify/assert"

func TestDanyel(t *testing.T) {
	filename := "configs/danyel.gitconfig"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Reading file %v failed", filename)
	}
	config, lineno, err := Parse(bytes)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10, int(lineno))
	assert.Equal(t, "Danyel Bayraktar", config["user.name"])
	assert.Equal(t, "cydrop@gmail.com", config["user.email"])
	assert.Equal(t, "subl -w", config["core.editor"])
	assert.Equal(t, `!git config --get-regexp 'alias.*' | colrm 1 6 | sed 's/[ ]/ = /' | sort`, config["alias.aliases"])
}

func TestInvalidKey(t *testing.T) {
	invalidConfig := ".name = Danyel"
	config, lineno, err := Parse([]byte(invalidConfig))
	assert.Equal(t, ErrInvalidKeyChar, err)
	assert.Equal(t, 2, int(lineno))
	assert.Equal(t, map[string]string{}, config)
}
