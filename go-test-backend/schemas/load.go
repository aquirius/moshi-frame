package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var dbName string

func init() {
	flag.StringVar(&dbName, "db", "sprout", "moshi")
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
	out := []string{}

	greenhouse, err := os.ReadFile("./schemas/greenhouses.sql")
	if err != nil {
		panic(err.Error())
	}
	userGreenhouse, err := os.ReadFile("./schemas/users-greenhouses.sql")
	if err != nil {
		panic(err.Error())
	}
	nutrients, err := os.ReadFile("./schemas/nutrients.sql")
	if err != nil {
		panic(err.Error())
	}
	plans, err := os.ReadFile("./schemas/plans.sql")
	if err != nil {
		panic(err.Error())
	}
	plants, err := os.ReadFile("./schemas/plants.sql")
	if err != nil {
		panic(err.Error())
	}
	pots, err := os.ReadFile("./schemas/pots.sql")
	if err != nil {
		panic(err.Error())
	}
	stacks, err := os.ReadFile("./schemas/stacks.sql")
	if err != nil {
		panic(err.Error())
	}
	users, err := os.ReadFile("./schemas/users.sql")
	if err != nil {
		panic(err.Error())
	}
	sprouts, err := os.ReadFile("./schemas/sprouts.sql")
	if err != nil {
		panic(err.Error())
	}

	out = append(out,
		string(greenhouse),
		string(nutrients),
		string(plans),
		string(plants),
		string(pots),
		string(stacks),
		string(users),
		string(userGreenhouse),
		string(sprouts))
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
		"sprout-backend",
		"mysql",
		"-uroot",
		"-pmoshi",
		"sprout",
	}, []byte(s)); err != nil {
		fmt.Errorf("failed to load sprout sql schemas: %w", err)
	}

}
