package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/argon2"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go [encrypt|decrypt] <diretório>")
		return
	}

	fmt.Print("Digite a senha: ")
	var password string
	fmt.Scanln(&password)

	key := generateKey(password) // Gera a chave a partir da senha do usuário

	action := os.Args[1]
	directory := os.Args[2]

	if action == "encrypt" {
		err := processDirectory(directory, key, true)
		if err != nil {
			fmt.Println("Erro ao criptografar diretório:", err)
			return
		}
		err = lockDirectory(directory)
		if err != nil {
			fmt.Println("Erro ao trancar o diretório:", err)
			return
		}
		fmt.Println("Diretório criptografado e trancado com sucesso!")
	} else if action == "decrypt" {
		err := unlockDirectory(directory)
		if err != nil {
			fmt.Println("Erro ao destrancar o diretório:", err)
			return
		}
		err = processDirectory(directory, key, false)
		if err != nil {
			fmt.Println("Erro ao descriptografar diretório:", err)
			return
		}
		fmt.Println("Diretório descriptografado e destrancado com sucesso!")
	} else {
		fmt.Println("Ação inválida. Use 'encrypt' ou 'decrypt'.")
	}
}

// Gera uma chave segura a partir da senha do usuário
func generateKey(password string) []byte {
	salt := []byte("s@ltForExtraSecurity!") // Pode ser salvo ou gerado dinamicamente
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// Bloqueia o acesso ao diretório
func lockDirectory(directory string) error {
	err := os.Chmod(directory, 0000) // Remove permissões no Linux
	if err != nil {
		return err
	}
	if isWindows() {
		return hideDirectory(directory) // Oculta no Windows
	}
	return nil
}

// Desbloqueia o acesso ao diretório
func unlockDirectory(directory string) error {
	err := os.Chmod(directory, 0755) // Restaura permissões no Linux
	if err != nil {
		return err
	}
	if isWindows() {
		return unhideDirectory(directory) // Torna visível no Windows
	}
	return nil
}

// Oculta o diretório no Windows
func hideDirectory(directory string) error {
	attr := uint32(2)
	return syscall.SetFileAttributes(syscall.StringToUTF16Ptr(directory), attr)
}

// Torna o diretório visível no Windows
func unhideDirectory(directory string) error {
	attr := uint32(0)
	return syscall.SetFileAttributes(syscall.StringToUTF16Ptr(directory), attr)
}

// Verifica se está rodando no Windows
func isWindows() bool {
	return os.PathSeparator == '\\'
}
