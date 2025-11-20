package services

import "picturebook/server/internal/models"

func sampleBooks() []models.Book {
    pagesA := []models.Page{}
    for i := 0; i < 8; i++ {
        pagesA = append(pagesA, models.Page{ID: "A-" + itoa(i), Index: i, ImageURL: "", Duration: nil})
    }
    pagesB := []models.Page{}
    for i := 0; i < 6; i++ {
        pagesB = append(pagesB, models.Page{ID: "B-" + itoa(i), Index: i, ImageURL: "", Duration: nil})
    }
    pagesC := []models.Page{}
    for i := 0; i < 10; i++ {
        pagesC = append(pagesC, models.Page{ID: "C-" + itoa(i), Index: i, ImageURL: "", Duration: nil})
    }
    return []models.Book{
        {ID: "book_a", Title: "森林探险", CoverURL: "cover_a", AgeMin: 3, AgeMax: 6, Tags: []string{"森林","动物"}, PopularityScore: 0.82, ThemeKeywords: []string{"自然","探险"}, IsEditorPick: true, Pages: pagesA, Status: "published", CategoryID: "default"},
        {ID: "book_b", Title: "海洋之旅", CoverURL: "cover_b", AgeMin: 5, AgeMax: 9, Tags: []string{"海洋","旅行"}, PopularityScore: 0.74, ThemeKeywords: []string{"自然","探索"}, IsEditorPick: true, Pages: pagesB, Status: "published", CategoryID: "default"},
        {ID: "book_c", Title: "太空奇遇记", CoverURL: "cover_c", AgeMin: 6, AgeMax: 10, Tags: []string{"太空","科学"}, PopularityScore: 0.91, ThemeKeywords: []string{"宇宙","冒险"}, IsEditorPick: false, Pages: pagesC, Status: "published", CategoryID: "default"},
        {ID: "book_d", Title: "农场的一天", CoverURL: "cover_d", AgeMin: 3, AgeMax: 7, Tags: []string{"农场","动物"}, PopularityScore: 0.65, ThemeKeywords: []string{"生活","自然"}, IsEditorPick: false, Pages: pagesB, Status: "published", CategoryID: "default"},
        {ID: "book_e", Title: "恐龙乐园", CoverURL: "cover_e", AgeMin: 4, AgeMax: 8, Tags: []string{"恐龙","历史"}, PopularityScore: 0.88, ThemeKeywords: []string{"史前","冒险"}, IsEditorPick: true, Pages: pagesA, Status: "published", CategoryID: "default"},
        {ID: "book_f", Title: "城市探秘", CoverURL: "cover_f", AgeMin: 5, AgeMax: 10, Tags: []string{"城市","建筑"}, PopularityScore: 0.59, ThemeKeywords: []string{"生活","探索"}, IsEditorPick: false, Pages: pagesA, Status: "published", CategoryID: "default"},
    }
}

func itoa(i int) string { return fmtInt(i) }

func fmtInt(i int) string {
    if i == 0 { return "0" }
    s := ""
    n := i
    for n > 0 {
        d := n % 10
        s = string('0'+d) + s
        n /= 10
    }
    return s
}