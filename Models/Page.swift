import Foundation

struct PageModel: Identifiable, Codable, Equatable {
    let id: String
    let index: Int
    let imageURL: String?
    let localName: String?
}