import SwiftUI

@main
struct BookApp: App {
    @State private var path = NavigationPath()
    var body: some Scene {
        WindowGroup {
            NavigationStack(path: $path) {
                HomeView()
                    .navigationDestination(for: Book.self) { book in
                        ReaderView(book: book)
                    }
            }
            .navigationViewStyle(.stack)
            .onReceive(NotificationCenter.default.publisher(for: .navigateToBook)) { note in
                if let next = note.object as? Book {
                    path.append(next)
                }
            }
        }
    }
}