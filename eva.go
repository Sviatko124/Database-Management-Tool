package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"math/rand"
	"time"
	"golang.org/x/term"
	"os/signal"
    "syscall"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	ID          int
	Title       string
	Keywords    string
	AttackStep  string
	Explanation string
	Commands    string
	Notes       string
	CreatedAt   string
	UpdatedAt   string
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return home
}

func initDB() *sql.DB {
	dataDir := filepath.Join(getHomeDir(), ".eva")
	os.MkdirAll(dataDir, 0755)
	dbPath := filepath.Join(dataDir, "eva.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			keywords TEXT,
			attack_step TEXT,
			explanation TEXT,
			commands TEXT,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		os.Exit(1)
	}

	return db
}
func prompt(message string) string {
    fmt.Print("\033[36m" + message + "\033[0m")
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        panic(err)
    }
    
    cleanup := func() {
        term.Restore(int(os.Stdin.Fd()), oldState)
        fmt.Println("\n\033[32mGoodbye!\033[0m")
        os.Exit(0)
    }
    
    go func() {
        <-sigChan
        cleanup()
    }()
    
    terminal := term.NewTerminal(os.NewFile(os.Stdin.Fd(), ""), "")
    
    line, err := terminal.ReadLine()
    if err != nil {
        cleanup()
        return ""
    }

    signal.Stop(sigChan)
    
    term.Restore(int(os.Stdin.Fd()), oldState)
    
    return line
}

