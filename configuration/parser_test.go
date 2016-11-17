package configuration_test

import (
	. "github.com/homedepot/github-webhook/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	"strings"
)

func TestParserMinimalYaml(t *testing.T) {

	data := []byte(`
triggers:
  - name: master
    event: push
    run: {path: a/b }
`)
	config, err := LoadBytes(data)

	assert.Nil(t, err)

	assert.Equal(t,len(config.Triggers), 1 )
	assert.Equal(t,config.Triggers[0].Name, "master" )
	assert.Equal(t,config.Triggers[0].Event, "push" )
	assert.Equal(t,config.Triggers[0].Run.Path, "a/b" )

}

func TestParserBadEvent(t *testing.T) {

	data := []byte(`
triggers:
  - name: master
    event: something
    run: {path: a/b }
`)
	_, err := LoadBytes(data)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err),"unsupported github event type" ))

}

func TestParserMissingEvent(t *testing.T) {

	data := []byte(`
triggers:
  - name: master
    run: {path: a/b }
`)
	_, err := LoadBytes(data)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err),"missing github event type" ))

}

func TestParserMissingName(t *testing.T) {

	data := []byte(`
triggers:
  - event: push
    run: {path: a/b }
`)
	_, err := LoadBytes(data)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err),"missing trigger name" ))

}

func TestParserMissingRun(t *testing.T) {

	data := []byte(`
triggers:
  - name: trigger
    event: push
`)
	_, err := LoadBytes(data)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err),"missing trigger run" ))

}

func TestParserRunArgs(t *testing.T) {

	data := []byte(`
triggers:
  - name: master
    event: push
    run:
      path: a/b
      args:
        - -exc
        - |
          a
          b
`)
	config, err := LoadBytes(data)

	assert.Nil(t, err)

	assert.Equal(t,len(config.Triggers), 1 )
	assert.Equal(t,config.Triggers[0].Run.Args, []string{"-exc","a\nb\n"} )

}

func TestParserLoadFile(t *testing.T) {

	config, err := LoadFile("../testdata/triggers.yaml")

	assert.Nil(t, err)

	assert.Equal(t,len(config.Triggers), 1 )
	assert.Equal(t,config.Triggers[0].Name, "master" )
	assert.Equal(t,config.Triggers[0].Event, "push" )
	assert.Equal(t,config.Triggers[0].Run.Path, "a/b" )
	assert.Equal(t,config.Triggers[0].Run.Args, []string{"-exc","a\nb"} )
	assert.Equal(t,config.Triggers[0].Rules["ref"], "refs/heads/master" )
	assert.Equal(t,config.Triggers[0].Rules["head_commit.author.username"], "username" )
	assert.Equal(t,config.Triggers[0].Rules["pusher.name"], "username" )
	assert.Equal(t,config.Triggers[0].Rules["commits.modified"], "README.md" )

}