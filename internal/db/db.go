package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(storagePath string) error {
	dbPath := filepath.Join(storagePath, "stats.db")

	// 确保目录存在
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}

	// 性能优化：启用 WAL 模式
	// WAL 模式允许并发读写，显著提高性能
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
	}

	for _, pragma := range pragmas {
		if _, err := DB.Exec(pragma); err != nil {
			return fmt.Errorf("执行 PRAGMA 失败 (%s): %w", pragma, err)
		}
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	return createTables()
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS visits (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ip TEXT,
            path TEXT,
            user_agent TEXT,
            referer TEXT,
            country TEXT,
            region TEXT,
            city TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS downloads (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            file_name TEXT,
            launcher TEXT,
            version TEXT,
            ip TEXT,
            country TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE INDEX IF NOT EXISTS idx_visits_created_at ON visits(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_downloads_created_at ON downloads(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_downloads_file_name ON downloads(file_name)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("创建表失败: %w, query: %s", err, query)
		}
	}
	return nil
}
