Perfeito. Então vamos seguir **exatamente no ciclo TDD** do material, na ordem certa:

1. **Red** → escrever um teste que falha.
2. **Green** → escrever o mínimo de código para o teste passar.
3. **Refactor** → melhorar o código sem quebrar os testes.

E no seu caso, o capítulo ainda mostra um detalhe importante: às vezes o “mínimo para passar” revela que o teste está incompleto. Isso faz parte do TDD.

Vou te orientar como se estivéssemos escrevendo juntos, arquivo por arquivo.

---

# Etapa 1 — criar o caminho feliz primeiro

Antes de testar cancelamento, precisamos provar que o servidor consegue buscar e devolver dados normalmente.

Comece criando algo assim:

```text
contexto/
├── context.go
├── context_test.go
└── go.mod
```

No `context.go`, coloque só o necessário:

```go
package contexto

import (
	"fmt"
	"net/http"
)

type Store interface {
	Fetch() string
}

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, store.Fetch())
	}
}
```

Agora vem o TDD de verdade.

## RED — escrever o teste

No `context_test.go`:

```go
package contexto

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubStore struct {
	response string
}

func (s *StubStore) Fetch() string {
	return s.response
}

func TestServer(t *testing.T) {
	data := "olá, mundo"

	store := &StubStore{
		response: data,
	}

	server := Server(store)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Body.String() != data {
		t.Errorf(
			`resultado "%s", esperado "%s"`,
			response.Body.String(),
			data,
		)
	}
}
```

Esse teste representa o comportamento:

```text
quando eu faço uma requisição
        ↓
o Server chama a Store
        ↓
a Store retorna "olá, mundo"
        ↓
o Server escreve "olá, mundo"
```

O material começa exatamente com esse caminho feliz antes de introduzir cancelamento. :chatgpt-content-reference{index="0"}

Execute:

```bash
go test
```

ou:

```bash
go test -v
```

Se passar, beleza.

---

# Etapa 2 — agora nasce um novo requisito

Agora queremos:

> Se a requisição for cancelada, a Store também deve parar o trabalho.

Esse é um comportamento novo.

Então começa um novo ciclo TDD.

## RED — primeiro mudamos nosso teste

Agora a `Store` precisa ter uma forma de ser cancelada.

Então alteramos:

```go
type Store interface {
	Fetch() string
	Cancel()
}
```

Nesse momento seu código provavelmente quebra.

E isso é normal.

Nosso `StubStore` não implementa mais a interface porque está faltando:

```go
Cancel()
```

Então transformamos ele em um **SpyStore**.

Por quê?

Porque agora queremos observar se alguma coisa aconteceu.

```go
type SpyStore struct {
	response  string
	cancelled bool
}
```

Implementamos:

```go
func (s *SpyStore) Fetch() string {
	return s.response
}
```

E:

```go
func (s *SpyStore) Cancel() {
	s.cancelled = true
}
```

Agora temos uma Store que consegue registrar:

```text
Cancel foi chamado?

true / false
```

---

# Etapa 3 — simular uma operação lenta

Precisamos que `Fetch()` demore.

Senão não daria tempo de cancelar a requisição.

Altere:

```go
func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)

	return s.response
}
```

Importe:

```go
import "time"
```

Mentalmente:

```text
Fetch começa
 ↓
espera 100ms
 ↓
retorna dados
```

---

# Etapa 4 — escrever o teste de cancelamento

Agora colocamos dois casos dentro de `TestServer`.

Primeiro:

```go
func TestServer(t *testing.T) {
	data := "olá, mundo"

	t.Run("retorna dados da store", func(t *testing.T) {

		store := &SpyStore{
			response: data,
		}

		server := Server(store)

		request := httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(
				`resultado "%s", esperado "%s"`,
				response.Body.String(),
				data,
			)
		}
	})
```

Agora o segundo:

