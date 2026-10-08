# GERADOR DE PLANILHAS — GO + WAILS

## 1. OBJETIVO

Desenvolva um aplicativo desktop profissional chamado **Planly** — Gerador de Planilhas, utilizando:

* **Go** como backend
* **Wails** como framework desktop
* **React + TypeScript** para frontend
* **Vite** para build do frontend
* **Tailwind CSS** para estilização
* **shadcn/ui** para componentes
* **ExcelJS** ou biblioteca equivalente para geração/manipulação de planilhas Excel
* Aplicação desktop Windows inicialmente
* Arquitetura preparada para futura expansão

O aplicativo deve permitir que o usuário **crie, configure, visualize, editar e exportar planilhas Excel** através de uma interface gráfica moderna, sem precisar abrir o Excel para montar a estrutura.

A aplicação deve funcionar localmente, sem necessidade de internet ou servidor externo.

---

# 2. PRINCÍPIO DO SISTEMA

O usuário deverá conseguir criar uma planilha através de um fluxo visual.

Exemplo:

1. Criar novo projeto
2. Informar nome da planilha
3. Criar colunas
4. Definir tipos dos dados
5. Inserir dados
6. Configurar formatação
7. Visualizar a planilha
8. Gerar arquivo `.xlsx`
9. Escolher onde salvar
10. Abrir a pasta ou arquivo gerado

O sistema deve ser extremamente simples para usuários que não possuem conhecimento avançado de Excel.

---

# 3. INTERFACE

Criar uma interface moderna, profissional e limpa.

Idioma:

**Português do Brasil (pt-BR)**

Não utilizar textos em inglês na interface.

### Layout principal

Utilizar uma estrutura semelhante a:

```text
┌─────────────────────────────────────────────────────────────┐
│ PLANLY                              Novo   Abrir   Exportar │
├──────────────┬──────────────────────────────────────────────┤
│              │                                              │
│ PROJETO      │              ÁREA DE TRABALHO                │
│              │                                              │
│ 📄 Planilha  │                                              │
│ 📊 Dados     │             Tabela / Planilha               │
│ 🎨 Estilo    │                                              │
│ ⚙ Config.    │                                              │
│              │                                              │
│              │                                              │
├──────────────┴──────────────────────────────────────────────┤
│ Linhas: 150    Colunas: 12                 Pronto           │
└─────────────────────────────────────────────────────────────┘
```

---

# 4. DASHBOARD INICIAL

Ao abrir o aplicativo, apresentar uma tela inicial com:

### Ações principais

* Nova planilha
* Abrir projeto
* Abrir planilha existente
* Planilhas recentes

### Informações

Mostrar:

* Últimos projetos
* Quantidade de planilhas criadas
* Último arquivo utilizado
* Data da última modificação

### Botão principal

**+ Criar nova planilha**

---

# 5. CRIAÇÃO DE PLANILHA

Ao clicar em "Nova planilha", abrir um assistente.

## Etapa 1 — Informações

Campos:

* Nome da planilha
* Nome do arquivo
* Nome da aba
* Descrição

Exemplo:

```text
Nome:
Controle de Estoque

Arquivo:
controle_estoque.xlsx

Aba:
Estoque
```

---

# 6. DEFINIÇÃO DAS COLUNAS

Permitir adicionar, remover e reorganizar colunas.

Cada coluna deverá possuir:

* Nome
* Tipo
* Largura
* Obrigatório
* Valor padrão
* Formatação

Tipos disponíveis:

* Texto
* Número inteiro
* Número decimal
* Moeda
* Porcentagem
* Data
* Data e hora
* Booleano
* CPF
* CNPJ
* Telefone
* E-mail
* Código
* Fórmula

Exemplo:

```text
┌──────────────┬──────────────┬──────────┬───────────┐
│ Coluna       │ Tipo         │ Obrig.   │ Largura   │
├──────────────┼──────────────┼──────────┼───────────┤
│ Código       │ Código       │ Sim      │ 100       │
│ Produto      │ Texto        │ Sim      │ 250       │
│ Quantidade   │ Inteiro      │ Sim      │ 120       │
│ Valor        │ Moeda        │ Não      │ 120       │
│ Total        │ Fórmula      │ Não      │ 150       │
└──────────────┴──────────────┴──────────┴───────────┘
```

Permitir:

* Arrastar para reorganizar
* Duplicar coluna
* Editar coluna
* Excluir coluna

---

# 7. EDITOR DE PLANILHA

Criar um editor semelhante a uma planilha.

Exemplo:

```text
     A          B             C          D
1 | Código | Produto       | Quantidade | Valor
2 | 001    | Arroz         | 10         | 25,90
3 | 002    | Feijão        | 20         | 8,50
4 | 003    | Açúcar        | 15         | 4,90
```

Funcionalidades:

