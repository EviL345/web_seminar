package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

// Структуры данных
type Recipe struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	ChefID      int      `json:"chef_id"`
	ChefName    string   `json:"chef_name"`
	VideoURL    string   `json:"video_url"`
	CreatedAt   string   `json:"created_at"`
}

type Chef struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Speciality  string  `json:"speciality"`
	Rating      float64 `json:"rating"`
	Avatar      string  `json:"avatar"`
	Description string  `json:"description"`
}

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Preferences string `json:"preferences"`
}

type MasterClass struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ChefID      int    `json:"chef_id"`
	ChefName    string `json:"chef_name"`
	DateTime    string `json:"datetime"`
	Duration    int    `json:"duration"`
	Price       int    `json:"price"`
	MaxStudents int    `json:"max_students"`
	Description string `json:"description"`
}

type Subscription struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	ChefID int `json:"chef_id"`
}

type Enrollment struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	MasterClassID int    `json:"master_class_id"`
	EnrolledAt    string `json:"enrolled_at"`
}

// База данных
var db *sql.DB

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func initDB() {
	var err error

	// Удаляем старую базу данных если она существует и повреждена
	dbPath := "./cooking_platform.db"
	if _, err := os.Stat(dbPath); err == nil {
		// Попытаемся открыть существующую базу
		testDB, testErr := sql.Open("sqlite", dbPath)
		if testErr == nil {
			var testCount int
			testQueryErr := testDB.QueryRow("SELECT COUNT(*) FROM sqlite_master").Scan(&testCount)
			testDB.Close()
			if testQueryErr != nil {
				log.Println("Existing database is corrupted, removing...")
				os.Remove(dbPath)
			}
		}
	}

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("Database connected successfully")
	createTables()
	seedData()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS chefs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			speciality TEXT,
			rating REAL DEFAULT 0,
			avatar TEXT,
			description TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			preferences TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			ingredients TEXT,
			chef_id INTEGER,
			video_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (chef_id) REFERENCES chefs (id)
		)`,
		`CREATE TABLE IF NOT EXISTS master_classes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			chef_id INTEGER,
			datetime DATETIME,
			duration INTEGER,
			price INTEGER,
			max_students INTEGER,
			description TEXT,
			FOREIGN KEY (chef_id) REFERENCES chefs (id)
		)`,
		`CREATE TABLE IF NOT EXISTS subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			chef_id INTEGER,
			UNIQUE(user_id, chef_id),
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (chef_id) REFERENCES chefs (id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			master_class_id INTEGER,
			attended_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, master_class_id),
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (master_class_id) REFERENCES master_classes (id)
		)`,
	}

	for i, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating table %d: %v", i+1, err)
		} else {
			log.Printf("Table %d created successfully", i+1)
		}
	}
}

func seedData() {
	// Проверяем, есть ли уже данные
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM chefs").Scan(&count)
	if err != nil {
		log.Printf("Error checking existing data: %v", err)
		return
	}

	if count > 0 {
		log.Println("Data already exists, skipping seed")
		return
	}

	log.Println("Seeding database...")

	// Добавляем поваров
	chefs := []Chef{
		{Name: "Гордон Рамзи", Speciality: "Европейская кухня", Rating: 4.9, Avatar: "https://via.placeholder.com/150", Description: "Мишленовский шеф-повар с мировым именем"},
		{Name: "Юлия Высоцкая", Speciality: "Русская кухня", Rating: 4.7, Avatar: "https://via.placeholder.com/150", Description: "Популярный телеведущий и кулинар"},
		{Name: "Джейми Оливер", Speciality: "Итальянская кухня", Rating: 4.8, Avatar: "https://via.placeholder.com/150", Description: "Британский повар, ресторатор и автор кулинарных книг"},
	}

	for _, chef := range chefs {
		_, err := db.Exec("INSERT INTO chefs (name, speciality, rating, avatar, description) VALUES (?, ?, ?, ?, ?)",
			chef.Name, chef.Speciality, chef.Rating, chef.Avatar, chef.Description)
		if err != nil {
			log.Printf("Error inserting chef %s: %v", chef.Name, err)
		}
	}

	// Добавляем рецепты
	recipes := []Recipe{
		{Title: "Говядина Веллингтон", Description: "Классическое английское блюдо", Ingredients: []string{"говядина", "тесто слоеное", "грибы", "паштет"}, ChefID: 1, VideoURL: "https://www.youtube.com/watch?v=example1"},
		{Title: "Борщ украинский", Description: "Традиционный славянский суп", Ingredients: []string{"свекла", "капуста", "морковь", "лук", "мясо"}, ChefID: 2, VideoURL: "https://www.youtube.com/watch?v=example2"},
		{Title: "Паста Карбонара", Description: "Римская паста с беконом и яйцами", Ingredients: []string{"спагетти", "бекон", "яйца", "пармезан", "черный перец"}, ChefID: 3, VideoURL: "https://www.youtube.com/watch?v=example3"},
		{Title: "Ризотто с грибами", Description: "Кремовое итальянское ризотто", Ingredients: []string{"рис арборио", "грибы", "лук", "вино белое", "пармезан"}, ChefID: 3, VideoURL: "https://www.youtube.com/watch?v=example4"},
	}

	for _, recipe := range recipes {
		ingredientsJSON, _ := json.Marshal(recipe.Ingredients)
		_, err := db.Exec("INSERT INTO recipes (title, description, ingredients, chef_id, video_url) VALUES (?, ?, ?, ?, ?)",
			recipe.Title, recipe.Description, string(ingredientsJSON), recipe.ChefID, recipe.VideoURL)
		if err != nil {
			log.Printf("Error inserting recipe %s: %v", recipe.Title, err)
		}
	}

	// Добавляем мастер-классы
	masterClasses := []MasterClass{
		{Title: "Секреты идеального стейка", ChefID: 1, DateTime: "2024-06-01 18:00", Duration: 120, Price: 5000, MaxStudents: 15, Description: "Научитесь готовить стейк как настоящий профессионал"},
		{Title: "Домашняя выпечка", ChefID: 2, DateTime: "2024-06-02 16:00", Duration: 180, Price: 3500, MaxStudents: 20, Description: "Традиционные рецепты русской выпечки"},
		{Title: "Итальянская паста", ChefID: 3, DateTime: "2024-06-03 19:00", Duration: 90, Price: 4000, MaxStudents: 12, Description: "Готовим пасту с нуля до подачи"},
	}

	for _, mc := range masterClasses {
		_, err := db.Exec("INSERT INTO master_classes (title, chef_id, datetime, duration, price, max_students, description) VALUES (?, ?, ?, ?, ?, ?, ?)",
			mc.Title, mc.ChefID, mc.DateTime, mc.Duration, mc.Price, mc.MaxStudents, mc.Description)
		if err != nil {
			log.Printf("Error inserting master class %s: %v", mc.Title, err)
		}
	}

	// Добавляем пользователей
	users := []User{
		{Username: "foodlover", Email: "food@example.com", Preferences: "итальянская,европейская"},
		{Username: "homecook", Email: "home@example.com", Preferences: "русская,домашняя"},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (username, email, preferences) VALUES (?, ?, ?)",
			user.Username, user.Email, user.Preferences)
		if err != nil {
			log.Printf("Error inserting user %s: %v", user.Username, err)
		}
	}

	log.Println("Database seeded successfully")
}

// API обработчики
func getRecipes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT r.id, r.title, r.description, r.ingredients, r.chef_id, c.name, r.video_url, r.created_at 
		FROM recipes r 
		JOIN chefs c ON r.chef_id = c.id
	`)
	if err != nil {
		log.Printf("Error querying recipes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		var ingredientsJSON string
		err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Description, &ingredientsJSON,
			&recipe.ChefID, &recipe.ChefName, &recipe.VideoURL, &recipe.CreatedAt)
		if err != nil {
			log.Printf("Error scanning recipe: %v", err)
			continue
		}
		json.Unmarshal([]byte(ingredientsJSON), &recipe.Ingredients)
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func createRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		log.Println("1")
		return
	}

	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("2")
		return
	}

	ingredientsJSON, _ := json.Marshal(recipe.Ingredients)
	result, err := db.Exec("INSERT INTO recipes (title, description, ingredients, chef_id, video_url) VALUES (?, ?, ?, ?, ?)",
		recipe.Title, recipe.Description, string(ingredientsJSON), recipe.ChefID, recipe.VideoURL)

	if err != nil {
		log.Printf("Error creating recipe: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("3")
		return
	}

	id, _ := result.LastInsertId()
	recipe.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
	log.Println("4")
}

