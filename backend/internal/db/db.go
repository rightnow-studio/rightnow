package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Init(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate failed: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			text_key TEXT NOT NULL DEFAULT '',
			duration_seconds INTEGER NOT NULL DEFAULT 60,
			category TEXT NOT NULL DEFAULT 'work',
			is_preset BOOLEAN NOT NULL DEFAULT FALSE,
			created_by TEXT NOT NULL DEFAULT '',
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_task_pool (
			user_id TEXT NOT NULL DEFAULT '',
			task_id TEXT NOT NULL DEFAULT '',
			added_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, task_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (task_id) REFERENCES tasks(id)
		);`,
		`CREATE TABLE IF NOT EXISTS task_history (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL DEFAULT '',
			task_id TEXT NOT NULL DEFAULT '',
			started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			ended_at DATETIME NOT NULL DEFAULT '1970-01-01T00:00:00Z',
			outcome TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (task_id) REFERENCES tasks(id)
		);`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return seedPresetTasks(db)
}

func seedPresetTasks(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE is_preset = TRUE").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	presetTasks := []struct {
		id, textKey, category string
		duration              int
	}{
		{"preset_001", "task.open_first_desktop_file", "work", 60},
		{"preset_002", "task.type_period_in_new_doc", "work", 30},
		{"preset_003", "task.close_all_notifications", "work", 60},
		{"preset_004", "task.write_todays_date", "work", 30},
		{"preset_005", "task.mark_first_email_read", "work", 30},
		{"preset_006", "task.move_three_files_to_folder", "work", 60},
		{"preset_007", "task.adjust_screen_brightness", "work", 10},
		{"preset_008", "task.close_first_browser_tab", "work", 10},
		{"preset_009", "task.put_nearest_clothes_in_washer", "home", 30},
		{"preset_010", "task.take_kitchen_trash_out", "home", 60},
		{"preset_011", "task.smooth_blanket_and_pillow", "home", 120},
		{"preset_012", "task.pick_up_socks_to_basket", "home", 30},
		{"preset_013", "task.lay_out_tomorrows_clothes", "home", 120},
		{"preset_014", "task.move_nearest_bowl_to_sink", "home", 30},
		{"preset_015", "task.fold_snack_bag_into_drawer", "home", 60},
		{"preset_016", "task.pick_up_three_things_from_floor", "home", 120},
		{"preset_017", "task.send_period_to_first_contact", "communication", 30},
		{"preset_018", "task.send_emoji_to_top_chat", "communication", 30},
		{"preset_019", "task.find_first_A_contact", "communication", 60},
		{"preset_020", "task.clear_all_red_badges", "communication", 60},
		{"preset_021", "task.put_on_nearest_shoes_and_tiptoe", "health", 60},
		{"preset_022", "task.stretch_arms_overhead_five_times", "health", 60},
		{"preset_023", "task.walk_to_nearest_window_look_out", "health", 60},
		{"preset_024", "task.pour_water_and_drink_three_sips", "health", 60},
		{"preset_025", "task.flip_phone_screen_down", "health", 10},
		{"preset_026", "task.close_eyes_count_breaths_to_five", "health", 60},
		{"preset_027", "task.write_start_on_paper", "admin", 30},
		{"preset_028", "task.stack_top_three_papers", "admin", 120},
		{"preset_029", "task.take_out_and_put_back_cards", "admin", 60},
		{"preset_030", "task.empty_bag_and_put_back", "admin", 120},
	}
	for _, t := range presetTasks {
		_, err := db.Exec(
			"INSERT INTO tasks (id, text_key, duration_seconds, category, is_preset) VALUES (?, ?, ?, ?, TRUE)",
			t.id, t.textKey, t.duration, t.category,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
