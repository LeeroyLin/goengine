package cli

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"strings"
)

type Cmd struct {
	root       cobra.Command
	welcomeStr string

	children map[string]*Cmd

	strMap    map[*string]string
	boolMap   map[*bool]bool
	intMap    map[*int]int
	int64Map  map[*int64]int64
	uint32Map map[*uint32]uint32
}

func NewRootCmd(welcomeStr string) *Cmd {
	c := &Cmd{
		root: cobra.Command{
			Use:   "root-cli",
			Short: "root-cli",
			Long:  "root-cli",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Run root cmd.")
			},
		},
		welcomeStr: welcomeStr,
		children:   make(map[string]*Cmd),
	}

	return c
}

func (c *Cmd) Run() {
	// 初始化 readline（优化输入体验，支持历史记录、光标移动）
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",         // 命令行提示符
		HistoryFile:     ".cli_history", // 历史记录保存路径（当前目录下）
		AutoComplete:    nil,            // 可扩展自动补全
		InterruptPrompt: "^C",           // Ctrl+C 中断提示
		EOFPrompt:       "exit",         // Ctrl+D 等价于 exit 命令
	})
	if err != nil {
		elog.Error("[Cmd] init readline failed. err", err)
		return
	}
	defer rl.Close() // 程序退出时关闭 readline

	// 打印欢迎信息
	elog.Info("[Cmd]========")
	elog.Info(c.welcomeStr)
	elog.Info("[Cmd]========")

	// 监听
	c.listen(rl)
}

func (c *Cmd) NewSubCmd(name, short, long string, runHandler func(cmd *cobra.Command)) *Cmd {
	sub := &Cmd{
		root: cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
		},
		children: make(map[string]*Cmd),
	}

	sub.root.Run = func(cmd *cobra.Command, args []string) {
		runHandler(cmd)

		sub.resetAllFlags()
	}

	c.children[sub.root.Use] = sub
	c.root.AddCommand(&sub.root)

	return sub
}

func (c *Cmd) AddSubCmd(sub *Cmd) {
	c.root.AddCommand(&sub.root)
}

func (c *Cmd) listen(rl *readline.Instance) {
	for {
		// 读取用户输入（自动处理换行，支持 Ctrl+C 中断）
		input, err := rl.Readline()
		if err != nil {
			// 若输入 EOF（Ctrl+D）或其他错误，终止循环
			if errors.Is(err, readline.ErrInterrupt) {
				elog.Info("[Cmd] exit")
				break
			}

			elog.Info("[Cmd] wrong input, err:", err)
			continue
		}

		// 处理空输入（直接跳过）
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 4. 将用户输入拆分为「命令参数列表」（模拟 os.Args）
		cmdArgs := append([]string{c.root.Use}, strings.Fields(input)...)

		// 5. 调用 Cobra 解析命令并执行
		// 重置根命令的参数，避免多次执行残留状态
		c.root.SetArgs(cmdArgs[1:])
		if err := c.root.Execute(); err != nil {
			elog.Error("[Cmd] command execute failed, err:", err)
		}
	}
}
