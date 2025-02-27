package git_commands

import (
	"io/ioutil"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type FileCommands struct {
	*common.Common

	cmd    oscommands.ICmdObjBuilder
	config *ConfigCommands
	os     FileOSCommand
}

type FileOSCommand interface {
	Getenv(string) string
}

func NewFileCommands(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
	config *ConfigCommands,
	osCommand FileOSCommand,
) *FileCommands {
	return &FileCommands{
		Common: common,
		cmd:    cmd,
		config: config,
		os:     osCommand,
	}
}

// Cat obtains the content of a file
func (self *FileCommands) Cat(fileName string) (string, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", nil
	}
	return string(buf), nil
}

func (self *FileCommands) GetEditCmdStr(filename string, lineNumber int) (string, error) {
	editor := self.UserConfig.OS.EditCommand

	if editor == "" {
		editor = self.config.GetCoreEditor()
	}
	if editor == "" {
		editor = self.os.Getenv("GIT_EDITOR")
	}
	if editor == "" {
		editor = self.os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = self.os.Getenv("EDITOR")
	}
	if editor == "" {
		if err := self.cmd.New("which vi").DontLog().Run(); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return "", errors.New("No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
	}

	templateValues := map[string]string{
		"editor":   editor,
		"filename": self.cmd.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	editCmdTemplate := self.UserConfig.OS.EditCommandTemplate
	return utils.ResolvePlaceholderString(editCmdTemplate, templateValues), nil
}
