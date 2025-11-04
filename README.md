# <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"/> Go CEP 🔎

<img src="go-cep.jpg" alt="Golang" width="200" />

Aplicação simples em **Go** que expõe uma página HTML com um formulário de CEP e consome a API pública do **[ViaCEP](https://viacep.com.br/)** para preencher automaticamente endereço, bairro, cidade e UF.

- Servidor HTTP em Go
- Template HTML + JavaScript
- Consumo de API externa (ViaCEP) do lado do backend
- JSON de ida e volta
- UX básica de auto-preenchimento por CEP

---

## ⚙️ Tecnologias utilizadas

- **Go (Golang)** – `net/http`, `html/template`, `encoding/json`, `regexp`
- **ViaCEP** – API pública de CEP para o Brasil
- **HTML + JavaScript** – formulário e chamada assíncrona à API interna `/api/cep`

---

## 🚀 Funcionalidades

- Servidor HTTP em Go rodando em `http://localhost:9090`
- Página com formulário contendo:
  - Campo de **CEP**
  - Campos de **logradouro**, **bairro**, **cidade**, **UF**
- Ao digitar o CEP e:
  - clicar em **“Buscar CEP”**, ou
  - sair do campo (evento `blur`),
  
  o frontend faz uma requisição `fetch` para:

  ```text
  GET /api/cep?cep={CEP_LIMPO}
  Host: localhost:9090
