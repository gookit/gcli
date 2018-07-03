# auto-completion 脚本编写

## bash 环境

bash 环境下的命令自动补全脚本

- 示例文件： [auto-completion.bash](auto-completion.bash)

### 说明

`complete -F` 后面接一个函数，该函数将输入三个参数：

1. 要补全的命令名
2. 当前光标所在的词
3. 当前光标所在的词的前一个词

生成的补全结果需要存储到COMPREPLY变量中，以待bash获取。

`complete` 选项参数：

- `-F function` 指定补全函数名
- `-A file` 表示默认的动作是补全文件名，也即是如果bash找不到补全的内容，就会默认以文件名进行补全

参数接收：

```bash
local cur prev

// 方式 1

_get_comp_words_by_ref -n = cur prev

// 方式 2

pre="$3"
cur="$2"

// 方式 3

pre=${COMP_WORDS[COMP_CWORD-1]} # COMP_WORDS变量是一个数组，存储着当前输入所有的词
cur=${COMP_WORDS[COMP_CWORD]}
```

### 参考链接

- https://segmentfault.com/a/1190000002968878

## zsh 环境

zsh 环境下的命令自动补全脚本

提示：

- `echo $fpath` zsh在启动时会加载`$fpath`路径下的脚本文件,可以在这些文件夹下找文件参考

### 参考链接 

- https://segmentfault.com/a/1190000002994217
- https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org
