# Golang Binary Size Analyzer

![Result example](/docs/graph.png?raw=true)

It's quite slow for now. I know.

Examle usage:

```
gba --tree --level 2 github.com/cli/cli/cmd/gh

github.com/cli/cli/cmd/gh up to 59.1 MiB
├── github.com/cli/cli/command up to 58.9 MiB
├────── github.com/AlecAivazis/survey/v2 up to 14.9 MiB
├────── github.com/AlecAivazis/survey/v2/core up to 12.7 MiB
├────── github.com/cli/cli/api up to 53.2 MiB
├────── github.com/cli/cli/context up to 54.8 MiB
├────── github.com/cli/cli/git up to 11.1 MiB
├────── github.com/cli/cli/internal/cobrafish up to 17.0 MiB
├────── github.com/cli/cli/internal/config up to 54.5 MiB
├────── github.com/cli/cli/internal/ghrepo up to 9.2 MiB
├────── github.com/cli/cli/internal/run up to 9.8 MiB
├────── github.com/cli/cli/pkg/githubtemplate up to 11.7 MiB
├────── github.com/cli/cli/pkg/httpmock up to 26.6 MiB
├────── github.com/cli/cli/pkg/surveyext up to 15.0 MiB
├────── github.com/cli/cli/pkg/text up to 5.4 MiB
├────── github.com/cli/cli/utils up to 36.2 MiB
├────── github.com/google/shlex up to 9.5 MiB
├────── github.com/spf13/cobra up to 16.9 MiB
├────── github.com/spf13/pflag up to 14.1 MiB
├────── golang.org/x/crypto/ssh/terminal up to 9.9 MiB
├── github.com/cli/cli/internal/config up to 54.5 MiB
├────── github.com/cli/cli/api up to 53.2 MiB
├────── github.com/cli/cli/auth up to 25.5 MiB
├────── github.com/mitchellh/go-homedir up to 7.9 MiB
├────── gopkg.in/yaml.v3 up to 11.4 MiB
├── github.com/cli/cli/update up to 54.4 MiB
├────── github.com/cli/cli/api up to 53.2 MiB
├────── github.com/hashicorp/go-version up to 10.2 MiB
├────── gopkg.in/yaml.v3 up to 11.4 MiB
├── github.com/cli/cli/utils up to 36.2 MiB
├────── github.com/briandowns/spinner up to 11.3 MiB
├────── github.com/charmbracelet/glamour up to 34.9 MiB
├────── github.com/cli/cli/internal/run up to 9.8 MiB
├────── github.com/cli/cli/pkg/browser up to 9.9 MiB
├────── github.com/cli/cli/pkg/text up to 5.4 MiB
├────── github.com/mattn/go-colorable up to 10.5 MiB
├────── github.com/mattn/go-isatty up to 9.8 MiB
├────── github.com/mgutz/ansi up to 11.0 MiB
├────── golang.org/x/crypto/ssh/terminal up to 9.9 MiB
├── github.com/mgutz/ansi up to 11.0 MiB
├── github.com/spf13/cobra up to 16.9 MiB

```

```
gba --dot github.com/cli/cli/cmd/gh

dot -Tsvg graph.dot > graph.svg
```

Next steps:

1. Sorted list output
2. Graph output to .svg
3. Graph output to .png
4. Top N most heavy dependencies
5. Omit dependencies less that N bytes

