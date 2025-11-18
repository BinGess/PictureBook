import Foundation

final class RecommendationService {
    static let shared = RecommendationService()

    func score(for book: Book, age: Int, theme: [String]) -> Double {
        let popularity = book.popularityScore
        let ageFit = book.ageFits(age)
        let themeSet = Set(theme)
        let bookSet = Set(book.themeKeywords + book.tags)
        let inter = Double(themeSet.intersection(bookSet).count)
        let union = Double(themeSet.union(bookSet).count)
        let themeSimilarity = union > 0 ? inter / union : 0
        return 0.5 * popularity + 0.3 * ageFit + 0.2 * themeSimilarity
    }

    func recommend(for current: Book, from all: [Book], age: Int, limit: Int = 5) -> [Book] {
        let base = Set(current.themeKeywords + current.tags)
        let filtered = all.filter { $0.id != current.id }
        let ranked = filtered.map { ($0, score(for: $0, age: age, theme: Array(base))) }
            .sorted { $0.1 > $1.1 }
        let top = Array(ranked.prefix(limit)).map { $0.0 }
        if top.count < limit {
            let rest = filtered.sorted { $0.popularityScore > $1.popularityScore }
            return Array((top + rest).uniqued().prefix(limit))
        }
        return top
    }
}

extension Array where Element: Hashable {
    func uniqued() -> [Element] {
        var set = Set<Element>()
        var arr: [Element] = []
        for e in self where !set.contains(e) {
            set.insert(e)
            arr.append(e)
        }
        return arr
    }
}