import Foundation

struct Book: Identifiable, Codable, Equatable {
    let id: String
    let title: String
    let coverURL: String
    let ageMin: Int
    let ageMax: Int
    let tags: [String]
    let popularityScore: Double
    let themeKeywords: [String]
    let isEditorPick: Bool
    let pages: [PageModel]
}

extension Book {
    func ageFits(_ age: Int) -> Double {
        if age < ageMin { return 0 }
        if age > ageMax { return 0 }
        let span = max(ageMax - ageMin, 1)
        let pos = Double(age - ageMin) / Double(span)
        return 1 - abs(pos - 0.5) * 2
    }
}