func confirm(message string) bool {
	response := prompt(message + " (y/n): ")
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func addEntry(db *sql.DB) {
	fmt.Println("\n\033[34mAdding New Entry\033[0m")
	
	title := prompt("Title: ")
	keywords := prompt("Keywords (comma-separated): ")
	attackStep := prompt("Attack Step: ")
	explanation := prompt("Explanation: ")
	commands := prompt("Commands (use \\n for new lines): ")
	notes := prompt("Additional Notes: ")

	if strings.TrimSpace(notes) == "" {
		notes = "No additional notes."
	}

	commands = strings.ReplaceAll(commands, "\\n", "\n")
	explanation = strings.ReplaceAll(explanation, "\\n", "\n")
	notes = strings.ReplaceAll(notes, "\\n", "\n")

	_, err := db.Exec(`
		INSERT INTO entries (title, keywords, attack_step, explanation, commands, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`, title, keywords, attackStep, explanation, commands, notes)

	if err != nil {
		fmt.Println("\033[31mError adding entry:", err, "\033[0m")
		return
	}

	fmt.Println("\033[32mEntry added successfully\033[0m")
}

func showEntryDetails(db *sql.DB, id string, showModifyPrompt bool) {
	var entry Entry
	err := db.QueryRow(`
		SELECT id, title, keywords, attack_step, explanation, commands, notes, created_at, updated_at 
		FROM entries WHERE id = ?
	`, id).Scan(&entry.ID, &entry.Title, &entry.Keywords, &entry.AttackStep, 
		&entry.Explanation, &entry.Commands, &entry.Notes, &entry.CreatedAt, &entry.UpdatedAt)

	if err != nil {
		fmt.Println("\033[31mEntry not found\033[0m")
		return
	}

	fmt.Println("\n\033[34mEntry Details\033[0m")
	fmt.Printf("\033[36mID:\033[0m %d\n", entry.ID)
	fmt.Printf("\033[36mTitle:\033[0m %s\n", entry.Title)
	fmt.Printf("\033[36mKeywords:\033[0m %s\n", entry.Keywords)
	fmt.Printf("\033[36mAttack Step:\033[0m %s\n", entry.AttackStep)
	fmt.Printf("\033[36mExplanation:\033[0m %s\n", entry.Explanation)
	fmt.Printf("\033[36mCommands:\033[0m\n%s\n", entry.Commands)
	fmt.Printf("\033[36mNotes:\033[0m %s\n", entry.Notes)
	fmt.Printf("\033[36mCreated:\033[0m %s\n", entry.CreatedAt)
	fmt.Printf("\033[36mLast Updated:\033[0m %s\n", entry.UpdatedAt)

	if showModifyPrompt && confirm("\nWould you like to modify this entry?") {
		modifyEntry(db, id)
	}
}

func modifyEntry(db *sql.DB, id string) {
    fmt.Println("\n\033[34mModify Entry\033[0m")
    fmt.Println("Which section would you like to modify?")
    fmt.Println("1. Title")
    fmt.Println("2. Keywords")
    fmt.Println("3. Attack Step")
    fmt.Println("4. Explanation")
    fmt.Println("5. Commands")
    fmt.Println("6. Notes")
    fmt.Println("7. Delete Entry")
    fmt.Println("8. Cancel")

    choice := prompt("Select an option (1-8): ")

    if choice == "8" {
        fmt.Println("\033[33mModification cancelled.\033[0m")
        return
    }

    if choice == "7" {
        deleteEntry(db, id)
        return
    }

	var field, currentContent string
	var row *sql.Row

	switch choice {
	case "1":
		field = "title"
		row = db.QueryRow("SELECT title FROM entries WHERE id = ?", id)
	case "2":
		field = "keywords"
		row = db.QueryRow("SELECT keywords FROM entries WHERE id = ?", id)
	case "3":
		field = "attack_step"
		row = db.QueryRow("SELECT attack_step FROM entries WHERE id = ?", id)
	case "4":
		field = "explanation"
		row = db.QueryRow("SELECT explanation FROM entries WHERE id = ?", id)
	case "5":
		field = "commands"
		row = db.QueryRow("SELECT commands FROM entries WHERE id = ?", id)
	case "6":
		field = "notes"
		row = db.QueryRow("SELECT notes FROM entries WHERE id = ?", id)
	default:
		fmt.Println("\033[33mInvalid choice. Modification cancelled.\033[0m")
		return
	}

	err := row.Scan(&currentContent)
	if err != nil {
		fmt.Println("\033[31mError retrieving current content:", err, "\033[0m")
		return
	}

	fmt.Println("\n\033[36mCurrent content:\033[0m")
	fmt.Println(currentContent)

	newContent := prompt("\nEnter new content: ")

	if field == "title" && strings.TrimSpace(newContent) == "" {
		fmt.Println("\033[33mTitle cannot be empty. Modification cancelled.\033[0m")
		return
	}

	if field == "notes" && strings.TrimSpace(newContent) == "" {
		newContent = "No additional notes."
	}

	if field == "explanation" || field == "commands" || field == "notes" {
		newContent = strings.ReplaceAll(newContent, "\\n", "\n")
	}

	if !confirm("\nDo you want to save these changes?") {
		fmt.Println("\033[33mModification cancelled.\033[0m")
		return
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE entries SET %s = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", field), newContent, id)
	if err != nil {
		fmt.Println("\033[31mError updating entry:", err, "\033[0m")
		return
	}

	fmt.Println("\033[32mEntry updated successfully\033[0m")
	showEntryDetails(db, id, false)
}

func deleteEntry(db *sql.DB, id string) {
    if !confirm("\n\033[31mWarning: This will permanently delete the entry. Are you sure?\033[0m") {
        fmt.Println("\033[33mDeletion cancelled.\033[0m")
        return
    }

    verificationCode := rand.Intn(9000) + 1000
    fmt.Printf("\n\033[31mTo confirm deletion, please enter this 4-digit code: %d\033[0m\n", verificationCode)
    
    userInput := prompt("Enter verification code: ")
    
    if userInput != fmt.Sprintf("%d", verificationCode) {
        fmt.Println("\033[31mIncorrect verification code. Deletion cancelled.\033[0m")
        return
    }

    tx, err := db.Begin()
    if err != nil {
        fmt.Println("\033[31mError starting transaction:", err, "\033[0m")
        return
    }

    result, err := tx.Exec("DELETE FROM entries WHERE id = ?", id)
    if err != nil {
        tx.Rollback()
        fmt.Println("\033[31mError deleting entry:", err, "\033[0m")
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        fmt.Println("\033[31mError checking deletion status:", err, "\033[0m")
        return
    }

    if rowsAffected == 0 {
        tx.Rollback()
        fmt.Println("\033[31mNo entry found with that ID.\033[0m")
        return
    }

    _, err = tx.Exec(`
        UPDATE entries 
        SET id = id - 1 
        WHERE id > ?
    `, id)
    
    if err != nil {
        tx.Rollback()
        fmt.Println("\033[31mError reindexing entries:", err, "\033[0m")
        return
    }

    _, err = tx.Exec(`
        UPDATE sqlite_sequence 
        SET seq = (SELECT MAX(id) FROM entries) 
        WHERE name = 'entries'
    `)

    if err != nil {
        tx.Rollback()
        fmt.Println("\033[31mError resetting sequence:", err, "\033[0m")
        return
    }

    err = tx.Commit()
    if err != nil {
        fmt.Println("\033[31mError committing changes:", err, "\033[0m")
        return
    }

    fmt.Println("\033[32mEntry successfully deleted and database reindexed.\033[0m")
}

func searchEntries(db *sql.DB) {
    fmt.Println("\n\033[34mSearch Entries\033[0m")
    searchInput := prompt("Enter search terms (separate multiple terms with spaces): ")
    
    terms := strings.Fields(searchInput)
    
    query := `
        SELECT id, title, keywords, attack_step, explanation, commands, notes
        FROM entries
        WHERE 1=1
    `
    params := []interface{}{}
    
    for _, term := range terms {
        term = "%" + term + "%"
        query += `
            AND (
                title LIKE ? 
                OR keywords LIKE ?
                OR attack_step LIKE ?
                OR explanation LIKE ?
                OR commands LIKE ?
                OR notes LIKE ?
            )`
        params = append(params, term, term, term, term, term, term)
    }

    rows, err := db.Query(query, params...)
    if err != nil {
        fmt.Println("\033[31mError searching entries:", err, "\033[0m")
        return
    }
    defer rows.Close()

    fmt.Println("\n\033[36mID\tTitle\t\t\t\tKeywords\t\t\t\tAttack Step\033[0m")
    fmt.Println("----------------------------------------------------------------------------------------")

    found := false
    for rows.Next() {
        found = true
        var entry Entry
        err := rows.Scan(&entry.ID, &entry.Title, &entry.Keywords, &entry.AttackStep,
            &entry.Explanation, &entry.Commands, &entry.Notes)
        if err != nil {
            fmt.Println("\033[31mError reading entry:", err, "\033[0m")
            continue
        }
        fmt.Printf("%d\t%-30s\t%-30s\t%-30s\n", 
            entry.ID, truncate(entry.Title, 30), truncate(entry.Keywords, 30), truncate(entry.AttackStep, 30))
    }

    if !found {
        fmt.Println("\033[33mNo entries found.\033[0m")
        return
    }

    if confirm("\nWould you like to see full details of any entry?") {
        id := prompt("Enter entry ID: ")
        showEntryDetails(db, id, true)
    }
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func main() {
    rand.Seed(time.Now().UnixNano())

    db := initDB()
    defer db.Close()

	for {
		fmt.Println("\n\033[34mEVA - Cheatsheet/Notes Database Tool\033[0m")
		fmt.Println("1. Search database")
		fmt.Println("2. Add entry")
		fmt.Println("3. Exit")

		choice := prompt("Select an option (1-3): ")

		switch choice {
		case "1":
			searchEntries(db)
		case "2":
			addEntry(db)
		case "3":
			fmt.Println("\033[32mGoodbye!\033[0m")
			os.Exit(0)
		}
	}
}
