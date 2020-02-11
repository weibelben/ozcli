package ozcli

import (
	"bytes"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	completionLong = `
		Output shell completion code for the specified shell (bash or zsh).
		The shell code must be evaluated to provide interactive
		completion of ozcli commands. This can be done by sourcing it from
		the .bash_profile.

		Note: this requires the bash-completion framework.

		If bash-completion is not installed on Linux, please install the 'bash-completion' package
		via your distribution's package manager.

		Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2`

	completionExample = `

		# Load the ozcli completion code for bash into the current shell
		source <(ozcli completion bash)

		# Load the ozcli completion code for zsh[1] into the current shell
		source <(ozcli completion zsh)`
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:     "completion SHELL",
	Short:   "Output shell completion code for the specified shell (bash or zsh)",
	Long:    completionLong,
	Example: completionExample,
	Run: func(cmd *cobra.Command, args []string) {
		RunCompletion(os.Stdout, cmd, args)
	},
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"bash", "zsh"},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

// RunCompletion checks given arguments and executes command
func RunCompletion(out io.Writer, cmd *cobra.Command, args []string) {
	shell := args[0]
	switch shell {
	case "bash":
		_ = runCompletionBash(out, cmd.Parent())
	case "zsh":
		_ = runCompletionZsh(out, cmd.Parent())
	default:
		log.Fatal("UnExpected shell type.")
	}
}

func runCompletionBash(out io.Writer, ozcli *cobra.Command) error {

	return ozcli.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, ozcli *cobra.Command) error {
	zshHead := "#compdef ozcli\n"

	_, _ = out.Write([]byte(zshHead))

	zshInitialization := `
__ozcli_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__ozcli_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__ozcli_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__ozcli_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__ozcli_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__ozcli_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__ozcli_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__ozcli_filedir() {
	local RET OLD_IFS w qw

	__ozcli_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi

	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"

	IFS="," __ozcli_debug "RET=${RET[@]} len=${#RET[@]}"

	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__ozcli_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__ozcli_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
	printf %q "$1"
    fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__ozcli_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__ozcli_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__ozcli_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__ozcli_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__ozcli_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__ozcli_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__ozcli_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	_, _ = out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	_ = ozcli.GenBashCompletion(buf)
	_, _ = out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}

__ozcli_bash_source <(__ozcli_convert_bash_to_zsh)
_complete ozcli 2>/dev/null
`
	_, _ = out.Write([]byte(zshTail))
	return nil
}
