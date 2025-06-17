# Contributing to VietRedis Server

Cảm ơn bạn đã quan tâm đến việc đóng góp cho VietRedis Server! Đây là dự án cá nhân của Namxely, nhưng tôi rất hoan nghênh mọi đóng góp và feedback từ cộng đồng.

## 👨‍💻 Về dự án

VietRedis Server là một implementation Redis được phát triển bởi **Namxely** ([@namxely](https://github.com/namxely)), tập trung vào:
- Hiệu năng cao cho thị trường Việt Nam
- Vietnamese developer community
- Learning và portfolio showcase
- Production-ready features

## 🤝 Cách đóng góp

### Báo cáo Bug

Nếu bạn tìm thấy bug, vui lòng tạo một issue với thông tin sau:

- **Mô tả bug**: Mô tả rõ ràng và ngắn gọn về bug
- **Cách tái tạo**: Các bước để tái tạo bug
  1. Bước 1
  2. Bước 2
  3. Bước 3...
- **Kết quả mong đợi**: Kết quả bạn mong đợi sẽ xảy ra
- **Kết quả thực tế**: Kết quả thực sự xảy ra
- **Environment**: 
  - OS: [e.g. Windows 11, Ubuntu 20.04]
  - Go version: [e.g. 1.18.0]
  - Version: [e.g. v1.0.0]

### Đề xuất Feature

Để đề xuất feature mới:

- Tạo issue với label "enhancement"
- Mô tả rõ ràng về feature bạn muốn
- Giải thích tại sao feature này hữu ích
- Cung cấp ví dụ về cách sử dụng (nếu có)

### Pull Requests

1. **Fork repository**
2. **Tạo branch mới** từ `main`:
   ```bash
   git checkout -b feature/ten-feature-moi
   ```
3. **Thực hiện thay đổi** với code chất lượng cao
4. **Viết tests** cho code mới (nếu áp dụng)
5. **Chạy existing tests** để đảm bảo không breaking:
   ```bash
   go test ./...
   ```
6. **Format code**:
   ```bash
   go fmt ./...
   ```
7. **Commit thay đổi** với commit message rõ ràng:
   ```bash
   git commit -m "Add: Thêm feature xyz"
   ```
8. **Push lên branch**:
   ```bash
   git push origin feature/ten-feature-moi
   ```
9. **Tạo Pull Request** với mô tả chi tiết

## Code Style

### Go Conventions

- Tuân theo [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Sử dụng `go fmt` để format code
- Sử dụng `go vet` để check code
- Viết comment cho public functions/methods
- Sử dụng error handling đúng cách

### Commit Messages

Sử dụng format sau cho commit message:

```
Type: Mô tả ngắn gọn (dưới 50 ký tự)

Mô tả chi tiết hơn (nếu cần)
- Chi tiết 1
- Chi tiết 2

Fixes #issue_number (nếu có)
```

**Types:**
- `Add`: Thêm feature mới
- `Fix`: Sửa bug  
- `Update`: Cập nhật existing feature
- `Remove`: Xóa code/feature
- `Refactor`: Refactor code
- `Docs`: Cập nhật documentation
- `Test`: Thêm/sửa tests

### Project Structure

Khi thêm file mới, hãy đặt vào thư mục phù hợp:

```
├── main.go              # Entry point
├── config/              # Configuration management
├── database/            # Core database implementation
├── redis/               # Redis protocol implementation
├── cluster/             # Cluster mode implementation
├── datastruct/          # Data structures
├── interface/           # Interface definitions
├── lib/                 # Utility libraries
├── pubsub/              # Pub/Sub implementation
├── tcp/                 # TCP server implementation
└── aof/                 # AOF persistence
```

## Development Setup

### Prerequisites

- Go 1.18 hoặc mới hơn
- Git
- Redis CLI (để testing)

### Local Development

```bash
# Clone your fork
git clone https://github.com/namxely/Build-your-own-viet-redis.git
cd Build-your-own-viet-redis

# Install dependencies
go mod tidy

# Build
go build -o my-redis-server

# Run
./my-redis-server

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

### Testing

- Viết unit tests cho logic mới
- Viết integration tests cho features lớn
- Sử dụng table-driven tests khi phù hợp
- Test coverage nên > 70%

### Documentation

- Cập nhật README.md khi thêm feature mới
- Thêm comments cho complex logic
- Cập nhật CHANGELOG.md
- Thêm examples trong documentation

## Review Process

1. **Automated checks**: CI sẽ chạy tests và checks
2. **Code review**: Maintainer sẽ review code
3. **Testing**: Nên test manually nếu cần
4. **Merge**: Sau khi approved, PR sẽ được merge

## Questions?

Nếu có câu hỏi về contributing, hãy tạo issue với label "question" hoặc liên hệ:

- Email: dev.namxely@gmail.com
- GitHub Issues: [Create an issue](https://github.com/namxely/Build-your-own-viet-redis/issues)
- Telegram: @NamxelyDev

## Code of Conduct

- Respectful communication
- Constructive feedback
- Help each other learn
- Focus on the code, not the person

## License

Bằng cách contribute, bạn đồng ý rằng contributions sẽ được licensed dưới GPL-3.0 License.

---

Cảm ơn bạn đã đóng góp! 🚀
