package main

import "fmt"

func englishPrompt(diff string) string {
	return fmt.Sprintf(`Analyze this git diff and provide a short and descriptive commit message and also a longer, detailed one.

The output should follow this template: 

(action)[file or system part]: short description.

	Longer and more detailed description here 

Diff:
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
