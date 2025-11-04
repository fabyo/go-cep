package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
)

// Struct da resposta do ViaCEP
type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Erro       bool   `json:"erro,omitempty"`
}

var tpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html lang="pt-BR">
<head>
	<meta charset="utf-8">
	<title>Form CEP em Go</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; }
		label { display: block; margin-top: 10px; }
		input { padding: 5px; width: 300px; }
		#status { margin-top: 10px; color: #555; font-size: 0.9rem; }
	</style>
</head>
<body>
	<h1>Formulário com CEP (Go + ViaCEP)</h1>

	<label>
		CEP:
		<input type="text" id="cep" placeholder="Digite o CEP" maxlength="9" />
		<button type="button" onclick="buscarCEP()">Buscar CEP</button>
	</label>

	<label>
		Logradouro:
		<input type="text" id="logradouro" />
	</label>

	<label>
		Bairro:
		<input type="text" id="bairro" />
	</label>

	<label>
		Cidade:
		<input type="text" id="cidade" />
	</label>

	<label>
		UF:
		<input type="text" id="uf" />
	</label>

	<div id="status"></div>

	<script>
	async function buscarCEP() {
		const cepInput = document.getElementById('cep');
		const statusEl = document.getElementById('status');

		let cep = cepInput.value.replace(/\D/g, '');

		if (cep.length !== 8) {
			statusEl.textContent = 'CEP inválido. Use 8 dígitos.';
			return;
		}

		statusEl.textContent = 'Buscando CEP...';

		try {
			const resp = await fetch('/api/cep?cep=' + cep);
			if (!resp.ok) {
				statusEl.textContent = 'Erro na requisição: ' + resp.status;
				return;
			}
			const data = await resp.json();

			if (data.error) {
				statusEl.textContent = 'Erro: ' + data.error;
				return;
			}

			document.getElementById('logradouro').value = data.logradouro || '';
			document.getElementById('bairro').value     = data.bairro || '';
			document.getElementById('cidade').value     = data.localidade || '';
			document.getElementById('uf').value         = data.uf || '';

			statusEl.textContent = 'CEP carregado com sucesso.';
		} catch (e) {
			console.error(e);
			statusEl.textContent = 'Erro inesperado ao buscar CEP.';
		}
	}

	// opcional: busca ao sair do campo CEP (blur)
	document.getElementById('cep').addEventListener('blur', function() {
		if (this.value.trim() !== '') {
			buscarCEP();
		}
	});
	</script>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/api/cep", cepHandler)

	url := "http://localhost:9090"
	log.Println("Servidor rodando em", url)

	// abre o navegador
	go openBrowser(url)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, "erro ao renderizar template", http.StatusInternalServerError)
	}
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	cep = onlyDigits(cep)

	if len(cep) != 8 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "CEP inválido, precisa ter 8 dígitos",
		})
		return
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Falha ao consultar ViaCEP",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"error": fmt.Sprintf("ViaCEP retornou status %d", resp.StatusCode),
		})
		return
	}

	var vc ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&vc); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Erro ao decodificar resposta do ViaCEP",
		})
		return
	}

	if vc.Erro {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"error": "CEP não encontrado",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"cep":        vc.Cep,
		"logradouro": vc.Logradouro,
		"bairro":     vc.Bairro,
		"localidade": vc.Localidade,
		"uf":         vc.Uf,
	})
}

func onlyDigits(s string) string {
	re := regexp.MustCompile(`\D+`)
	return re.ReplaceAllString(s, "")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	case "linux":
		exec.Command("xdg-open", url).Start()
	default:
		log.Println("Não sei abrir o navegador automaticamente nesse SO. Acesse:", url)
	}
}
