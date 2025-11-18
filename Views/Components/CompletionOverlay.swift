import SwiftUI

struct CompletionOverlay: View {
    let recommended: [Book]
    let onSelect: (Book) -> Void
    var body: some View {
        VStack(spacing: 16) {
            Text("你已读完啦")
                .font(.largeTitle).bold()
            Text("为你推荐这些绘本")
                .font(.headline)
                .foregroundStyle(.secondary)
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 12) {
                    ForEach(recommended) { book in
                        Button {
                            onSelect(book)
                        } label: {
                            BookCard(book: book)
                                .frame(width: 240)
                        }
                    }
                }
                .padding(.horizontal, 24)
            }
        }
        .padding(24)
        .frame(maxWidth: .infinity)
        .background(.ultraThickMaterial)
        .clipShape(RoundedRectangle(cornerRadius: 24))
        .shadow(radius: 10)
        .padding(.horizontal, 40)
    }
}