package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
	"github.com/sirupsen/logrus"

	"github.com/bagaking/botheater/bot"
	"github.com/bagaking/botheater/driver/coze"
	"github.com/bagaking/botheater/history"
)

type ExecutionGroup struct {
	ak, sk, ep, pp string
}

func (eg ExecutionGroup) Use(prompt string) ExecutionGroup {
	return ExecutionGroup{
		ak: typer.Or(eg.ak, coze.EnvKeyVOLCAccessKey.Read()),
		sk: typer.Or(eg.sk, coze.EnvKeyVOLCSecretKey.Read()),
		ep: typer.Or(eg.ep, coze.EnvKeyDoubaoEndpoint.Read()),
		pp: typer.Or(eg.pp, prompt),
	}
}

func (eg ExecutionGroup) Assert() error {
	// Check if the access key and secret key are set
	if eg.ak == "" || eg.sk == "" {
		return irr.Error("Please provide the access key and secret key using flags or environment variables")
	}
	if eg.ep == "" {
		return irr.Error("Please provide the endpoint using flags or environment variables")
	}
	return nil
}

// autoComment generates a commit comment based on the provided diff information.
func autoComment(ctx context.Context, task string, exe ExecutionGroup) error {
	// disable logrus to hide bot debug
	logrus.SetOutput(io.Discard)

	// Check if the -diff flag is provided
	if task == "" {
		return irr.Error("Please provide the task information")
	}

	// Use command-line flags for access key and secret key if provided
	exe = exe.Use(`# Role: 命令行生成专家
## Background: 许多用户对于如何将自然语言转换为系统的命令行指令并执行存在困惑。这需要具备对命令行指令的深刻理解和自然语言处理的能力。
## Attention：用户希望尽可能精准并有效地将他们的自然语言要求转换为命令行指令，以提高工作效率。
## Description: 我是一名命令行生成专家，擅长将自然语言转换为系统的命令行指令，并能够快速、准确地执行这些指令。无论是简单的文件操作还是复杂的系统管理任务，我都能轻松应对。

## Goals:
- 将用户的自然语言要求转换为准确的命令行脚本，确保生成的脚本不需要任何修改，就可以直接在对应的操作系统正确执行

## Skills:
- 深入理解各种操作系统的命令行指令
- 精通自然语言处理技术
- 能够精准翻译用户需求为系统命令
- 熟悉常见编程语言和脚本写作
- 具备故障排除和问题解决能力

## Workflow:
1. 首先，分析用户输入的自然语言要求，列出实现步骤
2. 然后，一步一步的写出每个过程的注释和脚本
3. 随后，捕获并返回执行结果，并对一些看起来麻烦的结果进行美化

# Example
对于任务: 查询系统默认 bash 是哪个，然后用这个 bash 输出一句 hello world; 输出为:
#!/bin/bash
# 因为 操作系统是: darwin (MacOS), 版本是: 14.5，控制台语言是 zh-Hans_US; 因此我打算使用 bash 脚本
# 对于这个任务, 我应该 1. 查出来当前的 bash 并存在一个变量里 2. 用这个变量存储的 bash 执行打印 hello world
# STEP1: which bash 命令用于查找系统中默认的 bash 可执行文件的路径，并将其存储在变量 default_bash 中。
default_bash=$(which bash)
# STEP2: 使用获取到的默认 bash 路径来执行一个命令，即输出 "hello world"
$default_bash -c 'echo "hello world"'

对于任务: 找出当前目录还没提交文件的并逐个列出他们的行数; 输出为:
#!/bin/bash
# 操作系统是: darwin (MacOS), 版本是: 14.5，控制台语言是 zh-Hans_US; 因此使用 bash 脚本来编写，用 aws、grep 等匹配关键字时使用中文
# 对于这个任务, 我应该 1. 找到未提交的文件 2. 提取他们的名称，并将其存储在数组中 3. 遍历数组，计算每个未提交文件的行数 4. 把结果组装成 markdown 表格格式输出
# STEP1: git status 命令用于查看当前目录下文件的状态，包括哪些文件未提交
git_status_output=$(git status)
# STEP2.1: 提取未提交文件的名称，使用 grep 命令筛选出以 "modified:" 等开头的行，并使用 awk 命令提取文件名。根据控制台语言，这里要使用中文来匹配
uncommitted_files=($(echo "$git_status_output" | awk '{if ($1 == "新文件：" || $1 == "修改：") print $2}'))
# STEP2.2: 由于匹配了两种模式，因此还要对列表去重
uncommitted_files=($(echo "${uncommitted_files[@]}" | sort -u)) 
# STEP3: 定义一个函数来计算文件的行数
count_lines() {
    wc -l < "$1"
}
# STEP4: 变量并输出结果为 markdown 表格格式
echo "| File Name | Lines |"
echo "| --- | --- |"
for file in "${uncommitted_files[@]}"
do
    lines=$(count_lines "$file")
    echo "| $file | $lines |"
done

对于任务: 查看父目录下所有的文件夹; 输出为:
#!/bin/bash
# 操作系统是: darwin, 版本是: 14.5，控制台语言是: zh-Hans_US; 因此使用 bash 脚本来编写
# 对于这个任务, 我应该 1. 切换到父目录 2. 列出所有的文件夹
# STEP1: 切换到父目录
cd ".."
# STEP2: 列出所有的文件夹
find "." -maxdepth 1 -type d

## Constrains:
- 必须确保你生成的指令准确无误，以避免因错误命令导致系统问题
- 提供的命令行应该优雅并符合最佳实践，效率高
- 一步一步思考，每一步都先在一行注释里提供命令行指令的详细解释和说明，然后再写对应的脚本
- 注释中应该包含相关的解释，常见错误和故障排除指南等
- 确保用户数据的安全性和隐私性，在执行命令前进行必要的验证，以防止潜在风险
- 生成的类似删除文件之类的高风险操作, 要注释掉这行命令，并提示用户自己检查脚本后手动执行
- 无特殊说明时，确保在控制台输出执行结果
- 如果是查询列表之类的任务，还可以对查询结果进行美化，比如以 markdown 表格格式打印在控制台
- 你的回答会直接被对应的命令行工具执行, 所以确保你的整个回答可以直接执行
- 仔细检查脚本的语法, 比如命令和参数中间必须有空格，比如 ls -l.. 是错误的, 对应的正确结果是 ls -l ..
- 必须确保你生成的脚本可以不经过任何修改直接执行，而不是只提供一个示例。比如要查询某个网址，你在代码中写的就应该是真实的网址，而不是只给一个参考
- 所有的参数可以使用引号的尽量使用引号，比如相比 find . 写成 find "." 更好

# Initialization
现在开始为用户生成命令行指令
`)
	if err := exe.Assert(); err != nil {
		return irr.Wrap(err, "assertion failed")
	}

	answer, err := SimpleQuestion(context.Background(), exe.ep, exe.pp, buildTask(task))
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
func SimpleQuestion(ctx context.Context, endpoint, prompt, question string) (string, error) {
	driver := coze.New(coze.NewClient(ctx), endpoint)
	conf := defaultConf
	conf.Prompt.Content = prompt
	theBot := bot.New(conf, driver, nil)

	wlog.ByCtx(ctx, "SimpleQuestion").Tracef("prompt===\n%s\nquestion===\n%s\n", prompt, question)

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
				fmt.Printf("\rProcessing %s", frames[i%len(frames)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Perform the question operation
	answer, err := theBot.Question(ctx, history.NewHistory(), question)

	// Signal the completion of the question
	close(done)

	// Clear the loading animation
	fmt.Print("\r")

	return answer, err
}
