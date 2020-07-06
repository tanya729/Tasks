package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalln("Too few args")
	}
	env, err := ReadDir(args[0])
	if err != nil {
		log.Fatalf("Error : %s", err)
	}
	code := RunCmd(args[1:], env)
	os.Exit(code)
}

// RunCmd will run []string cmd with map[string]string env
func RunCmd(cmd []string, env map[string]string) int {
	// fmt.Printf("%v", cmd)
	command := exec.Command(cmd[0])
	if len(cmd) > 1 {
		command.Args = cmd[1:]
	}
	result := make([]string, 0, len(env))
	for key, value := range env {
		result = append(result, key+"="+value)
	}
	command.Env = append(os.Environ(), result...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}
	return 0
}

//ReadDir parse folder for enviroment values in files
func ReadDir(dir string) (map[string]string, error) {
	env := make(map[string]string)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return env, err
	}
	for _, file := range files {
		stat, err := os.Stat(dir + "/" + file.Name())
		if err != nil {
			return env, err
		}
		if stat.IsDir() == true {
			continue
		}
		content, err := ioutil.ReadFile(dir + "/" + file.Name())
		if err != nil {
			return env, err
		}
		if string(content) == "" {
			continue
		}
		env[file.Name()] = string(content)
	}
	return env, err
}
