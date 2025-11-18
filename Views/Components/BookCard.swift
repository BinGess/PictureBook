import SwiftUI

struct BookCard: View {
    let book: Book
    var body: some View {
        VStack(spacing: 8) {
            ZStack {
                Rectangle().fill(Color.gray.opacity(0.2))
                Text(book.title.prefix(1))
                    .font(.largeTitle).bold()
            }
            .frame(maxWidth: .infinity)
            .aspectRatio(4/3, contentMode: .fit)
            Text(book.title)
                .font(.headline)
                .foregroundStyle(.primary)
            HStack(spacing: 6) {
                Image(systemName: "person.fill")
                Text("\(book.ageMin)-\(book.ageMax)")
            }
            .font(.subheadline)
            .foregroundStyle(.secondary)
        }
        .padding(12)
        .background(RoundedRectangle(cornerRadius: 16).fill(Color.white))
        .shadow(color: .black.opacity(0.06), radius: 6, x: 0, y: 3)
    }
}