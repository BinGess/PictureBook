import XCTest
@testable import PictureBook

final class RecommendationTests: XCTestCase {
    func testScoreWeights() {
        let pages = [PageModel(id: "p", index: 0, imageURL: nil, localName: nil)]
        let a = Book(id: "a", title: "A", coverURL: "", ageMin: 3, ageMax: 6, tags: ["森林"], popularityScore: 0.8, themeKeywords: ["自然"], isEditorPick: false, pages: pages)
        let b = Book(id: "b", title: "B", coverURL: "", ageMin: 7, ageMax: 9, tags: ["城市"], popularityScore: 0.4, themeKeywords: ["生活"], isEditorPick: false, pages: pages)
        let s = RecommendationService.shared
        let scoreA = s.score(for: a, age: 5, theme: ["自然"])
        let scoreB = s.score(for: b, age: 5, theme: ["自然"])
        XCTAssertGreaterThan(scoreA, scoreB)
    }

    func testRecommendCount() {
        let pages = [PageModel(id: "p", index: 0, imageURL: nil, localName: nil)]
        let books = (0..<10).map { i in
            Book(id: "b\(i)", title: "B\(i)", coverURL: "", ageMin: 3, ageMax: 10, tags: ["标签\(i%3)"], popularityScore: Double.random(in: 0...1), themeKeywords: ["主题\(i%2)"], isEditorPick: i < 3, pages: pages)
        }
        let s = RecommendationService.shared
        let rec = s.recommend(for: books[0], from: books, age: 6, limit: 5)
        XCTAssertEqual(rec.count, 5)
        XCTAssertFalse(rec.contains(where: { $0.id == books[0].id }))
    }
}