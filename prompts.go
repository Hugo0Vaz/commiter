package main

import "fmt"

func englishPrompt(diff string) string {
	return fmt.Sprintf(`Using the commit message format delimited in triple quotes and the code diff provided in the end of the prompt analyze and generate a commit massage in the same format.

"""
(<type>)[optional scope]: <description>

[optional body (Longer description) ]
"""

the commit <types> available are as follows in the table:

|Commit Type|Title|Description|Emoji|
|---|---|---|:-:|
|'feat'|Features|A new feature|✨|
|'fix'|Bug Fixes|A bug Fix|🐛|
|'docs'|Documentation|Documentation only changes|📚|
|'style'|Styles|Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)|💎|
|'refactor'|Code Refactoring|A code change that neither fixes a bug nor adds a feature|📦|
|'perf'|Performance Improvements|A code change that improves performance|🚀|
|'test'|Tests|Adding missing tests or correcting existing tests|🚨|
|'build'|Builds|Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)|🛠|
|'ci'|Continuous Integrations|Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)|⚙️|
|'chore'|Chores|Other changes that don't modify src or test files|♻️|
|'revert'|Reverts|Reverts a previous commit|🗑|

%s`, diff)
}

func portuguesePrompt(diff string) string {
	return fmt.Sprintf(`Analise esse diff do git e forneça uma mensagem de confirmação curta e descritiva e outra mais longa e detalhada.

A saída deve seguir este modelo: 

(ação)[arquivo ou parte do sistema]: descrição curta.

	Descrição longa e mais detalhada aqui 

Diff:
%s`, diff)
}
