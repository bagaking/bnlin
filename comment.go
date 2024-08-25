package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"
	"github.com/sirupsen/logrus"

	"github.com/bagaking/botheater/bot"
	"github.com/bagaking/botheater/history"
)

//go:embed role_prompt.md
var RolePrompt string

//go:embed role_few_shot_example.txt
var FowShotExample string

// autoComment generates a commit comment based on the provided diff information.
func autoComment(ctx context.Context, task string, exe ExecutionGroup) error {
	// disable logrus to hide bot debug
	logrus.SetOutput(io.Discard)

	// Check if the -diff flag is provided
	if task == "" {
		return irr.Error("Please provide the task information")
	}

	// Use command-line flags for access key and secret key if provided
	exe = exe.Use(strings.Replace(RolePrompt, "{{role_few_shot_example.txt}}", FowShotExample, -1))
	if err := exe.Assert(); err != nil {
		return irr.Wrap(err, "assertion failed")
	}

	answer, err := SimpleQuestion(ctx, exe.Bot(ctx), buildTask(task))
	if err != nil {
		return irr.Wrap(err, "failed to generate comment")
	}

	// Print the generated comment
	// fmt.Println(comment)

	fmt.Println("=== Command Start ===")
	if err = execute(answer); err != nil {
		return irr.Wrap(err, "failed to execute command")
	}
	fmt.Println("=== Command Finish ===")
	return nil
}

func buildTask(task string) string {
	os, ver, lang := getOSInfo()
	return fmt.Sprintf(`我的操作系统是: %s, 版本是: %s, 控制台的语言是: %s
请为我生成能直接运行，并谨慎仔细的命令行脚本, 完成以下任务:
%s`,
		os, ver, lang,
		task)
}

// SimpleQuestion sends a question to the bot and returns the generated answer.
//
// ctx: The context for the request.
// endpoint: The endpoint of the bot service.
// prompt: The prompt for generating the answer.
// question: The question to be answered.
//
// Returns the generated answer and any error encountered.
func SimpleQuestion(ctx context.Context, bot *bot.Bot, question string) (string, error) {
	wlog.ByCtx(ctx, "SimpleQuestion").Tracef("prompt===\n%s\nquestion===\n%s\n",
		bot.Prompt,
		question,
	)

	// Create a channel to signal the completion of the question
	done := make(chan struct{})

	// Start a goroutine to display the loading animation
	go func() {
		frames := []string{"-", "\\", "|", "/"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\rAI Processing %s", frames[i%len(frames)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Perform the question operation
	answer, err := bot.Question(ctx, history.NewHistory(), question)

	// Clear the loading animation
	fmt.Print("\r")
	fmt.Print("\n")

	// Signal the completion of the question
	close(done)

	return answer, err
}
