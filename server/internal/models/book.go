package models

type Page struct {
    ID        string  `json:"id"`
    Index     int     `json:"index"`
    ImageURL  string  `json:"imageURL"`
    Duration  *int    `json:"duration_hint,omitempty"`
}

type Book struct {
    ID              string   `json:"id"`
    Title           string   `json:"title"`
    CoverURL        string   `json:"coverURL"`
    AgeMin          int      `json:"ageMin"`
    AgeMax          int      `json:"ageMax"`
    Tags            []string `json:"tags"`
    PopularityScore float64  `json:"popularityScore"`
    ThemeKeywords   []string `json:"themeKeywords"`
    IsEditorPick    bool     `json:"isEditorPick"`
    Pages           []Page   `json:"pages"`
    Status          string   `json:"status"`
}