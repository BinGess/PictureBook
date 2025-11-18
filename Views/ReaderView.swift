import SwiftUI

struct ReaderView: View {
    let book: Book
    @State private var index: Int = 0
    @State private var showCompletion: Bool = false
    @State private var recommended: [Book] = []

    var body: some View {
        ZStack(alignment: .bottom) {
            TabView(selection: $index) {
                ForEach(Array(book.pages.enumerated()), id: \._0) { i, page in
                    pageView(page)
                        .tag(i)
                        .onAppear {
                            if i == book.pages.count - 1 {
                                computeRecommendation()
                                withAnimation { showCompletion = true }
                            }
                        }
                }
            }
            .tabViewStyle(.page(indexDisplayMode: .never))

            ProgressBar(progress: progress)
                .padding(.horizontal, 40)
                .padding(.vertical, 24)

            if showCompletion {
                CompletionOverlay(recommended: recommended) { next in
                    showCompletion = false
                    index = 0
                    navigateTo(next)
                }
                .transition(.move(edge: .bottom).combined(with: .opacity))
                .padding(.bottom, 40)
            }
        }
        .navigationTitle(book.title)
        .navigationBarTitleDisplayMode(.inline)
    }

    private var progress: Double {
        guard book.pages.count > 0 else { return 0 }
        return Double(index + 1) / Double(book.pages.count)
    }

    @ViewBuilder
    private func pageView(_ page: PageModel) -> some View {
        ZStack {
            Rectangle().fill(Color.gray.opacity(0.1))
            Text("第 \(page.index + 1) 页")
                .font(.largeTitle)
        }
        .ignoresSafeArea()
    }

    private func computeRecommendation() {
        let age = max(book.ageMin, min(book.ageMax, (book.ageMin + book.ageMax) / 2))
        Task {
            if let res: [Book] = try? await APIClient.fetchRecommendations(id: book.id, age: age, limit: 5) {
                recommended = res
            } else {
                let all = ContentService.shared.books
                recommended = RecommendationService.shared.recommend(for: book, from: all, age: age, limit: 5)
            }
        }
    }

    private func navigateTo(_ next: Book) {
        NotificationCenter.default.post(name: .navigateToBook, object: next)
    }
}

extension Notification.Name {
    static let navigateToBook = Notification.Name("navigateToBook")
}