```go
t.Run("cancela a store se a requisição for cancelada", func(t *testing.T) {

	store := &SpyStore{
		response: data,
	}

	server := Server(store)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	cancellingCtx, cancel := context.WithCancel(
		request.Context(),
	)

	time.AfterFunc(
		5*time.Millisecond,
		cancel,
	)

	request = request.WithContext(cancellingCtx)

	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if !store.cancelled {
		t.Errorf("store não foi avisada para cancelar")
	}
})
```

Esse teste faz exatamente isso:

```text
requisição começa
       ↓
Fetch demora 100ms

depois de 5ms
       ↓
cancel()
       ↓
context cancelado

esperado:
       ↓
store.Cancel()
```

O capítulo usa justamente `context.WithCancel`, agenda o `cancel` após 5 ms e injeta esse novo contexto na requisição. :chatgpt-content-reference{index="1"}

Agora execute:

```bash
go test -v
```

O teste deve falhar.

Algo parecido com:

```text
store não foi avisada para cancelar
```

Isso é o:

# RED

E isso é bom.

Você acabou de provar:

> “Meu sistema ainda não possui esse comportamento.”

---

# Etapa 5 — GREEN: fazer o mínimo possível

Agora vem uma parte muito importante do capítulo.

O TDD diz:

> escreva a menor quantidade de código possível para fazer o teste passar.

Então poderíamos fazer essa implementação absurda:

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		store.Cancel()

		fmt.Fprint(w, store.Fetch())
	}
}
```

Rode:

```bash
go test
```

O teste de cancelamento passa.

Por quê?

Porque:

```go
store.Cancel()
```

sempre acontece.

Só que isso criou um problema.

Mesmo uma requisição normal faz:

```text
Cancel
↓
Fetch
```

O que não faz sentido.

E aqui está um dos pontos mais legais desse capítulo: a implementação mínima fez o teste passar, mas revelou que **nosso conjunto de testes ainda não descrevia todo o comportamento esperado**. :chatgpt-content-reference{index="2"}

---

# Etapa 6 — melhorar o teste antes do código

Agora precisamos adicionar outra regra:

> Em uma requisição normal, a Store NÃO deve ser cancelada.

Dentro do teste feliz:

```go
t.Run("retorna dados da store", func(t *testing.T) {

	store := &SpyStore{
		response: data,
	}

	server := Server(store)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Body.String() != data {
		t.Errorf(
			`resultado "%s", esperado "%s"`,
			response.Body.String(),
			data,
		)
	}

	if store.cancelled {
		t.Error("não deveria ter cancelado a store")
	}
})
```

Agora execute:

```bash
go test
```

O caminho feliz deve falhar.

Por quê?

Porque nossa implementação atual faz:

```go
store.Cancel()
```

sempre.

Voltamos para:

# RED

---

# Etapa 7 — GREEN novamente

Agora precisamos implementar uma lógica verdadeira.

Primeiro pegamos:

```go
ctx := r.Context()
```

Depois criamos um channel para receber o resultado:

```go
data := make(chan string, 1)
```

Depois executamos `Fetch()` em paralelo:

```go
go func() {
	data <- store.Fetch()
}()
```

E agora usamos `select`:

```go
select {
case d := <-data:
	fmt.Fprint(w, d)

case <-ctx.Done():
	store.Cancel()
}
```

A função completa:

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {

		case d := <-data:
			fmt.Fprint(w, d)

		case <-ctx.Done():
			store.Cancel()
		}
	}
}
```

Execute:

```bash
go test -v
```

Agora esperamos:

```text
PASS
```

O comportamento agora é:

```text
              Server
                 |
        -------------------
        |                 |
        v                 v
    Fetch()          ctx.Done()
        |                 |
        |                 |
       dados           cancelou
        |                 |
        v                 v
   responde        store.Cancel()
```

O `select` está esperando duas coisas e executa o caso que estiver pronto primeiro. Essa é exatamente a implementação intermediária mostrada no material. :chatgpt-content-reference{index="3"}