* Editar células
* Selecionar células
* Selecionar linha
* Selecionar coluna
* Copiar
* Colar
* Recortar
* Desfazer
* Refazer
* Adicionar linha
* Excluir linha
* Adicionar coluna
* Excluir coluna
* Redimensionar coluna
* Redimensionar linha
* Ordenar
* Filtrar
* Pesquisar

---

# 8. FORMATAÇÃO

Criar painel de formatação.

Permitir configurar:

### Texto

* Fonte
* Tamanho
* Negrito
* Itálico
* Sublinhado
* Alinhamento
* Cor do texto

### Célula

* Cor de fundo
* Bordas
* Alinhamento vertical
* Alinhamento horizontal
* Quebra de texto

### Número

* Número
* Moeda
* Porcentagem
* Data
* Data e hora

---

# 9. CABEÇALHO

Permitir configurar o cabeçalho da planilha.

Opções:

* Nome
* Fonte
* Tamanho
* Negrito
* Cor
* Fundo
* Alinhamento
* Bordas

Permitir definir automaticamente a primeira linha como cabeçalho.

---

# 10. FÓRMULAS

Implementar suporte a fórmulas.

Exemplos:

```text
=B2*C2

=SOMA(D2:D20)

=MÉDIA(D2:D20)

=MÁXIMO(D2:D20)

=MÍNIMO(D2:D20)
```

O sistema deve preservar as fórmulas no arquivo `.xlsx`.

Criar botão:

**Inserir fórmula**

E disponibilizar funções comuns:

* SOMA
* MÉDIA
* MÍNIMO
* MÁXIMO
* CONT.SE
* SE
* PROCV
* SOMASE
* ARRED
* HOJE

---

# 11. MODELOS PRONTOS

Criar biblioteca de modelos.

Exemplos:

### Estoque

Colunas:

* Código
* Produto
* Categoria
* Quantidade
* Estoque mínimo
* Valor unitário
* Valor total

### Vendas

* Data
* Cliente
* Produto
* Quantidade
* Valor unitário
* Desconto
* Total

### Clientes

* Código
* Nome
* CPF/CNPJ
* Telefone
* E-mail
* Cidade
* Estado

### Financeiro

* Data
* Descrição
* Categoria
* Tipo
* Valor
* Status

### Funcionários

* Código
* Nome
* CPF
* Cargo
* Salário
* Data de admissão

### Controle de produtos

* Código
* Produto
* Categoria
* Unidade
* Quantidade
* Valor
* Fornecedor

---

# 12. IMPORTAÇÃO

Permitir importar:

* `.xlsx`
* `.xls`
* `.csv`

Ao importar, detectar automaticamente:

* Cabeçalho
* Colunas
* Tipos de dados
* Quantidade de linhas

Mostrar uma prévia antes de confirmar.

---

# 13. EXPORTAÇÃO

Principal formato:

```text
.xlsx
```

Também permitir:

```text
.csv
```

Futuramente:

```text
.pdf
```

Ao exportar:

1. Validar dados
2. Gerar arquivo
3. Mostrar progresso
4. Exibir resultado

Exemplo:

```text
Planilha gerada com sucesso!

Arquivo:
controle_estoque.xlsx

Linhas: 1.250
Colunas: 8

[ Abrir arquivo ]
[ Abrir pasta ]
[ Fechar ]
```

---

# 14. CONFIGURAÇÕES DO EXCEL

Permitir configurar:

* Nome da aba
* Congelar primeira linha
* Congelar primeira coluna
* AutoFiltro
* Formatação automática
* Largura automática das colunas
* Altura das linhas
* Cabeçalho
* Rodapé
* Número de páginas
* Orientação da página
* Margens
* Área de impressão

---

# 15. VALIDAÇÃO

Antes de gerar a planilha, verificar:

* Colunas sem nome
* Campos obrigatórios vazios
* Valores inválidos
* Datas inválidas
* Fórmulas inválidas
* Tipos incompatíveis

Mostrar erros de forma clara.

Exemplo:

```text
⚠ Existem problemas na planilha

3 erros encontrados:

Linha 15:
Quantidade inválida.

Linha 27:
Data inválida.

Linha 42:
Campo "Produto" obrigatório.
```

Permitir clicar no erro e ir diretamente para a célula.

---

# 16. PROJETOS

O usuário poderá salvar projetos.

Extensão:

```text
.planly
```

Um projeto deve armazenar:

* Estrutura
* Colunas
* Dados
* Formatação
* Fórmulas
* Configurações
* Modelos utilizados

Permitir:

* Novo
* Abrir
* Salvar
* Salvar como
* Duplicar
* Renomear
* Excluir

---

# 17. AUTOSAVE

Implementar salvamento automático.

Exemplo:

```text
Salvamento automático ativado
Último salvamento: agora
```

