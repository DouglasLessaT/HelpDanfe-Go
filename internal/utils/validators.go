package utils

import (
	"regexp"
	"strconv"
)

// ValidarChaveAcesso valida se a chave de acesso tem 44 dígitos
func ValidarChaveAcesso(chave string) bool {
	if len(chave) != 44 {
		return false
	}

	// Verifica se contém apenas números
	matched, _ := regexp.MatchString(`^\d{44}$`, chave)
	return matched
}

// ValidarCNPJ valida um CNPJ
func ValidarCNPJ(cnpj string) bool {
	// Remove caracteres não numéricos
	cnpj = regexp.MustCompile(`[^\d]`).ReplaceAllString(cnpj, "")

	if len(cnpj) != 14 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	if regexp.MustCompile(`^(\d)\1{13}$`).MatchString(cnpj) {
		return false
	}

	// Calcula primeiro dígito verificador
	sum := 0
	weight := 5
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weight
		weight--
		if weight < 2 {
			weight = 9
		}
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}

	// Calcula segundo dígito verificador
	sum = 0
	weight = 6
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weight
		weight--
		if weight < 2 {
			weight = 9
		}
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}

	// Verifica se os dígitos verificadores estão corretos
	expectedDigit1, _ := strconv.Atoi(string(cnpj[12]))
	expectedDigit2, _ := strconv.Atoi(string(cnpj[13]))

	return digit1 == expectedDigit1 && digit2 == expectedDigit2
}

// ValidarCPF valida um CPF
func ValidarCPF(cpf string) bool {
	// Remove caracteres não numéricos
	cpf = regexp.MustCompile(`[^\d]`).ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	if regexp.MustCompile(`^(\d)\1{10}$`).MatchString(cpf) {
		return false
	}

	// Calcula primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}

	// Calcula segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}

	// Verifica se os dígitos verificadores estão corretos
	expectedDigit1, _ := strconv.Atoi(string(cpf[9]))
	expectedDigit2, _ := strconv.Atoi(string(cpf[10]))

	return digit1 == expectedDigit1 && digit2 == expectedDigit2
}

// ValidarCodigoBarras valida um código de barras
func ValidarCodigoBarras(codigo string) bool {
	// Remove espaços
	codigo = regexp.MustCompile(`\s`).ReplaceAllString(codigo, "")

	// Verifica se tem 44 ou 47 dígitos (padrão FEBRABAN)
	if len(codigo) != 44 && len(codigo) != 47 {
		return false
	}

	// Verifica se contém apenas números
	matched, _ := regexp.MatchString(`^\d+$`, codigo)
	return matched
}

// ValidarLinhaDigitavel valida uma linha digitável
func ValidarLinhaDigitavel(linha string) bool {
	// Remove espaços e pontos
	linha = regexp.MustCompile(`[\s\.]`).ReplaceAllString(linha, "")

	// Verifica se tem 47 dígitos
	if len(linha) != 47 {
		return false
	}

	// Verifica se contém apenas números
	matched, _ := regexp.MatchString(`^\d{47}$`, linha)
	return matched
}

// ValidarEmail valida um email
func ValidarEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidarTelefone valida um telefone
func ValidarTelefone(telefone string) bool {
	// Remove caracteres não numéricos
	telefone = regexp.MustCompile(`[^\d]`).ReplaceAllString(telefone, "")

	// Verifica se tem 10 ou 11 dígitos
	if len(telefone) != 10 && len(telefone) != 11 {
		return false
	}

	return true
}
