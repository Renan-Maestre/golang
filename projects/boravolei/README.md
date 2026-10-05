# 🏐 Bora Vôlei

Documento inicial de produto e arquitetura — proposta para o MVP, sem implementação ainda.

## Ideia

O Bora Vôlei organiza jogos de vôlei e automatiza inscrições pagas. Um administrador publica o jogo com data, horário, local, limite de participantes e preço. Os usuários consultam os jogos, reservam uma vaga e pagam por Pix. Após confirmação do pagamento pelo backend, entram na lista de participantes. O aplicativo mantém histórico de jogos, inscrições e pagamentos.

## Escopo inicial

- Cadastro, login e perfil básico.
- Papéis `admin` e `player`; cadastro público sempre cria `player`.
- Administrador cria, publica, edita, cancela e conclui jogos.
- Usuários autenticados consultam jogos e participantes confirmados.
- Inscrição com reserva temporária, Pix dinâmico, QR Code e copia e cola.
- Confirmação automática por webhook e consulta de reconciliação.
- Histórico pessoal e painel administrativo de participantes e valores.
- Cancelamento e devolução com rastreabilidade.

Fora do MVP: chat, sorteio de times, ranking, recorrência, múltiplos recebedores, marketplace e aplicativo nativo. React é a opção provável para a interface; ainda não é uma decisão fechada.

## Premissas propostas

O MVP terá um único grupo organizador e uma única conta recebedora. Administradores do sistema administram os jogos desse grupo. Se cada organizador receber em sua própria conta, será necessário redesenhar permissões, onboarding e integração financeira antes de implementar.

O valor informado pelo administrador é o preço bruto da inscrição, em BRL. Exemplo: `1000` centavos = R$ 10,00. A tarifa do gateway reduz o valor líquido recebido. O app não adiciona taxa própria no MVP.

Reserva sugerida: 15 minutos, configurável e limitada pelo prazo de inscrição e pelo vencimento suportado pelo gateway. O preço fica congelado na inscrição. Após a primeira reserva, bloquear alterações de preço; alterações de local/data devem ser auditadas e informadas aos participantes.

Política inicial proposta: cancelamento pelo usuário até o prazo configurado devolve integralmente o valor da inscrição; cancelamento do jogo devolve todas as inscrições pagas. Confirmar essa política e eventual custo de devolução antes de produção.

## Opções de gateway Pix

Pesquisa em 05/10/2026. Tarifas dependem da conta, contratação e negociação; validar a proposta efetiva antes da escolha. Gratuidade de Pix por chave não implica gratuidade de cobrança dinâmica via API.

| Opção | Integração | Custo verificado | Avaliação para o projeto |
| --- | --- | --- | --- |
| Efí Bank | API Pix com OAuth2 e certificado cliente; configuração inicial mais trabalhosa | Tabela pública: 1,19% para QR Code dinâmico/Pix Cob via API, sujeito à contratação | Primeira candidata pelo baixo custo publicado para ticket de R$ 10 |
| Mercado Pago | Checkout API com geração de Pix; alternativa conhecida para o checkout | Tarifa atual de Pix via checkout não confirmada nesta pesquisa; consultar condições da conta | Comparar proposta efetiva com Efí; candidata se simplicidade compensar custo |
| Asaas | API REST de cobrança e obtenção de QR Code dinâmico | Tarifa dinâmica não confirmada; 100 Pix gratuitos publicados se referem a chave/QR estático | Candidata pela organização de cobranças; atenção a eventual taxa fixa por pagamento |

Na tarifa publicada da Efí, R$ 10 × 1,19% = R$ 0,119: aproximadamente R$ 0,12 por inscrição. Arredondamento e valor líquido devem seguir o extrato do provedor. Exemplo hipotético: tarifa fixa de R$ 1 representaria 10% de uma inscrição de R$ 10; esse exemplo não é uma cotação de nenhum fornecedor.

**Recomendação provisória:** testar Efí em homologação e comparar propostas de Mercado Pago e Asaas. Implementar apenas um gateway no MVP, atrás de uma interface interna. A decisão permanece aberta.

Fontes oficiais:

- Efí — tarifas: https://sejaefi.com.br/tarifas
- Efí — credenciais, certificados e autorização: https://dev.efipay.com.br/docs/api-pix/credenciais/
- Mercado Pago — Pix: https://www.mercadopago.com.br/developers/pt/docs/checkout-api-payments/integration-configuration/integrate-pix
- Asaas — Pix: https://docs.asaas.com/docs/pix
- Asaas — condições de gratuidade: https://central.ajuda.asaas.com/hc/pt-br/articles/32040230167067-Quais-s%C3%A3o-as-taxas-para-receber-via-Pix

