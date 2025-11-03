// ///////////////////////////////////////
// edit.go - Edit a command line in vi
// Mike Schilli, 2020 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"github.com/kballard/go-shellquote"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func editCmd(args *[]string) {
	tmp, err := ioutil.TempFile("/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	qargs := []string{}
	for _, arg := range *args {
		qargs = append(qargs, shellquote.Join(arg))
	}

	b := []byte(strings.Join(qargs, " "))
	err = ioutil.WriteFile(
		tmp.Name(), b, 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("vi", tmp.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	str, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		panic(err)
	}
	line := strings.TrimSuffix(string(str), "\n")
	*args, err = shellquote.Split(line)
	if err != nil {
		panic(err)
	}
}
