package services

import (
    "sort"
    "picturebook/server/internal/models"
)

func sortBooks(items []models.Book, age int, sortBy string, themeOpt ...[]string) {
    sort.Slice(items, func(i, j int) bool {
        a := items[i]
        b := items[j]
        switch sortBy {
        case "popular":
            if a.PopularityScore == b.PopularityScore {
                return a.AgeMin < b.AgeMin
            }
            return a.PopularityScore > b.PopularityScore
        case "age_fit":
            af := ageFit(a, age)
            bf := ageFit(b, age)
            if af == bf {
                return a.PopularityScore > b.PopularityScore
            }
            return af > bf
        case "combinedWithTheme":
            theme := []string{}
            if len(themeOpt) > 0 {
                theme = themeOpt[0]
            }
            sa := score(a, age, theme)
            sb := score(b, age, theme)
            return sa > sb
        default:
            sa := 0.5*a.PopularityScore + 0.5*ageFit(a, age)
            sb := 0.5*b.PopularityScore + 0.5*ageFit(b, age)
            return sa > sb
        }
    })
}

func sortPages(items []models.Page) {
    sort.Slice(items, func(i, j int) bool { return items[i].Index < items[j].Index })
}