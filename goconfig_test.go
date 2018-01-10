package goconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, 1, int(lineno))
	assert.Equal(t, map[string]string{}, config)
}

func TestNoNewLine(t *testing.T) {
	validConfig := "[user] name = Danyel"
	config, lineno, err := Parse([]byte(validConfig))
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, int(lineno))
	assert.Equal(t, map[string]string{"user.name": "Danyel"}, config)
}

func TestExtended(t *testing.T) {
	validConfig := `[http "https://my-website.com"] sslVerify = false`
	config, lineno, err := Parse([]byte(validConfig))
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, int(lineno))
	assert.Equal(t, map[string]string{`http.https://my-website.com.sslverify`: "false"}, config)
}

func ExampleParse() {
	gitconfig := "configs/danyel.gitconfig"
	bytes, err := ioutil.ReadFile(gitconfig)
	if err != nil {
		log.Fatalf("Couldn't read file %v\n", gitconfig)
	}

	config, lineno, err := Parse(bytes)
	if err != nil {
		log.Fatalf("Error on line %d: %v\n", lineno, err)
	}
	fmt.Println()
	fmt.Println(lineno)
	fmt.Println(config["user.name"])
	fmt.Println(config["user.email"])
	// Output:
	// 10
	// Danyel Bayraktar
	// cydrop@gmail.com
}

func BenchmarkParse(b *testing.B) {
	gitconfig := "configs/danyel.gitconfig"
	bytes, err := ioutil.ReadFile(gitconfig)
	if err != nil {
		b.Fatalf("Couldn't read file %v: %s\n", gitconfig, err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(bytes)
	}
}
