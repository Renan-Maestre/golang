<div align="center">

<img src="https://go.dev/images/go-logo-blue.svg" width="220" alt="Go Logo">

# 🐹 Aprendendo Go com Testes

### Estudos de Golang utilizando testes, TDD e boas práticas de desenvolvimento

![](https://2763912287-files.gitbook.io/~/files/v0/b/gitbook-legacy-files/o/assets%2F-Lia9CiG1cfWmh7Adpdu%2F-Lia9TbxTuAr7XbyNb3I%2F-Lia9ambvvlPGkPz-Q7f%2Fred-green-blue-gophers-smaller.png?generation=1561860928341453&alt=media)

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Tests](https://img.shields.io/badge/Tests-Go_Testing-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/testing)
[![TDD](https://img.shields.io/badge/TDD-Test_Driven_Development-success?style=for-the-badge)](https://en.wikipedia.org/wiki/Test-driven_development)
[![Learning](https://img.shields.io/badge/Status-Aprendendo-yellow?style=for-the-badge)](#)

</div>

---

## 📚 Sobre o projeto

Este repositório registra minha evolução no aprendizado da linguagem **Go (Golang)** utilizando uma abordagem baseada em testes.

O objetivo não é somente aprender a sintaxe da linguagem, mas entender como desenvolver aplicações de maneira organizada, simples, testável e confiável.

Grande parte dos estudos utiliza como referência o projeto:

**Learn Go with Tests**

- 🇧🇷 [Aprenda Go com Testes](https://larien.gitbook.io/aprenda-go-com-testes/)
- 🇺🇸 [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)
- 🧪 [Go by Example](https://gobyexample.com/)
- 📖 [Documentação oficial do Go](https://go.dev/doc/)
- 📦 [Go Packages](https://pkg.go.dev/)

---

## 🎯 Objetivos

Durante os estudos, pretendo desenvolver conhecimento sobre:

- 🐹 Fundamentos da linguagem Go
- 🧪 Testes unitários
- 🔴🟢🔵 TDD — Test Driven Development
- 📦 Packages e módulos
- 🧱 Structs
- 🔌 Interfaces
- 👉 Ponteiros
- 🗺️ Maps
- 📚 Arrays e Slices
- ⚠️ Tratamento de erros
- 🧩 Injeção de dependência
- 🎭 Mocks
- ⚡ Concorrência
- 🔀 Goroutines
- 📡 Channels
- ⏳ Context
- 🔒 Mutex
- 🔄 WaitGroup
- 🌐 APIs HTTP
- 📄 JSON
- 🔌 WebSockets
- 📁 Organização de projetos Go
- 📊 Benchmarks
- 📈 Cobertura de testes
- 🧹 Refatoração

---

# 🧪 Por que aprender Go utilizando testes?

Uma das características interessantes do Go é possuir suporte nativo para testes através do package:

```go
testing
```

Isso permite executar testes sem depender inicialmente de frameworks externos.

Um arquivo de teste utiliza normalmente o padrão:

```text
arquivo_test.go
```

Por exemplo:

```text
adicionar.go
adicionar_test.go
```

Podemos então executar:

```bash
go test
```

Ou visualizar mais informações:

```bash
go test -v
```

Também podemos testar todos os packages do projeto:

```bash
go test ./...
```

---

# 🔴 🟢 🔵 Ciclo TDD

O **TDD — Test Driven Development** trabalha normalmente com um ciclo simples:

```text
        ┌─────────────┐
        │ 🔴 RED      │
        │ Criar teste │
        └──────┬──────┘
               │
               ▼
        ┌─────────────┐
        │ 🟢 GREEN    │
        │ Fazer passar│
        └──────┬──────┘
               │
               ▼
        ┌─────────────┐
        │ 🔵 REFACTOR │
        │ Melhorar    │
        └──────┬──────┘
               │
               └───────────────► Repetir
```

### 🔴 RED

Primeiro escrevemos um teste que representa o comportamento desejado.

O teste inicialmente deve falhar.

### 🟢 GREEN

Implementamos somente o necessário para fazer o teste passar.

### 🔵 REFACTOR

Com os testes passando, podemos melhorar a implementação mantendo o comportamento existente.

Depois repetimos o processo.

---

# 💡 Exemplo simples

## Código

```go
package inteiros

func Adicionar(x, y int) int {
	return x + y
}
```

## Teste

```go
package inteiros

import "testing"

func TestAdicionar(t *testing.T) {
	resultado := Adicionar(2, 2)
	esperado := 4

	if resultado != esperado {
		t.Errorf(
			"resultado %d, esperado %d",
			resultado,
			esperado,
		)
	}
}
```

Executando:

```bash
go test
```

Resultado esperado:

```text
PASS
ok
```

---

# 🗂️ Estrutura de estudos

A ideia deste repositório é separar cada conceito estudado em seu próprio package/diretório.

Exemplo:

```text
.
├── ola-mundo/
│   ├── hello.go
│   └── hello_test.go
│
├── inteiros/
│   ├── adicionar.go
│   └── adicionar_test.go
│
├── iteracao/
│   ├── repetir.go
│   └── repetir_test.go
│
├── arrays-slices/
│   ├── soma.go
│   └── soma_test.go
│
├── structs/
│
├── ponteiros/
│
├── maps/
│
├── dependency-injection/
│
├── mocks/
│
├── concorrencia/
│
├── select/
│
├── sync/
│
├── context/
│
├── go.mod
└── README.md
```

A estrutura poderá evoluir conforme novos conceitos forem estudados.

---