func getChefs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, speciality, rating, avatar, description FROM chefs")
	if err != nil {
		log.Printf("Error querying chefs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chefs []Chef
	for rows.Next() {
		var chef Chef
		err := rows.Scan(&chef.ID, &chef.Name, &chef.Speciality, &chef.Rating, &chef.Avatar, &chef.Description)
		if err != nil {
			log.Printf("Error scanning chef: %v", err)
			continue
		}
		chefs = append(chefs, chef)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chefs)
}

func getMasterClasses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT mc.id, mc.title, mc.chef_id, c.name, mc.datetime, mc.duration, mc.price, mc.max_students, mc.description 
		FROM master_classes mc 
		JOIN chefs c ON mc.chef_id = c.id
		ORDER BY mc.datetime
	`)
	if err != nil {
		log.Printf("Error querying master classes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var masterClasses []MasterClass
	for rows.Next() {
		var mc MasterClass
		err := rows.Scan(&mc.ID, &mc.Title, &mc.ChefID, &mc.ChefName, &mc.DateTime, &mc.Duration, &mc.Price, &mc.MaxStudents, &mc.Description)
		if err != nil {
			log.Printf("Error scanning master class: %v", err)
			continue
		}
		masterClasses = append(masterClasses, mc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masterClasses)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email, preferences FROM users")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Preferences)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (username, email, preferences) VALUES (?, ?, ?)",
		user.Username, user.Email, user.Preferences)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func generateShoppingList(w http.ResponseWriter, r *http.Request) {
	recipeID := r.URL.Query().Get("recipe_id")
	if recipeID == "" {
		http.Error(w, "recipe_id is required", http.StatusBadRequest)
		return
	}

	var ingredientsJSON string
	err := db.QueryRow("SELECT ingredients FROM recipes WHERE id = ?", recipeID).Scan(&ingredientsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	var ingredients []string
	json.Unmarshal([]byte(ingredientsJSON), &ingredients)

	response := map[string]interface{}{
		"recipe_id":     recipeID,
		"shopping_list": ingredients,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Получаем предпочтения пользователя
	var preferences string
	db.QueryRow("SELECT preferences FROM users WHERE id = ?", userID).Scan(&preferences)

	// Получаем поваров, на которых подписан пользователь
	subscribedChefs := []int{}
	rows, _ := db.Query("SELECT chef_id FROM subscriptions WHERE user_id = ?", userID)
	for rows.Next() {
		var chefID int
		rows.Scan(&chefID)
		subscribedChefs = append(subscribedChefs, chefID)
	}
	rows.Close()

	// Формируем рекомендации на основе предпочтений и подписок
	query := `
		SELECT mc.id, mc.title, mc.chef_id, c.name, mc.datetime, mc.duration, mc.price, mc.max_students, mc.description 
		FROM master_classes mc 
		JOIN chefs c ON mc.chef_id = c.id
		WHERE mc.datetime > datetime('now')
	`

	args := []interface{}{}
	if len(subscribedChefs) > 0 {
		placeholders := strings.Repeat("?,", len(subscribedChefs)-1) + "?"
		query += " AND (mc.chef_id IN (" + placeholders + ")"
		for _, chefID := range subscribedChefs {
			args = append(args, chefID)
		}

		if preferences != "" {
			query += " OR c.speciality LIKE ?)"
			args = append(args, "%"+preferences+"%")
		} else {
			query += ")"
		}
	} else if preferences != "" {
		query += " AND c.speciality LIKE ?"
		args = append(args, "%"+preferences+"%")
	}

	query += " ORDER BY mc.datetime LIMIT 10"

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error querying recommendations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recommendations []MasterClass
	for rows.Next() {
		var mc MasterClass
		err := rows.Scan(&mc.ID, &mc.Title, &mc.ChefID, &mc.ChefName, &mc.DateTime, &mc.Duration, &mc.Price, &mc.MaxStudents, &mc.Description)
		if err != nil {
			continue
		}
		recommendations = append(recommendations, mc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var sub Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT OR IGNORE INTO subscriptions (user_id, chef_id) VALUES (?, ?)", sub.UserID, sub.ChefID)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "subscribed"})
}

func enrollInMasterClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли места на мастер-классе
	var currentEnrollments, maxStudents int
	err := db.QueryRow(`
		SELECT COUNT(uh.id), mc.max_students 
		FROM master_classes mc 
		LEFT JOIN user_history uh ON mc.id = uh.master_class_id 
		WHERE mc.id = ?
		GROUP BY mc.id, mc.max_students
	`, enrollment.MasterClassID).Scan(&currentEnrollments, &maxStudents)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking enrollment capacity: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if currentEnrollments >= maxStudents {
		http.Error(w, "No available spots", http.StatusConflict)
		return
	}

	// Записываем пользователя
	_, err = db.Exec("INSERT OR IGNORE INTO user_history (user_id, master_class_id) VALUES (?, ?)",
		enrollment.UserID, enrollment.MasterClassID)

	if err != nil {
		log.Printf("Error enrolling user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "enrolled"})
}

func getUserHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT uh.id, uh.user_id, uh.master_class_id, mc.title, c.name, uh.attended_at
		FROM user_history uh
		JOIN master_classes mc ON uh.master_class_id = mc.id
		JOIN chefs c ON mc.chef_id = c.id
		WHERE uh.user_id = ?
		ORDER BY uh.attended_at DESC
	`, userID)

	if err != nil {
		log.Printf("Error querying user history: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type HistoryEntry struct {
		ID            int    `json:"id"`
		UserID        int    `json:"user_id"`
		MasterClassID int    `json:"master_class_id"`
		ClassTitle    string `json:"class_title"`
		ChefName      string `json:"chef_name"`
		AttendedAt    string `json:"attended_at"`
	}

	var history []HistoryEntry
	for rows.Next() {
		var entry HistoryEntry
		err := rows.Scan(&entry.ID, &entry.UserID, &entry.MasterClassID,
			&entry.ClassTitle, &entry.ChefName, &entry.AttendedAt)
		if err != nil {
			continue
		}
		history = append(history, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func getUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT s.id, s.user_id, s.chef_id, c.name, c.speciality, c.rating
		FROM subscriptions s
		JOIN chefs c ON s.chef_id = c.id
		WHERE s.user_id = ?
	`, userID)

	if err != nil {
		log.Printf("Error querying user subscriptions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserSubscription struct {
		ID         int     `json:"id"`
		UserID     int     `json:"user_id"`
		ChefID     int     `json:"chef_id"`
		ChefName   string  `json:"chef_name"`
		Speciality string  `json:"speciality"`
		ChefRating float64 `json:"chef_rating"`
	}

	var subscriptions []UserSubscription
	for rows.Next() {
		var sub UserSubscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.ChefID,
			&sub.ChefName, &sub.Speciality, &sub.ChefRating)
		if err != nil {
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	type Stats struct {
		TotalRecipes       int `json:"total_recipes"`
		TotalChefs         int `json:"total_chefs"`
		TotalUsers         int `json:"total_users"`
		TotalMasterClasses int `json:"total_master_classes"`
		TotalEnrollments   int `json:"total_enrollments"`
	}

	var stats Stats

	db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&stats.TotalRecipes)
	db.QueryRow("SELECT COUNT(*) FROM chefs").Scan(&stats.TotalChefs)
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	db.QueryRow("SELECT COUNT(*) FROM master_classes").Scan(&stats.TotalMasterClasses)
	db.QueryRow("SELECT COUNT(*) FROM user_history").Scan(&stats.TotalEnrollments)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func searchRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "search query is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT r.id, r.title, r.description, r.ingredients, r.chef_id, c.name, r.video_url, r.created_at 
		FROM recipes r 
		JOIN chefs c ON r.chef_id = c.id
		WHERE r.title LIKE ? OR r.description LIKE ? OR r.ingredients LIKE ?
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		var ingredientsJSON string
		err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Description, &ingredientsJSON,
			&recipe.ChefID, &recipe.ChefName, &recipe.VideoURL, &recipe.CreatedAt)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(ingredientsJSON), &recipe.Ingredients)
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	http.ServeFile(w, r, "index.html")
}
func main() {
	initDB()
	defer db.Close()

	// API маршруты
	http.HandleFunc("/api/recipes", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getRecipes(w, r)
		} else if r.Method == "POST" {
			createRecipe(w, r)
		}
	}))

	http.HandleFunc("/api/chefs", corsMiddleware(getChefs))
	http.HandleFunc("/api/masterclasses", corsMiddleware(getMasterClasses))
	http.HandleFunc("/api/users", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getUsers(w, r)
		} else if r.Method == "POST" {
			createUser(w, r)
		}
	}))

	http.HandleFunc("/api/shopping-list", corsMiddleware(generateShoppingList))
	http.HandleFunc("/api/recommendations", corsMiddleware(getRecommendations))
	http.HandleFunc("/api/subscribe", corsMiddleware(subscribe))
	http.HandleFunc("/api/enroll", corsMiddleware(enrollInMasterClass))
	http.HandleFunc("/api/user-history", corsMiddleware(getUserHistory))
	http.HandleFunc("/api/user-subscriptions", corsMiddleware(getUserSubscriptions))
	http.HandleFunc("/api/stats", corsMiddleware(getStats))
	http.HandleFunc("/api/search", corsMiddleware(searchRecipes))

	// Главная страница
	http.HandleFunc("/", corsMiddleware(homePage))

	// Статические файлы (если нужно)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fmt.Println("🍳 Кулинарная платформа запущена на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