## Arquitetura proposta

Monólito modular em Go com API HTTP e worker, PostgreSQL e frontend React futuro. API e worker podem compartilhar o mesmo código e rodar como processos separados. Não começar com microserviços, Redis ou broker externo sem necessidade demonstrada.

Módulos: autenticação, usuários, jogos, inscrições, pagamentos e auditoria. Handlers HTTP validam entrada; serviços executam regras; repositórios persistem; adaptador do gateway concentra chamadas externas.

Stack sugerida: Go em versão estável suportada a definir no bootstrap, `net/http`, PostgreSQL com `pgx`, SQL explícito e migrations versionadas. `sqlc` é opcional. Docker Compose para ambiente local. React + TypeScript com Vite é uma proposta para a etapa de frontend.

Persistir tarefas e eventos financeiros no PostgreSQL. Worker busca tarefas com lease/bloqueio, retries limitados, backoff e registro de falhas. Goroutines devem ter limites, contexto e encerramento controlado; não são mecanismo de durabilidade.

## Fluxo de inscrição e pagamento

1. Usuário autenticado solicita inscrição com `Idempotency-Key`.
2. Em transação curta, bloquear a linha do jogo, validar publicação, prazo e capacidade. Contar confirmados e reservas ativas. Criar inscrição/reserva, tentativa de pagamento e tarefa de criação da cobrança. Confirmar transação antes de acessar o gateway.
3. API retorna `202` enquanto a cobrança está sendo criada. Worker cria Pix com referência e chave estáveis, salva identificador externo, vencimento, QR Code e copia e cola. Frontend consulta a inscrição com intervalo e backoff.
4. Usuário paga. Webhook autenticado é persistido de forma durável e deduplicado antes de responder sucesso. Worker consulta o pagamento no provedor e valida recebedor, referência, moeda, valor e situação.
5. Em transação com bloqueio do jogo e da inscrição, registrar pagamento e confirmar a vaga uma única vez. A lista exibe somente inscrições confirmadas.
6. Worker libera reservas vencidas, reconcilia pagamentos pendentes e processa devoluções. O navegador nunca confirma pagamento.

Se a chamada de criação sofrer timeout, consultar a referência externa antes de repetir. Não abrir outra cobrança enquanto a tentativa anterior estiver em situação incerta.

**Pagamento tardio:** pagamento realizado dentro da validade, mas com webhook atrasado, precisa ser reconhecido na reconciliação. Caso a vaga já tenha sido liberada, confirmar somente se ainda existir capacidade e o jogo permitir entrada; caso contrário, marcar para devolução. Nunca ultrapassar a capacidade ou ignorar dinheiro recebido. Pagamentos de inscrições canceladas/jogos cancelados também seguem para devolução.

## Regras e estados

- Jogo: `draft`, `published`, `cancelled`, `completed`. Lotação é calculada, não um estado persistido.
- Inscrição: `pending_payment`, `confirmed`, `expired`, `cancelled`, `refund_pending`, `refunded`.
- Pagamento: `creating`, `pending`, `paid`, `failed`, `expired`, `refund_pending`, `refunded`.
- `failed` significa falha conhecida; timeout incerto exige reconciliação.
- Devolução solicitada não significa devolução concluída; concluir só após confirmação do gateway.
- Uma inscrição por usuário e jogo, com tentativas de pagamento distintas. Reativação de inscrição expirada/cancelada só sem pagamento pendente ou devolução em andamento e após nova verificação de capacidade.
- Cada pagamento externo e evento processado deve ter identificação única. Eventos duplicados ou fora de ordem não podem regredir um pagamento pago para pendente.
- Cancelamento libera a vaga em transação e registra a obrigação financeira quando houver pagamento. Não excluir histórico.
- Não reduzir capacidade abaixo de confirmados e reservas ativas. Não inscrever em jogo cancelado, concluído ou fora do prazo.
- Novas inscrições são bloqueadas no encerramento. Conclusão do jogo exige reconciliação das reservas; devoluções continuam processáveis depois disso.

## Modelo de dados inicial