Criar recuperação automática caso o aplicativo seja fechado inesperadamente.

---

# 18. BACKEND GO

O backend deve ser responsável por:

* Gerenciamento dos projetos
* Manipulação dos arquivos
* Geração do Excel
* Importação
* Exportação
* Validação
* Processamento de dados
* Sistema de templates
* Configurações
* Autosave

Utilizar arquitetura organizada.

Exemplo:

```text
planly/
│
├── backend/
│   ├── app.go
│   ├── models/
│   ├── services/
│   ├── excel/
│   ├── projects/
│   ├── templates/
│   ├── validation/
│   └── utils/
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── layouts/
│   │   ├── hooks/
│   │   ├── services/
│   │   ├── stores/
│   │   └── types/
│   │
│   └── ...
│
├── build/
├── wails.json
├── go.mod
└── README.md
```

---

# 19. COMUNICAÇÃO GO ↔ FRONTEND

Utilizar os métodos do Wails para comunicação entre Go e React.

Exemplo conceitual:

```go
func CreateProject(...)
func OpenProject(...)
func SaveProject(...)
func ImportSpreadsheet(...)
func ExportSpreadsheet(...)
func GenerateSpreadsheet(...)
func ValidateSpreadsheet(...)
func GetTemplates(...)
```

O frontend não deve manipular diretamente arquivos do sistema operacional.

Operações de arquivo devem passar pelo backend Go.

---

# 20. SISTEMA DE ARQUIVOS

Utilizar os diálogos nativos do sistema operacional através do Wails.

Exemplos:

**Abrir arquivo**

```text
Selecionar arquivo
```

**Salvar arquivo**

```text
Salvar planilha como...
```

Nunca exigir que o usuário digite manualmente caminhos de arquivos.

---

# 21. PERFORMANCE

O aplicativo deve ser preparado para planilhas grandes.

Objetivo inicial:

* 100.000+ linhas
* 100+ colunas

Evitar renderizar milhares de elementos DOM simultaneamente.

Utilizar virtualização da tabela quando necessário.

Operações pesadas devem ser executadas no Go para não travar a interface.

---

# 22. THREADS E PROCESSAMENTO

Operações demoradas como:

* Importação
* Exportação
* Geração
* Validação
* Processamento de milhares de linhas

devem ocorrer de forma assíncrona.

A interface deve continuar responsiva.

Mostrar:

```text
Processando...

████████████████░░░░ 80%

8.000 / 10.000 linhas
```

---

# 23. BANCO DE DADOS

Não utilizar banco de dados externo inicialmente.

Os projetos devem ser armazenados localmente.

Caso seja necessário persistir configurações do aplicativo, utilizar SQLite ou arquivo JSON local.

Arquitetar o sistema de forma que um banco possa ser adicionado futuramente.

---

# 24. TEMA

Implementar:

### Tema claro

Interface limpa e profissional.

### Tema escuro

Interface escura moderna.

Permitir:

```text
Sistema
Claro
Escuro
```

---

# 25. ATALHOS

Implementar atalhos:

```text
Ctrl + N    Novo projeto
Ctrl + O    Abrir projeto
Ctrl + S    Salvar
Ctrl + Shift + S    Salvar como
Ctrl + Z    Desfazer
Ctrl + Y    Refazer
Ctrl + C    Copiar
Ctrl + V    Colar
Ctrl + X    Recortar
Ctrl + F    Pesquisar
Ctrl + E    Exportar
```

---

# 26. MENU

Criar menu:

```text
Arquivo
├── Novo
├── Abrir
├── Salvar
├── Salvar como
├── Importar
├── Exportar
└── Sair

Editar
├── Desfazer
├── Refazer
├── Copiar
├── Colar
└── Recortar

Planilha
├── Nova linha
├── Nova coluna
├── Formatação
├── Ordenar
├── Filtrar
└── Validar

Exibir
├── Tema
├── Zoom
└── Painéis

Ajuda
├── Atalhos
├── Sobre
└── Verificar atualizações
```

---

# 27. ARQUITETURA

Utilizar separação clara de responsabilidades.

### Go

```text
Application
 ├── ProjectService
 ├── SpreadsheetService
 ├── ExcelService
 ├── ImportService
 ├── ExportService
 ├── ValidationService
 └── TemplateService
```

### React

```text
Pages
 ↓
Components
 ↓
Hooks
 ↓
Stores
 ↓
Wails API
 ↓
Go Services
```

Não colocar regras de negócio diretamente nos componentes React.

---

# 28. TRATAMENTO DE ERROS

Nunca deixar erros técnicos aparecerem diretamente ao usuário.

Errado:

```text
panic: runtime error...
```

Correto:

```text
Não foi possível gerar a planilha.

Verifique os dados e tente novamente.
```

