import Foundation

struct APIClient {
    static var baseURL: URL = {
        if let urlStr = ProcessInfo.processInfo.environment["BASE_URL"], let url = URL(string: urlStr) { return url }
        return URL(string: "http://localhost:8080")!
    }()

    static func get<T: Decodable>(_ path: String, query: [URLQueryItem] = []) async throws -> T {
        var comps = URLComponents(url: baseURL.appendingPathComponent(path), resolvingAgainstBaseURL: false)!
        if !query.isEmpty { comps.queryItems = query }
        let req = URLRequest(url: comps.url!)
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard let http = resp as? HTTPURLResponse, (200..<300).contains(http.statusCode) else { throw URLError(.badServerResponse) }
        return try JSONDecoder().decode(T.self, from: data)
    }

    static func fetchEditorPicks() async throws -> [Book] {
        try await get("/v1/editor-picks")
    }

    struct BooksResponse: Decodable { let items: [Book]; let paging: Paging }
    struct Paging: Decodable { let page: Int; let page_size: Int; let has_more: Bool }

    static func fetchBooks(page: Int, pageSize: Int, age: Int?, sort: String?) async throws -> BooksResponse {
        var q: [URLQueryItem] = [URLQueryItem(name: "page", value: String(page)), URLQueryItem(name: "page_size", value: String(pageSize))]
        if let age = age { q.append(URLQueryItem(name: "age", value: String(age))) }
        if let sort = sort { q.append(URLQueryItem(name: "sort", value: sort)) }
        return try await get("/v1/books", query: q)
    }

    static func fetchBook(id: String) async throws -> Book { try await get("/v1/books/\(id)") }

    static func fetchRecommendations(id: String, age: Int?, limit: Int = 5) async throws -> [Book] {
        var q: [URLQueryItem] = [URLQueryItem(name: "limit", value: String(limit))]
        if let age = age { q.append(URLQueryItem(name: "age", value: String(age))) }
        return try await get("/v1/books/\(id)/recommendations", query: q)
    }
}