| Entidade | Campos principais |
| --- | --- |
| `users` | id, name, email único, password_hash, role, created_at |
| `games` | id, created_by, title, description, venue_name, address, starts_at, ends_at, timezone, registration_deadline, cancellation_deadline, capacity, price_cents, currency, status |
| `registrations` | id, game_id, user_id, status, price_cents, currency, reserved_until, confirmed_at, cancelled_at, created_at |
| `payments` | id, registration_id, provider, provider_payment_id, attempt_number, idempotency_key, amount_cents, currency, status, expires_at, paid_at, refunded_at, provider_fee_cents |
| `webhook_events` | id, provider, deduplication_key, payment_reference, received_at, processing_status, attempts |
| `jobs` | id, kind, reference_id, status, run_at, attempts, lease_until, last_error |
| `audit_logs` | id, actor_id, action, entity_type, entity_id, created_at, metadata filtrada |

Valores monetários em inteiro de centavos (`int64`), nunca float. Datas em UTC (`timestamptz`) e fuso IANA para apresentação, inicialmente `America/Sao_Paulo`. QR Code não é permanente: guardar apenas o necessário com acesso restrito.

Constraints: unicidade `(game_id, user_id)`, `(provider, provider_payment_id)` quando preenchido, `(provider, deduplication_key)` e chave de idempotência por usuário/operação. Índices para jogos por status/data, inscrições por jogo/status e usuário/data, pagamentos por status e tarefas por status/run_at. Capacidade exige transação; índice único sozinho não resolve disputa da última vaga.

## Rotas propostas

Prefixo `/api/v1`. `player` inclui administrador; `admin` exige papel administrativo. Todas as operações de recursos pessoais validam propriedade no servidor.

| Método | Rota | Acesso | Finalidade |
| --- | --- | --- | --- |
| POST | `/auth/register` | Público | Cadastro como player |
| POST | `/auth/login` | Público | Login |
| POST | `/auth/refresh` | Sessão de renovação | Renovar acesso com rotação |
| POST | `/auth/logout` | Sessão | Revogar sessão |
| GET/PATCH | `/me` | Player | Ler/editar perfil permitido |
| GET | `/games` | Player | Jogos publicados com filtros e paginação |
| GET | `/games/{game_id}` | Player | Detalhes do jogo |
| GET | `/games/{game_id}/participants` | Player | Nomes de confirmados; sem contatos ou dados financeiros |
| POST | `/games/{game_id}/registrations` | Player | Reservar vaga e iniciar cobrança |
| GET | `/me/registrations` | Player | Histórico pessoal paginado |
| GET | `/registrations/{registration_id}` | Dono/admin | Estado e Pix disponível da inscrição |
| POST | `/registrations/{registration_id}/cancel` | Dono/admin | Cancelar conforme política |
| GET | `/payments/{payment_id}` | Dono/admin | Estado de pagamento |
| POST | `/webhooks/payments/{provider}` | Gateway autenticado | Receber eventos do provedor habilitado |
| GET | `/admin/games` | Admin | Jogos de todos os estados |
| POST | `/admin/games` | Admin | Criar rascunho |
| PATCH | `/admin/games/{game_id}` | Admin | Editar campos permitidos |
| POST | `/admin/games/{game_id}/publish` | Admin | Publicar |
| POST | `/admin/games/{game_id}/cancel` | Admin | Cancelar e programar devoluções |
| POST | `/admin/games/{game_id}/complete` | Admin | Concluir após jogo e reconciliação |
| GET | `/admin/games/{game_id}/registrations` | Admin | Inscritos e estados financeiros |
| GET | `/admin/games/{game_id}/summary` | Admin | Bruto, taxas conhecidas, devoluções e líquido |
| POST | `/admin/payments/{payment_id}/refund` | Admin | Solicitar devolução integral idempotente |
| GET | `/health/live` | Infraestrutura | Processo ativo |
| GET | `/health/ready` | Infraestrutura | Prontidão e banco disponível |

Não oferecer rota pública para marcar pagamento como pago. IDs de cobrança e valores enviados pelo cliente não são autoridade financeira.

### Exemplo de criação de jogo

`POST /api/v1/admin/games`

```json
{
  "title": "Vôlei de sábado",
  "venue_name": "Quadra Central",
  "address": "Endereço da quadra",
  "starts_at": "2026-10-10T18:00:00-03:00",
  "ends_at": "2026-10-10T20:00:00-03:00",
  "timezone": "America/Sao_Paulo",
  "registration_deadline": "2026-10-10T17:00:00-03:00",
  "cancellation_deadline": "2026-10-10T12:00:00-03:00",
  "capacity": 18,
  "price_cents": 1000,
  "currency": "BRL"
}
```

