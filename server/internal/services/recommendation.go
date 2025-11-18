package services

import "picturebook/server/internal/models"

func ageFit(b models.Book, age int) float64 {
    if age < b.AgeMin || age > b.AgeMax {
        return 0
    }
    span := b.AgeMax - b.AgeMin
    if span <= 0 {
        return 1
    }
    pos := float64(age-b.AgeMin) / float64(span)
    v := 1 - abs(pos-0.5)*2
    if v < 0 {
        return 0
    }
    return v
}

func themeSimilarity(a, b []string) float64 {
    set := map[string]struct{}{}
    for _, x := range a {
        set[x] = struct{}{}
    }
    inter := 0
    union := len(a)
    seen := map[string]struct{}{}
    for _, y := range b {
        if _, ok := seen[y]; ok {
            continue
        }
        seen[y] = struct{}{}
        union++
        if _, ok := set[y]; ok {
            inter++
        }
    }
    if union == 0 {
        return 0
    }
    return float64(inter) / float64(union)
}

func score(b models.Book, age int, theme []string) float64 {
    p := b.PopularityScore
    a := ageFit(b, age)
    t := themeSimilarity(theme, append(b.ThemeKeywords, b.Tags...))
    return 0.5*p + 0.3*a + 0.2*t
}

func Recommend(current models.Book, all []models.Book, age, limit int) []models.Book {
    base := append(current.ThemeKeywords, current.Tags...)
    filtered := make([]models.Book, 0, len(all))
    for _, b := range all {
        if b.ID != current.ID && b.Status == "published" {
            filtered = append(filtered, b)
        }
    }
    ranked := make([]models.Book, len(filtered))
    copy(ranked, filtered)
    sortBooks(ranked, age, "combinedWithTheme", base)
    if len(ranked) >= limit {
        return ranked[:limit]
    }
    return ranked
}

func abs(v float64) float64 { if v < 0 { return -v } ; return v }