package handlers

import (
    "os"
    "strconv"
)

func atoi(s string) int { i, _ := strconv.Atoi(s); return i }
func atoiDefault(s string, d int) int { if s == "" { return d }; i, err := strconv.Atoi(s); if err != nil { return d }; return i }
func atofDefault(s string, d float64) float64 { f, err := strconv.ParseFloat(s, 64); if err != nil { return d }; return f }
func splitCSV(s string) []string {
    out := []string{}
    cur := ""
    for i := 0; i < len(s); i++ {
        ch := s[i]
        if ch == ',' {
            if cur != "" { out = append(out, cur) }
            cur = ""
        } else {
            cur += string(ch)
        }
    }
    if cur != "" { out = append(out, cur) }
    return out
}

func uploadsDir() string {
    if v := os.Getenv("UPLOADS_DIR"); v != "" { return v }
    return "/app/uploads"
}