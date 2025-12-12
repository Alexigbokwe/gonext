# GoNext Framework Documentation

**The Scalable, Modular, and Developer-First Go Framework.**

GoNext is an opinionated framework designed to bring the elegance and structure of modern frameworks like NestJS and Laravel to the Go ecosystem. It leverages the raw performance of **Fiber** while providing a robust architecture for building scalable applications.

---

## 📚 Table of Contents
1.  [Quick Start](#-quick-start)
2.  [Architecture](#-architecture)
3.  [Dependency Injection](#-dependency-injection)
4.  [Configuration](#-configuration)
5.  [Database & ORM](#-database--orm)
6.  [Routing & Controllers](#-routing--controllers)
7.  [Middleware](#-middleware)
8.  [Security](#-security)
9.  [Queue & Background Jobs](#-queue--background-jobs)
10. [Observability](#-observability)
11. [CLI Reference](#-cli-reference)

---

## 🚀 Quick Start

### Installation
Ensure you have Go 1.20+ installed.

```bash
go install github.com/Alexigbokwe/gonext@latest
```

### Creating a New Project
Bootstrap a new application in seconds:

```bash
gonext new my-awesome-project
cd my-awesome-project
go mod tidy
go run main.go
```

The server will start on `http://localhost:5050`.

---

## 🏗 Architecture

GoNext uses a **Modular Architecture**. Instead of organizing files by type (controllers, services), we organize them by **Domain Feature** (e.g., `UserModule`, `ProductModule`).

### Directory Structure
```
├── app/                  # Core Framework Code
├── config/               # Configuration Structs & Logic
├── cmd/                  # Application Entry Points
├── app/                  # Your Application Logic
│   └── user/             # User Module
│       ├── controller/   # HTTP Handlers
│       ├── service/      # Business Logic
│       ├── repository/   # Data Access
│       ├── dto/          # Data Transfer Objects
│       └── module.go     # Module Wiring
├── main.go               # Entry Point
└── go.mod
```

### Modules
A Module is the container for your feature. It registers your components with the DI container.

```go
// app/user/module.go
func (m *UserModule) Register(container *app.Container) {
    container.Register(&service.UserService{})
    container.Register(&controller.UserController{})
}

func (m *UserModule) MountRoutes(router fiber.Router) {
    router.Get("/users", m.Controller.GetUsers)
}
```

---

## 💉 Dependency Injection

GoNext features a powerful, reflection-based DI container.

### Lifetimes
*   **Singleton**: One instance per app (stateless services).
*   **Scoped**: One instance per request (context-aware services).
*   **Transient**: New instance per injection.

### Usage
Simply tag your struct fields with \`inject:"type"\` or \`inject:"token"\`.

```go
type UserService struct {
    // Injects by type
    Repo *repository.UserRepository `inject:"type"`
    
    // Injects by token (e.g., generic interface)
    Queue queue.TaskQueue `inject:"queue"`
}
```

---

## ⚙️ Configuration

We support strict, type-safe configuration via **Viper**.

*   **File**: `config/appConfig.go`
*   **Source**: `.env` file or Environment Variables.

### Defining Config
```go
type Config struct {
    Server ServerConfig
    Queue  QueueConfig
}
```

### Injecting Config
```go
type MyService struct {
    Config *config.Config `inject:"type"`
}
```

---

## 🗄️ Database & ORM

GoNext is **Database Agnostic**. We believe you should choose the best tool for the job.

*   **Recommended**: `pgx` (Performance), `GORM` (Productivity), `Ent` (Type-safety).
*   **Setup**: Initialize your DB in `main.go` and register it in the container.

```go
// main.go
func main() {
    db := connectToDatabase() 
    container.Register(db) // Available to all repositories
}
```

---

## 🌐 Routing & Controllers

Controllers handle incoming HTTP requests. They should remain thin and delegate logic to Services.

```go
// @Summary Create User
// @Router /users [post]
func (c *UserController) Create(ctx *fiber.Ctx) error {
    var payload dto.CreateUserDTO
    if err := ctx.BodyParser(&payload); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }
    
    // Auto-Validation
    if err := app.ValidateStruct(payload); err != nil {
        return ctx.Status(422).JSON(err)
    }
    
    return c.Service.Create(payload)
}
```

---

---

## ⚡ Concurrency & Async

GoPro brings the ergonomics of `async/await` to Go, making concurrent programming intuitive and safe.

### The `Async` Helper
Use `utils.Async` to run a blocking operation in a goroutine and return a generic `Promise[T]`.

```go
import "goNext/app/utils"

func (s *Service) FetchData() (string, error) {
    // Start task
    promise := utils.Async(func(ctx context.Context) (string, error) {
        // ... some heavy work ...
        return "result", nil
    })

    // Do other work here...

    // Await result (blocks until done)
    return promise.Await()
}
```

### Flow Control
*   **Timeout**: `promise.WithTimeout(5 * time.Second)` automatically cancels the context if time is exceeded.
*   **Cancel**: Manually call `promise.Cancel()` to stop execution.

### Promise Combinators
*   **PromiseAll**: Wait for all tasks to succeed. Fails fast if one errors.
*   **PromiseResult**: Wait for the first task to complete (success or error).
*   **PromiseAllSettled**: Wait for all tasks, returning a list of successes and failures.

```go
p1 := utils.Async(fetchUser)
p2 := utils.Async(fetchOrders)

results, err := utils.PromiseAll(p1, p2)
```

### Fire-and-Forget
For lightweight background tasks that don't need persistence (like internal logging), use `RunBackground`.

```go
utils.RunBackground(func() error {
    // This runs in a goroutine with panic recovery
    return sendAnalytics()
})
```

---

## 🛡️ Security

GoNext aims to be secure by default.

### Features
*   **Authentication**: Built-in `JwtService` and `AuthGuard`.
*   **Headers**: Automated security headers via `Helmet`.
*   **CORS**: Configurable Cross-Origin Resource Sharing.
*   **Rate Limiting**: Built-in DDOS protection.

### Protecting Routes
```go
func (m *Module) MountRoutes(r fiber.Router) {
    guard := &security.AuthGuard{}
    m.Container.MustAutowire(guard)
    
    r.Get("/profile", guard.Middleware(), m.Controller.Profile)
}
```

---

## 📨 Queue & Background Jobs

GoNext provides a `TaskQueue` interface but **does not enforce a specific implementation**. You are free to plug in Redis, RabbitMQ, Kafka, or AWS SQS.

### 1. The Interface
Every queue provider must implement this interface:

```go
// app/queue/queue.go
type TaskQueue interface {
    Enqueue(typeName string, payload interface{}, opts ...interface{}) error
    RegisterHandler(typeName string, handler func(payload []byte) error)
    Start()
    Shutdown()
}
```

### 2. Implementation Recipes

Here are full implementations you can copy into your project.

#### Option A: Redis (Recommended)

1.  **Install**: `go get github.com/hibiken/asynq`
2.  **File**: `app/queue/redis.go`

```go
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"goNext/app/logger"
	"goNext/config"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type RedisTaskQueue struct {
	client *asynq.Client
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewRedisTaskQueue(cfg *config.Config) *RedisTaskQueue {
	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	if redisAddr == ":" {
		redisAddr = "localhost:6379"
	}

	redisOpt := asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password}

	return &RedisTaskQueue{
		client: asynq.NewClient(redisOpt),
		server: asynq.NewServer(redisOpt, asynq.Config{
			Concurrency: 10,
			Logger:      logger.Log.Sugar(),
		}),
		mux: asynq.NewServeMux(),
	}
}

func (q *RedisTaskQueue) Enqueue(typeName string, payload interface{}, opts ...interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(typeName, data)
	info, err := q.client.Enqueue(task)
	if err != nil {
		return err
	}
	logger.Log.Info("Enqueued task", zap.String("id", info.ID), zap.String("type", typeName))
	return nil
}

func (q *RedisTaskQueue) RegisterHandler(typeName string, handler func(payload []byte) error) {
	q.mux.HandleFunc(typeName, func(ctx context.Context, t *asynq.Task) error {
		return handler(t.Payload())
	})
}

func (q *RedisTaskQueue) Start() {
	go func() {
		if err := q.server.Run(q.mux); err != nil {
			logger.Log.Fatal("Queue server failed", zap.Error(err))
		}
	}()
}

func (q *RedisTaskQueue) Shutdown() {
	q.client.Close()
	q.server.Stop()
}
```

#### Option B: RabbitMQ

1.  **Install**: `go get github.com/rabbitmq/amqp091-go`
2.  **File**: `app/queue/rabbitmq.go`

```go
package queue

import (
	"encoding/json"
	"goNext/app/logger"
	"goNext/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQQueue struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	handlers  map[string]func([]byte) error
}

func NewRabbitQueue(url string, queueName string) *RabbitMQQueue {
	conn, _ := amqp.Dial(url)
	ch, _ := conn.Channel()
	ch.QueueDeclare(queueName, true, false, false, false, nil)
	
	return &RabbitMQQueue{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
		handlers:  make(map[string]func([]byte) error),
	}
}

func (r *RabbitMQQueue) Enqueue(typeName string, payload interface{}, opts ...interface{}) error {
	body, _ := json.Marshal(map[string]interface{}{
		"type": typeName, "payload": payload,
	})
	return r.channel.Publish("", r.queueName, false, false, amqp.Publishing{
		ContentType: "application/json", Body: body,
	})
}

func (r *RabbitMQQueue) RegisterHandler(typeName string, handler func(payload []byte) error) {
	r.handlers[typeName] = handler
}

func (r *RabbitMQQueue) Start() {
	msgs, _ := r.channel.Consume(r.queueName, "", true, false, false, false, nil)
	go func() {
		for d := range msgs {
			var msg struct { Type string; Payload json.RawMessage }
			json.Unmarshal(d.Body, &msg)
			if h, ok := r.handlers[msg.Type]; ok {
				h(msg.Payload)
			}
		}
	}()
}

func (r *RabbitMQQueue) Shutdown() {
	r.channel.Close()
	r.conn.Close()
}
```

#### Option C: Kafka

1.  **Install**: `go get github.com/segmentio/kafka-go`
2.  **Implementation**: Similar structure, using `kafka.Writer` and `kafka.Reader`. Wrap your payload in a struct that includes the task type, so your single consumer loop can dispatch to the correct handler.

### 3. Registration
In your `main.go`, initialize your chosen provider and bind it to the interface.

```go
func main() {
    // ...
    // Initialize Redis Queue
    redisQueue := queue.NewRedisQueue(config)
    redisQueue.Start()
    defer redisQueue.Shutdown()
    
    // Bind to Interface
    container.Bind("queue", redisQueue)
}
```

### 4. Usage in Services
```go
type EmailService struct {
    Queue queue.TaskQueue `inject:"queue"`
}

func (s *EmailService) Send() {
    s.Queue.Enqueue("email:send", map[string]string{"to": "foo@bar.com"})
}
```

---

---

## 📢 Events (Pub/Sub)

Decouple your modules using the built-in Event Bus.

### Usage

1.  **Define an Event**:
    ```go
    type UserRegisteredParams struct {
        UserID string
    }
    // Implement events.Event interface (Name() string)
    type UserRegistered struct {
        Payload UserRegisteredParams
    }
    func (e UserRegistered) Name() string { return "UserRegistered" }
    ```

2.  **Register a Listener**:
    in your `module.go` or `main.go`:
    ```go
    dispatcher := events.GetDispatcher()
    dispatcher.Register("UserRegistered", func(ctx context.Context, e events.Event) error {
        event := e.(UserRegistered)
        fmt.Println("User registered:", event.Payload.UserID)
        return nil
    })
    ```

3.  **Dispatch**:
    ```go
    dispatcher.Dispatch(ctx, UserRegistered{Payload: p})
    // Or Async
    dispatcher.DispatchAsync(UserRegistered{Payload: p})
    ```

---

## 📦 Caching

Unified API for caching throughout your app.

### Configuration
GoNext attempts to use **Redis** if configured in `.env`. Otherwise, it falls back to **Memory**.

```env
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Usage
Inject `cache.Store` into your service.

```go
type ProductService struct {
    Cache cache.Store `inject:"cache"` // Binds to Memory or Redis
}

func (s *ProductService) GetProduct(id string) {
    // Set
    s.Cache.Set(ctx, "product:"+id, product, 5*time.Minute)
    
    // Get
    var p Product
    s.Cache.Get(ctx, "product:"+id, &p)
    
    // Forget
    s.Cache.Forget(ctx, "product:"+id)
}
```

---

## ⏰ Task Scheduler

Run recurring jobs (Cron) with ease.

### Usage
The scheduler is started automatically in `main.go`. You can inject it or use the singleton.

```go
// In main.go or module registration
cronScheduler.Add("@every 1h", func() {
    logger.Log.Info("Cleaning up temporary files...")
    utils.Cleanup()
})
```

Supported specs:
*   `@every 1h30m`
*   `0 30 * * * *` (Standard Cron: sec min hour dom month dow)

---

## 📧 Mailer

Send transactional emails via SMTP.

### Configuration
```env
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=user
MAIL_PASSWORD=pass
MAIL_FROM_ADDRESS=noreply@example.com
```

### Usage
Inject `mail.Mailer`.

```go
type AuthService struct {
    Mailer mail.Mailer `inject:"mail"`
}

func (s *AuthService) SendWelcome() {
    s.Mailer.Send(
        []string{"user@example.com"},
        "Welcome to GoNext!",
        "Thanks for signing up.",
    )
}
```

---

---

## 📂 File Storage

Unified API for file operations. Supports **Local Disk** and easily extensible for **S3**.

### Usage
Inject `storage.Disk`.

```go
type DocumentService struct {
    Disk storage.Disk `inject:"storage"`
}

func (s *DocumentService) Upload(file *multipart.FileHeader) {
    // Saves to ./public/uploads/file.png
    url, err := s.Disk.Put(file, "uploads")
    
    // Check existence
    exists := s.Disk.Exists("uploads/file.png")
    
    // Get URL
    // http://localhost:5050/public/uploads/file.png
    link := s.Disk.URL(url)
}
```

---

## 🔐 Hashing

Standardized password hashing using **Bcrypt**.

### Usage
Inject `security.HashService`.

```go
type AuthService struct {
    Hasher security.HashService `inject:"hasher"`
}

func (s *AuthService) Register(password string) {
    // Hash
    hashed, _ := s.Hasher.Hash(password)
    
    // Compare
    match := s.Hasher.Compare(hashed, "secret")
}
```

---

## ✅ Validation

GoNext uses `go-playground/validator` but enhances it with **Friendly Error Messages**.

### Usage
In your controller:

```go
func (c *UserController) Create(ctx *fiber.Ctx) error {
    var dto CreateUserDTO
    ctx.BodyParser(&dto)
    
    // Returns detailed array of errors or nil
    if errors := app.ValidateStruct(dto); errors != nil {
         return ctx.Status(422).JSON(fiber.Map{
             "message": "Validation failed",
             "errors": errors,
         })
    }
    // ...
}
```

---

## 🔭 Observability

*   **Metrics**: Visit `/metrics` to see real-time request stats, latency distributions, and memory usage.
*   **Logging**: Structured JSON logs (Zap) allow for easy ingestion into ELK/Datadog.

---

## 🛠 CLI Reference

| Command | Description |
| :--- | :--- |
| `gonext new <name>` | Scaffolds a new project. |
| `gonext g module <name>` | Generates a new Domain Module. |
| `gonext g controller <name> <mod>` | Generates a Controller. |
| `gonext g service <name> <mod>` | Generates a Service. |
| `gonext g repository <name> <mod>` | Generates a Repository. |
| `gonext g dto <name> <mod>` | Generates a DTO struct. |
| `gonext doc` | Opens this offline documentation. |

---

*Documentation generated by GoNext CLI*
