package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

var dbName string

func init() {
	flag.StringVar(&dbName, "db", "milon", "milon")
	flag.Parse()
}

func run(name string, cmds []string, stdin []byte) error {
	cmd := exec.Command(name, cmds...)

	// capture outputs
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	assembleErr := func(message string, e error) error {
		return fmt.Errorf("%s: %s\n%s\n%s", message, stdout.String(), stderr.String(), e)
	}

	// fwd data to execution?
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}

	err := cmd.Start()
	if err != nil {
		return assembleErr("Failed to start command", err)
	}

	if err := cmd.Wait(); err != nil {
		return assembleErr("Failed to wait for command", err)
	}
	return nil
}

func main() {
	users, err := filepath.Abs("./schemas/users.sql")
	if err != nil {
		panic(err.Error())
	}

	if len(users) <= 0 {
		fmt.Errorf("no users schema found")
	}
	out := []string{}

	c, err := ioutil.ReadFile(users)
	if err != nil {
		panic(err.Error())
	}

	out = append(out, string(c))
	s := ""
	for _, v := range out {
		s += v
		if !strings.HasSuffix(v, ";") {
			s += ";"
		}
		s += "\n\n"
	}
	s = `
	SET FOREIGN_KEY_CHECKS=0;
	` + s + `
	SET FOREIGN_KEY_CHECKS=1;
	`

	if err := run("docker", []string{
		"exec",
		"-i",
		"go-test-backend",
		"mysql",
		"-umilon",
		"-pmilon",
		"milon",
	}, []byte(s)); err != nil {
		fmt.Errorf("failed to load milon sql schemas: %w", err)
	}

}