Registrar detalhes técnicos em logs locais.

---

# 29. LOGS

Criar sistema de logs.

Diretório:

```text
%APPDATA%/Planly/logs/
```

Registrar:

* Inicialização
* Erros
* Importações
* Exportações
* Falhas
* Recuperações

Nunca registrar dados sensíveis das planilhas sem necessidade.

---

# 30. SEGURANÇA

* Validar todos os arquivos importados
* Evitar path traversal
* Validar extensões
* Evitar execução de conteúdo externo
* Não executar macros automaticamente
* Não permitir que fórmulas importadas executem código externo
* Sanitizar nomes de arquivos
* Tratar arquivos corrompidos corretamente

---

# 31. UX

A interface deve priorizar simplicidade.

O usuário deve conseguir:

**Abrir → Criar → Preencher → Exportar**

sem precisar conhecer programação.

Utilizar:

* Tooltips
* Mensagens claras
* Estados vazios
* Feedback visual
* Loading
* Toasts
* Confirmações quando necessário

---

# 32. ESTADOS VAZIOS

Exemplo:

```text
Nenhuma planilha criada

Crie sua primeira planilha para começar.

[ + Criar planilha ]
```

---

# 33. PRIMEIRA EXECUÇÃO

Na primeira abertura:

```text
Bem-vindo ao Planly!

Crie planilhas profissionais de maneira
rápida e simples.

[ Criar minha primeira planilha ]
```

Não exigir login.

Não exigir cadastro.

Não exigir internet.

---

# 34. REQUISITOS TÉCNICOS

Utilizar versões estáveis e atuais das tecnologias no momento do desenvolvimento.

Não utilizar tecnologias desnecessárias.

Evitar:

* Docker
* Servidor externo
* API externa
* Sistema de login
* Dependências desnecessárias

O aplicativo deve ser realmente desktop.

---

# 35. BUILD

Preparar o projeto para gerar:

```text
Planly.exe
```

e posteriormente um instalador:

```text
Planly-Setup.exe
```

O processo de build deve ser documentado.

---

# 36. README

Criar README.md contendo:

* Descrição
* Tecnologias
* Requisitos
* Instalação
* Desenvolvimento
* Build
* Estrutura
* Comandos
* Solução de problemas

---

# 37. QUALIDADE DO CÓDIGO

O código deve ser:

* Limpo
* Modular
* Tipado
* Documentado quando necessário
* Fácil de manter
* Fácil de expandir

Não criar arquivos gigantes.

Não duplicar lógica.

Não utilizar código fictício ou funções vazias para simular funcionalidades.

Todas as funcionalidades apresentadas na interface devem possuir implementação real.

---

# 38. ORDEM DE DESENVOLVIMENTO

Desenvolver em etapas.

### FASE 1

* Inicialização Wails
* Go
* React
* TypeScript
* Tailwind
* shadcn/ui
* Layout
* Dashboard

### FASE 2

* Criar projeto
* Estrutura da planilha
* Editor de células
* Linhas
* Colunas

### FASE 3

* Formatação
* Tipos
* Fórmulas
* Ordenação
* Filtros

### FASE 4

* Importação Excel
* Importação CSV
* Exportação Excel
* Exportação CSV

### FASE 5

* Projetos `.planly`
* Autosave
* Recuperação
* Templates

### FASE 6

* Performance
* Virtualização
* Validação
* Logs
* Tratamento de erros

### FASE 7

* Polimento visual
* Atalhos
* Tema
* Testes
* Build Windows

---

# 39. TESTES

Criar testes para:

### Go

* Serviços
* Geração Excel
* Importação
* Exportação
* Validação
* Projetos

### Frontend

* Componentes
* Editor
* Formulários
* Estados
* Interações

Criar testes para casos como:

* Planilha vazia
* 1 linha
* 100.000 linhas
* Colunas vazias
* Dados inválidos
* Arquivo corrompido
* Fórmulas inválidas
* Exportação com sucesso
* Exportação interrompida

---

# 40. RESULTADO FINAL

Ao finalizar, entregar um aplicativo desktop funcional chamado:

# PLANLY

O resultado deve permitir que um usuário comum abra o aplicativo e consiga criar uma planilha Excel profissional sem precisar utilizar diretamente o Excel.

O sistema deve possuir uma experiência semelhante a um pequeno editor de planilhas, mas com foco em **criação rápida, geração automática e modelos prontos**.

Priorize:

1. Usabilidade
2. Performance
3. Interface profissional
4. Geração correta do Excel
5. Arquitetura limpa
6. Facilidade de manutenção
7. Preparação para futuras funcionalidades

Não apenas crie um protótipo visual.

**Implemente as funcionalidades de forma real e funcional.**

Ao terminar cada fase, valide o funcionamento antes de prosseguir para a próxima.
