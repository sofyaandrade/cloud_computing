package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Exercise struct {
	ID          int
	Name        string
	Category    string
	Description string
	CreatedAt   time.Time
}

var db *sql.DB
var tmpl = template.Must(template.New("index").Parse(pageHTML))

func main() {
	var err error

	db, err = connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/exercicios", exercisesHandler)
	http.HandleFunc("/editar-exercicios", editarExercicio)
	http.HandleFunc("/delete-exercicios", deletarExercicio)

	port := getenv("APP_PORT", "8080")
	log.Printf("Aplicação iniciada na porta %s", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectDB() (*sql.DB, error) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASSWORD", "postgres")
	name := getenv("DB_NAME", "treinosdb")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		name,
	)

	var database *sql.DB
	var err error

	for i := 1; i <= 20; i++ {
		database, err = sql.Open("postgres", dsn)

		if err == nil {
			err = database.Ping()
		}

		if err == nil {
			log.Println("Conexão com PostgreSQL realizada com sucesso")
			return database, nil
		}

		log.Printf("Aguardando banco de dados... tentativa %d/20", i)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("não foi possível conectar ao banco: %w", err)
}

func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS exercises (
		id SERIAL PRIMARY KEY,
		name VARCHAR(120) NOT NULL,
		category VARCHAR(80) NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/exercicios", http.StatusSeeOther)
}

func exercisesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		category := r.FormValue("category")
		description := r.FormValue("description")

		if name == "" || category == "" || description == "" {
			http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(
			"INSERT INTO exercises (name, category, description) VALUES ($1, $2, $3)",
			name,
			category,
			description,
		)

		if err != nil {
			log.Println("Erro ao cadastrar exercício:", err)
			http.Error(w, "Erro ao cadastrar exercício", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/exercicios", http.StatusSeeOther)
		return
	}

	rows, err := db.Query("SELECT id, name, category, description, created_at FROM exercises ORDER BY id DESC")
	if err != nil {
		log.Println("Erro ao consultar exercícios:", err)
		http.Error(w, "Erro ao consultar exercícios", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	exercises := []Exercise{}

	for rows.Next() {
		var ex Exercise

		err := rows.Scan(
			&ex.ID,
			&ex.Name,
			&ex.Category,
			&ex.Description,
			&ex.CreatedAt,
		)

		if err != nil {
			log.Println("Erro ao ler dados:", err)
			http.Error(w, "Erro ao ler dados", http.StatusInternalServerError)
			return
		}

		exercises = append(exercises, ex)
	}

	if err := rows.Err(); err != nil {
		log.Println("Erro nas linhas do resultado:", err)
		http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
		return
	}

	data := struct {
		Exercises []Exercise
	}{
		Exercises: exercises,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Erro ao renderizar página:", err)
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
	}
}

func editarExercicio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/exercicios", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")
	category := r.FormValue("category")
	description := r.FormValue("description")

	if id == "" || name == "" || category == "" || description == "" {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"UPDATE exercises SET name = $1, category = $2, description = $3 WHERE id = $4",
		name,
		category,
		description,
		id,
	)

	if err != nil {
		log.Println("Erro ao editar exercício:", err)
		http.Error(w, "Erro ao editar exercício", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/exercicios", http.StatusSeeOther)
}

func deletarExercicio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/exercicios", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")

	if id == "" {
		http.Error(w, "ID não informado", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM exercises WHERE id = $1", id)
	if err != nil {
		log.Println("Erro ao deletar exercício:", err)
		http.Error(w, "Erro ao deletar exercício", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/exercicios", http.StatusSeeOther)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

const pageHTML = `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>Plataforma de Treinos Online</title>

	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			background: #f4f6f8;
			color: #222;
		}

		header {
			background: #111827;
			color: white;
			padding: 24px;
			text-align: center;
		}

		main {
			max-width: 1100px;
			margin: 24px auto;
			padding: 0 16px;
		}

		.card {
			background: white;
			border-radius: 10px;
			padding: 20px;
			margin-bottom: 20px;
			box-shadow: 0 2px 8px rgba(0,0,0,0.08);
		}

		label {
			display: block;
			margin-top: 12px;
			font-weight: bold;
		}

		input,
		textarea {
			width: 100%;
			padding: 10px;
			margin-top: 6px;
			border: 1px solid #ccc;
			border-radius: 6px;
			box-sizing: border-box;
		}

		textarea {
			resize: vertical;
		}

		button {
			border: 0;
			padding: 10px 14px;
			border-radius: 6px;
			cursor: pointer;
			font-weight: bold;
		}

		.btn-cadastrar {
			margin-top: 16px;
			background: #2563eb;
			color: white;
		}

		.btn-editar {
			background: #16a34a;
			color: white;
			margin-bottom: 6px;
			width: 100%;
		}

		.btn-deletar {
			background: #dc2626;
			color: white;
			width: 100%;
		}

		table {
			width: 100%;
			border-collapse: collapse;
			margin-top: 12px;
		}

		th,
		td {
			border-bottom: 1px solid #ddd;
			padding: 10px;
			text-align: left;
			vertical-align: top;
		}

		th {
			background: #e5e7eb;
		}

		td input,
		td textarea {
			margin-top: 0;
		}

		.acoes {
			width: 120px;
		}

		.empty {
			color: #666;
		}
	</style>
</head>

<body>
	<header>
		<h1>Plataforma de Treinos Online</h1>
		<p>Cadastro, consulta, edição e exclusão de exercícios usando Go, PostgreSQL, Docker e Docker Compose</p>
	</header>

	<main>
		<section class="card">
			<h2>Cadastrar exercício</h2>

			<form method="POST" action="/exercicios">
				<label for="name">Nome do exercício</label>
				<input id="name" name="name" placeholder="Ex: Agachamento" required>

				<label for="category">Categoria</label>
				<input id="category" name="category" placeholder="Ex: Pernas, Peito, Costas" required>

				<label for="description">Descrição</label>
				<textarea id="description" name="description" rows="4" placeholder="Explique como executar o exercício" required></textarea>

				<button class="btn-cadastrar" type="submit">Cadastrar</button>
			</form>
		</section>

		<section class="card">
			<h2>Exercícios cadastrados</h2>

			{{if .Exercises}}
			<table>
				<thead>
					<tr>
						<th>ID</th>
						<th>Nome</th>
						<th>Categoria</th>
						<th>Descrição</th>
						<th class="acoes">Ações</th>
					</tr>
				</thead>

				<tbody>
					{{range .Exercises}}
					<tr>
						<td>{{.ID}}</td>

						<td>
							<input type="text" name="name" value="{{.Name}}" required form="editar-{{.ID}}">
						</td>

						<td>
							<input type="text" name="category" value="{{.Category}}" required form="editar-{{.ID}}">
						</td>

						<td>
							<textarea name="description" rows="2" required form="editar-{{.ID}}">{{.Description}}</textarea>
						</td>

						<td class="acoes">
							<form id="editar-{{.ID}}" method="POST" action="/editar-exercicios">
								<input type="hidden" name="id" value="{{.ID}}">
								<button class="btn-editar" type="submit">Editar</button>
							</form>

							<form method="POST" action="/delete-exercicios">
								<input type="hidden" name="id" value="{{.ID}}">
								<button class="btn-deletar" type="submit" onclick="return confirm('Tem certeza que deseja deletar este exercício?')">
									Deletar
								</button>
							</form>
						</td>
					</tr>
					{{end}}
				</tbody>
			</table>
			{{else}}
			<p class="empty">Nenhum exercício cadastrado ainda.</p>
			{{end}}
		</section>
	</main>
</body>
</html>
`
