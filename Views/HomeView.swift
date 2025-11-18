import SwiftUI

struct HomeView: View {
    @State private var editorPicks: [Book] = []
    @State private var items: [Book] = []
    @State private var page: Int = 1
    private let pageSize = 8

    private var pagedBooks: [Book] { items }

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 24) {
                if !editorPicks.isEmpty {
                    editorSection
                }
                gridSection
            }
            .padding(24)
        }
        .ignoresSafeArea(edges: .horizontal)
        .task {
            await ContentService.shared.loadRemoteHome(page: 1, pageSize: 24, age: nil)
            editorPicks = ContentService.shared.editorPicks
            items = ContentService.shared.books
        }
        .navigationTitle("绘本")
    }

    private var editorSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("编辑推荐")
                .font(.title2).bold()
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 16) {
                    ForEach(editorPicks) { book in
                        NavigationLink(value: book) {
                            BookCard(book: book).frame(width: 320)
                        }
                    }
                }
                .padding(.horizontal, 8)
            }
        }
    }

    private var gridSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("为你推荐")
                .font(.title2).bold()
            LazyVGrid(columns: Array(repeating: GridItem(.flexible(), spacing: 16), count: 3), spacing: 16) {
                ForEach(pagedBooks) { book in
                    NavigationLink(value: book) {
                        BookCard(book: book)
                    }
                    .onAppear {
                        if book.id == pagedBooks.last?.id {
                            Task { await loadMoreIfNeeded() }
                        }
                    }
                }
            }
        }
        .navigationDestination(for: Book.self) { book in
            ReaderView(book: book)
        }
    }

    private func loadMoreIfNeeded() async {
        if ContentService.shared.hasMore {
            page += 1
            await ContentService.shared.loadMore(page: page, pageSize: 24, age: nil)
            items = ContentService.shared.books
        }
    }
}