package main

import "fmt"

func englishPrompt(diff string) string {
	return fmt.Sprintf(`Using the following commit message format and the code diff provided in the end of the prompt analyze and generate a commit massage in the same format.

(<type>)[scope]: <description>\n<detailed description>

The <description> should be a short and informative description of the changes.

The <detailed description> should be longer and more descriptive of the changes.

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

Diff:
%s`, diff)
}

func portuguesePrompt(diff string) string {
	return fmt.Sprintf(`Usando o formato de mensagem a seguir e o diff de código fornecida no final do prompt, analise e gere uma mensagem de confirmação no mesmo formato.

(<type>)[escopo]: <description>\n<longer description>

A <description> deve ser curta e informativa.

A <longer description> deve ser mais longa e mais descritiva das mudanças

Os <types> de commit disponíveis são os seguintes na tabela:

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

Diff:
%s`, diff)
}
