# AGENTS.md — Bora Vôlei

## Objetivo e estado atual

Construir uma aplicação para organizar jogos de vôlei com inscrições pagas por Pix. Leia README.md antes de alterar o projeto. Neste momento há somente documentação inicial: arquitetura, rotas e políticas são propostas. Não declare API, frontend, testes ou comandos de execução como existentes sem verificar o repositório.

## Direção técnica

- Backend em Go, PostgreSQL, API REST versionada em `/api/v1`.
- Monólito modular; API e worker compartilham regras, com tarefas duráveis no banco.
- React + TypeScript é provável para o frontend; confirmar antes de criar `web/`.
- Um único recebedor no MVP. Gateway ainda aberto; Efí é candidata provisória.
- Preferir biblioteca padrão e SQL explícito; pgx é a proposta de acesso ao banco. Não adicionar Redis, microserviços ou filas externas sem necessidade demonstrada.
- Usar versão estável suportada de Go, fixada em go.mod e ambiente. Consultar documentação oficial atual ao integrar dependências/gateway.

## Regras que não podem ser quebradas

1. Somente backend confirma pagamento, após validação no provedor. Frontend, comprovante ou webhook não verificado não são confirmação.
2. Monetários são inteiros em centavos e moeda explícita. Valor da cobrança vem da inscrição congelada, não do cliente.
3. Reservas ativas e confirmados nunca excedem capacidade. Reserva, expiração, confirmação e cancelamento usam transações e ordem consistente de locks (jogo antes da inscrição).
4. Não segurar lock de banco durante chamadas ao gateway.
5. Inscrição única por usuário/jogo; tentativas de pagamento distintas e rastreáveis. Sem cobranças simultâneas ou recriação cega após timeout.
6. Criação de cobrança, webhook, confirmação e devolução são idempotentes. Chaves repetidas com payload diferente geram conflito.
7. Persistir webhook/tarefa antes de reconhecer sucesso. Duplicação e ordem invertida de eventos não podem corromper estados.
8. Pagamento recebido após liberação da vaga exige nova verificação de capacidade ou devolução; nunca perder registro financeiro.
9. Devolução é pendente até confirmação externa. Falhas e retries precisam permanecer recuperáveis.
10. Cadastro público cria player; toda ação verifica papel/propriedade no servidor. Não permitir alteração de role pelo perfil.
11. Não excluir registros financeiros e histórico para implementar cancelamentos.
12. Datas em UTC no banco, com fuso IANA para exibição.

## Organização e contratos

Siga a estrutura proposta no README, adaptando somente quando houver motivo concreto. Handlers cuidam do transporte; serviços das regras; repositórios do SQL; adaptador do gateway da integração externa. Evite abstrações sem uso.

Antes de implementar endpoints, criar/atualizar `api/openapi.yaml` com schemas, autorização, erros, paginação e idempotência. Não alterar contratos silenciosamente. Migrations versionadas, revisadas e sem reset destrutivo de banco existente.

Gateway interno deve permitir criar cobrança, consultar pagamento e solicitar/consultar devolução; autenticação de webhook continua específica do provedor. Implementar um gateway real e fake para testes. Não inventar SDK, endpoint, status ou assinatura: validar documentação oficial.

## Qualidade, segurança e desempenho

- Propagar context.Context, limites de pool e timeouts; reutilizar cliente HTTP.
- Goroutines limitadas e supervisionadas; não substituir jobs duráveis por goroutines soltas.
- SQL parametrizado, paginação e índices orientados às consultas. Otimizar com medições.
- Segredos, certificados e arquivos .env reais fora do git. .env.example contém somente placeholders.
- Logs estruturados com request_id; não registrar tokens, senhas, CPF, payload Pix completo ou dados pessoais desnecessários.
- Aplicar mecanismo oficial de autenticação do webhook. Não desabilitar validação TLS/mTLS para facilitar produção.
- Não executar pagamentos reais, publicar serviços ou usar credenciais de produção sem autorização específica.
- Não criar títulos/tarifas de fornecedor como garantias. Distinguir cotações verificadas e exemplos hipotéticos.

## Testes e entrega

Testar invariantes financeiras e concorrência, especialmente última vaga, duplicação, timeout, webhook atrasado, cancelamento concorrente e devolução com retry. Usar PostgreSQL isolado em integração e fake/sandbox financeiro.

Quando houver módulo Go executável, rodar formatação (`gofmt` nos arquivos alterados), `go test ./...` e `go vet ./...`. Usar `go test -race ./...` para mudanças concorrentes quando o ambiente suportar. Não afirmar sucesso de checks não executados; registrar bloqueios. Para frontend, usar scripts realmente definidos no projeto.

Ao entregar uma alteração, explicar comportamento resultante, validação feita e pendências materiais. Atualizar README e decisões quando mudar política ou contrato. Não construir funcionalidades fora do MVP sem solicitação.