`POST /api/v1/games/{game_id}/registrations`, sem valor no body, com `Idempotency-Key`.

```json
{
  "id": "registration_uuid",
  "status": "pending_payment",
  "reserved_until": "2026-10-05T17:45:00Z",
  "payment": { "id": "payment_uuid", "status": "creating" }
}
```

A consulta posterior retorna `pix_copy_paste`, `qr_code` e `expires_at` somente ao dono/admin. Reuso da chave com payload diferente retorna `409`; repetições iguais retornam a mesma operação.

Convenções: JSON em snake_case; IDs UUID; `201` criação, `202` processamento, `400` entrada inválida, `401` sem sessão, `403` acesso negado, `404` recurso ausente, `409` conflito/lotação, `429` limite. Erros com `code`, `message` e `request_id`, sem stack trace. Listas com cursor e limite máximo; filtros documentados em OpenAPI antes da implementação.

## Segurança e desempenho

- Autorização no backend por papel e propriedade; impedir autopromoção no PATCH de perfil.
- Senhas com hash apropriado; estratégia de sessão a fechar antes do módulo auth. Refresh rotacionado, revogável e armazenado com hash. Para cookies, usar HttpOnly/Secure e proteção CSRF adequada.
- Credenciais/certificados do gateway só no backend, fora do git; usar ambientes separados e privilégios mínimos.
- Webhook protegido pelo mecanismo oficial do provedor (assinatura, token ou mTLS, conforme integração); endpoint genérico não elimina validação específica.
- Validar tamanho de body, entradas, origens CORS e limites de login/inscrição. SQL parametrizado.
- Não expor CPF, email, telefone ou payload financeiro na lista de participantes. Coletar dados de pagador apenas quando exigidos pela integração.
- Pool de conexões limitado, timeouts HTTP/SQL/gateway, `context.Context`, graceful shutdown e clientes HTTP reutilizados.
- Não manter locks SQL durante chamadas de rede. Paginar listas, evitar N+1 e medir planos de consultas antes de introduzir cache.
- Logs estruturados e métricas de latência p95/p99, erros, pool, fila de tarefas, atraso de confirmação e devoluções pendentes. Não registrar segredos ou payloads pessoais completos.
- Fazer backup e testar recuperação. Definir meta de latência e carga após medir infraestrutura e volume esperado; Go sozinho não garante desempenho.

## Estrutura proposta de repositório

| Caminho | Responsabilidade |
| --- | --- |
| `cmd/api/` | Inicialização do servidor |
| `cmd/worker/` | Inicialização de tarefas duráveis |
| `internal/auth/`, `internal/users/` | Identidade e acesso |
| `internal/games/`, `internal/registrations/` | Regras de jogos e vagas |
| `internal/payments/` | Regras financeiras e interface do gateway |
| `internal/payments/providers/` | Adaptador do provedor escolhido |
| `internal/platform/` | Configuração, banco e observabilidade |
| `migrations/` | Evolução do schema |
| `api/openapi.yaml` | Contrato a criar |
| `web/` | Frontend futuro |
| `docs/decisions/` | Decisões e justificativas |

Os caminhos são uma proposta; não indicam código já existente.

## Plano de implementação e validação

1. Fechar recebedor, gateway, autenticação, prazo da reserva e política de cancelamento.
2. Criar bootstrap Go, banco, migrations, configuração, health checks e contrato OpenAPI.
3. Implementar autenticação, autorização e CRUD/publicação de jogos.
4. Implementar reservas e concorrência com gateway fake.
5. Integrar Pix em homologação, webhook durável, reconciliação e devolução.
6. Implementar interface de jogos, checkout e histórico; painel administrativo.
7. Validar testes, backup, observabilidade e fluxo financeiro antes do piloto.

Testes prioritários: dois usuários disputando a última vaga; webhook duplicado/fora de ordem; timeout criando cobrança; pagamento após expiração; webhook atrasado de pagamento dentro do prazo; cancelamento durante pagamento; falha/retry de devolução; acesso indevido; alteração de preço; recuperação de worker após reinício. Testes de integração usam PostgreSQL real isolado. Nunca movimentar dinheiro real em testes automáticos.

## Decisões ainda abertas

Gateway e tarifa contratada; conta recebedora; necessidade futura de múltiplos grupos; prazo de reserva; regras finais de cancelamento; campos obrigatórios do pagador; método de autenticação e recuperação de senha; frontend definitivo; hospedagem e volume esperado. Registrar respostas como decisões e atualizar este documento e AGENTS.md.
