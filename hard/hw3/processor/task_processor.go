package main

import (
	"bytes"
	"context"
	"fmt"
	"hard-hw1/models"
	"log"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func getCommand(task models.CodeTaskMessage) ([]string, error) {
	switch task.Translator {
	case "python3":
		return []string{
			"sh",
			"-c",
			"printf '%s' \"$CODE\" > /app/main.py && python3 /app/main.py",
		}, nil

	case "gcc":
		return []string{
			"sh",
			"-c",
			"printf '%s' \"$CODE\" > /app/main.c && gcc /app/main.c -o /app/program && /app/program",
		}, nil

	case "clang":
		return []string{
			"sh",
			"-c",
			"printf '%s' \"$CODE\" > /app/main.c && clang /app/main.c -o /app/program && /app/program",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported translator: %s", task.Translator)
	}
}

func processTask(task models.CodeTaskMessage) {
	cmd, err := getCommand(task)
	if err != nil {
		log.Printf("Failed to get command: %v", err)
		return
	}

	cli, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
		return
	}
	defer cli.Close()

	ctx := context.Background()

	response, err := cli.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				Image: "code-runner",
				Cmd:   cmd,
				Env: []string{
					"CODE=" + task.Code,
				},
			},
		},
	)
	if err != nil {
		log.Printf("Failed to create container: %v", err)
		return
	}

	log.Printf("Created container: %s", response.ID)
	defer func() {
		_, err := cli.ContainerRemove(
			ctx,
			response.ID,
			client.ContainerRemoveOptions{
				Force: true,
			},
		)
		if err != nil {
			log.Printf("Failed to remove container: %v", err)
		}
	}()

	_, err = cli.ContainerStart(
		ctx,
		response.ID,
		client.ContainerStartOptions{},
	)
	if err != nil {
		log.Printf("Failed to start container: %v", err)
		return
	}

	log.Printf("Started container: %s", response.ID)

	waitResult := cli.ContainerWait(
		ctx,
		response.ID,
		client.ContainerWaitOptions{},
	)

	result := <-waitResult.Result
	if result.Error != nil {
		log.Printf("Container wait error: %v", result.Error)
		return
	}

	log.Printf("Container exited with code: %d", result.StatusCode)

	logs, err := cli.ContainerLogs(
		ctx,
		response.ID,
		client.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		},
	)
	if err != nil {
		log.Printf("Failed to get container logs: %v", err)
		return
	}
	defer logs.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	_, err = stdcopy.StdCopy(&stdout, &stderr, logs)
	if err != nil {
		log.Printf("Failed to read container logs: %v", err)
		return
	}

	output := stdout.String() + stderr.String()

	log.Printf("Container output: %s", output)
	err = sendResult(task.TaskID, string(output))
	if err != nil {
		log.Printf("Failed to send result: %v", err)
		return
	}
}