---

# Etapa 8 — REFACTOR

Agora os testes passaram.

Então:

```text
RED ✅

GREEN ✅

agora:

REFACTOR
```

O capítulo primeiro começa melhorando os testes.

Em vez de escrever:

```go
if !store.cancelled {
	t.Errorf("store não foi avisada para cancelar")
}
```

podemos colocar isso dentro do próprio Spy.

O `SpyStore` vira:

```go
type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}
```

Criamos:

```go
func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()

	if !s.cancelled {
		s.t.Errorf("store não foi avisada para cancelar")
	}
}
```

E:

```go
func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()

	if s.cancelled {
		s.t.Errorf("store foi avisada para cancelar")
	}
}
```

Então o teste fica mais legível:

```go
store.assertWasCancelled()
```

e:

```go
store.assertWasNotCancelled()
```

Esse é o primeiro refactor sugerido no texto. :chatgpt-content-reference{index="4"}

---

# Mas o capítulo NÃO termina aqui

Agora vem a segunda parte importante.

O código funciona:

```go
case <-ctx.Done():
	store.Cancel()
```

Mas surge uma pergunta:

> O `Server` deveria realmente saber como cancelar a `Store`?

Imagine:

```text
Server
 ↓
Store
 ↓
Database
 ↓
API
 ↓
outro serviço
```

Nosso servidor teria que saber cancelar tudo?

Não.

A solução mais idiomática é fazer o `context` viajar pela aplicação:

```text
HTTP
 ↓
Server
 ↓ ctx
Store
 ↓ ctx
Database
 ↓ ctx
API
```

É por isso que o capítulo posteriormente muda a interface para:

```go
type Store interface {
	Fetch(ctx context.Context) (string, error)
}
```

:chatgpt-content-reference{index="5"}

E então o `Server` eventualmente fica muito mais simples:

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := store.Fetch(r.Context())

		if err != nil {
			return
		}

		fmt.Fprint(w, data)
	}
}
```

Essa é a refatoração final da arquitetura apresentada no capítulo. :chatgpt-content-reference{index="6"}

---

# A ordem correta para você escrever

Se você está acompanhando o livro na prática, faça **nessa ordem**:

```text
1. Store com Fetch()

2. Server chamando Fetch()

3. StubStore

4. teste:
   "retorna dados da store"

5. roda teste
   ✅ passa


6. novo requisito:
   cancelamento

7. muda Store:
   Fetch()
   Cancel()

8. transforma StubStore em SpyStore

9. adiciona:
   cancelled bool

10. Fetch demora 100ms

11. cria teste:
    "cancela se requisição cancelar"

12. roda teste
    ❌ RED


13. implementação mínima:
    store.Cancel()

14. roda teste
    ✅ novo teste passa


15. percebe falha de design

16. adiciona no caminho feliz:
    store NÃO deve cancelar

17. roda teste
    ❌ RED


18. implementa:
    ctx
    goroutine
    channel
    select

19. roda teste
    ✅ GREEN


20. refatora os asserts


21. questiona arquitetura


22. refatora Store para:

    Fetch(ctx context.Context)
       (string, error)


23. modifica SpyStore


24. altera testes


25. implementa Server:

    data, err :=
        store.Fetch(r.Context())


26. roda testes

27. ✅ GREEN
```

Esse é o fluxo do capítulo.

### E tem uma coisa importante

Você **não deveria simplesmente copiar a solução final**:

```go
Fetch(ctx context.Context)
```

desde o começo.

Porque aí você perde justamente a lição de TDD que o capítulo está tentando ensinar:

```text
problema
   ↓
teste
   ↓
falha
   ↓
mínimo necessário
   ↓
novo problema aparece
   ↓
teste melhora
   ↓
implementação melhora
   ↓
refactor
```

O objetivo não é somente chegar no código certo.

É aprender **como o código chegou naquela forma**.
