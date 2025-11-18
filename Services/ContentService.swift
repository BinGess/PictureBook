import Foundation

final class ContentService {
    static let shared = ContentService()

    private(set) var books: [Book] = []
    private(set) var editorPicks: [Book] = []
    private(set) var hasMore: Bool = false

    func loadLocal() {
        guard let url = Bundle.main.url(forResource: "books", withExtension: "json") else {
            books = SampleData.books
            return
        }
        do {
            let data = try Data(contentsOf: url)
            let decoded = try JSONDecoder().decode([Book].self, from: data)
            books = decoded
        } catch {
            books = SampleData.books
        }
    }

    func loadRemoteHome(page: Int = 1, pageSize: Int = 24, age: Int? = nil) async {
        do {
            let picks = try await APIClient.fetchEditorPicks()
            editorPicks = Array(picks.prefix(5))
        } catch {
            editorPicks = Array(SampleData.books.filter { $0.isEditorPick }.prefix(5))
        }
        do {
            let resp = try await APIClient.fetchBooks(page: page, pageSize: pageSize, age: age, sort: "popular")
            books = resp.items
            hasMore = resp.paging.has_more
        } catch {
            loadLocal()
            hasMore = false
        }
    }

    func loadMore(page: Int, pageSize: Int = 24, age: Int? = nil) async {
        guard hasMore else { return }
        do {
            let resp = try await APIClient.fetchBooks(page: page, pageSize: pageSize, age: age, sort: "popular")
            books.append(contentsOf: resp.items)
            hasMore = resp.paging.has_more
        } catch {
            hasMore = false
        }
    }
}

enum SampleData {
    static let books: [Book] = {
        let pagesA = (0..<8).map { PageModel(id: "A-\($0)", index: $0, imageURL: nil, localName: "a_\($0)") }
        let pagesB = (0..<6).map { PageModel(id: "B-\($0)", index: $0, imageURL: nil, localName: "b_\($0)") }
        let pagesC = (0..<10).map { PageModel(id: "C-\($0)", index: $0, imageURL: nil, localName: "c_\($0)") }
        return [
            Book(id: "book_a", title: "森林探险", coverURL: "cover_a", ageMin: 3, ageMax: 6, tags: ["森林","动物"], popularityScore: 0.82, themeKeywords: ["自然","探险"], isEditorPick: true, pages: pagesA),
            Book(id: "book_b", title: "海洋之旅", coverURL: "cover_b", ageMin: 5, ageMax: 9, tags: ["海洋","旅行"], popularityScore: 0.74, themeKeywords: ["自然","探索"], isEditorPick: true, pages: pagesB),
            Book(id: "book_c", title: "太空奇遇记", coverURL: "cover_c", ageMin: 6, ageMax: 10, tags: ["太空","科学"], popularityScore: 0.91, themeKeywords: ["宇宙","冒险"], isEditorPick: false, pages: pagesC),
            Book(id: "book_d", title: "农场的一天", coverURL: "cover_d", ageMin: 3, ageMax: 7, tags: ["农场","动物"], popularityScore: 0.65, themeKeywords: ["生活","自然"], isEditorPick: false, pages: pagesB),
            Book(id: "book_e", title: "恐龙乐园", coverURL: "cover_e", ageMin: 4, ageMax: 8, tags: ["恐龙","历史"], popularityScore: 0.88, themeKeywords: ["史前","冒险"], isEditorPick: true, pages: pagesA),
            Book(id: "book_f", title: "城市探秘", coverURL: "cover_f", ageMin: 5, ageMax: 10, tags: ["城市","建筑"], popularityScore: 0.59, themeKeywords: ["生活","探索"], isEditorPick: false, pages: pagesA)
        ]
    }()
}