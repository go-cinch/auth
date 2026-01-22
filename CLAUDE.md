# Development Guidelines

## 1. Logging with Context

When logging in **any** method that has a `context.Context` parameter, **always** use `log.WithContext(ctx)` to ensure trace correlation and proper context propagation. This applies to all layers: `service`, `biz`, `data`, `middleware`, `task`, etc.

```go
// Good - use WithContext when ctx is available
func (s *GameService) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (*emptypb.Empty, error) {
    // ...
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("create game failed")
        return nil, err
    }
    return &emptypb.Empty{}, nil
}

func (uc *GameUseCase) Create(ctx context.Context, item *CreateGame) error {
    // ...
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("create game failed")
        return err
    }
    return nil
}

func (ro gameRepo) Create(ctx context.Context, item *biz.CreateGame) error {
    // ...
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("create game failed")
        return err
    }
    return nil
}

// Bad - missing context in log
func (ro gameRepo) Create(ctx context.Context, item *biz.CreateGame) error {
    // ...
    if err != nil {
        log.WithError(err).Error("create game failed")  // Missing WithContext(ctx)
        return err
    }
    return nil
}
```

**Note:** `log.WithContext(ctx)` enables:
- Trace ID correlation with OpenTelemetry spans
- Request-scoped logging metadata
- Distributed tracing visibility

## 2. OpenTelemetry Tracing

All public methods in `service`, `biz`, and `data` layers should include OpenTelemetry tracing:

```go
func (s *SomeService) MethodName(ctx context.Context, req *Request) (*Response, error) {
    tr := otel.Tracer("service")  // or "biz" or "data" depending on layer
    ctx, span := tr.Start(ctx, "MethodName")
    defer span.End()

    // method implementation...
}
```

**Tracer names by layer:**
- Service layer: `otel.Tracer("service")`
- Business layer: `otel.Tracer("biz")`
- Data layer: `otel.Tracer("data")`

**Span naming:** Use only the method name (e.g., `"Create"`, `"Find"`, `"Update"`), not the full qualified name.

## 3. Service Layer Parameter Handling

Use `copierx.Copy` to handle request and response parameters in the service layer. Avoid verbose manual field-by-field assignments.

```go
// Good - use copierx.Copy for Create
func (s *GameService) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (*emptypb.Empty, error) {
    r := &biz.CreateGame{}
    copierx.Copy(r, req)

    if err := s.game.Create(ctx, r); err != nil {
        return nil, err
    }
    return &emptypb.Empty{}, nil
}

// Good - use copierx.Copy for Update
func (s *GameService) UpdateGame(ctx context.Context, req *v1.UpdateGameRequest) (*emptypb.Empty, error) {
    r := &biz.UpdateGame{}
    copierx.Copy(r, req)
    if err := s.game.Update(ctx, r); err != nil {
        return nil, err
    }
    return &emptypb.Empty{}, nil
}

// Bad - manual field assignments
func (s *GameService) UpdateGame(ctx context.Context, req *v1.UpdateGameRequest) (*emptypb.Empty, error) {
    r := &biz.UpdateGame{
        ID:   req.GetId(),
        Name: req.Name,
        Type: req.Type,
    }
    // ... more manual assignments
}
```

**For paginated queries:**
```go
var p page.Page
r := &biz.FindGame{}
copierx.Copy(&p, req.GetPage())
copierx.Copy(r, req)
```

## 4. No Nil Checks for Request Parameters

Service layer request parameters (`req`) generally do not need nil checks. The framework ensures requests are valid before reaching the handler.

```go
// Good - no nil check needed
func (s *GameService) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (*emptypb.Empty, error) {
    r := &biz.CreateGame{}
    copierx.Copy(r, req)
    // ...
}

// Bad - unnecessary nil check
func (s *GameService) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (*emptypb.Empty, error) {
    if req == nil {
        return nil, ErrIllegalParameter(ctx, "req")
    }
    // ...
}
```

## 5. No Nil Checks for Wire-Injected Dependencies

Dependencies injected via Wire in `service`, `biz`, and `data` layers do not need nil checks. Wire guarantees all dependencies are properly initialized at startup.

```go
// Good - trust Wire injection
func (s *GameService) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (*emptypb.Empty, error) {
    return s.game.Create(ctx, r)  // s.game is guaranteed by Wire
}

// Bad - unnecessary nil check for Wire-injected dependency
func (s *GameService) CreateGame(ctx context.Context, req *v1.CreateGameRequest) (*emptypb.Empty, error) {
    if s.game == nil {
        return nil, ErrInternal(ctx, "game usecase not configured")
    }
    return s.game.Create(ctx, r)
}
```

## 6. Use copierx.Copy for List Results

In Find methods, use `copierx.Copy` to copy list results directly instead of manual for loop iteration.

```go
// Good - use copierx.Copy for list
func (s *GameService) FindGame(ctx context.Context, req *v1.FindGameRequest) (*v1.FindGameReply, error) {
    // ... business logic call ...

    rp := &v1.FindGameReply{
        Page: &params.Page{},
    }
    copierx.Copy(&rp.Page, &r.Page)
    copierx.Copy(&rp.List, list)
    return rp, nil
}

// Bad - manual for loop iteration
func (s *GameService) FindGame(ctx context.Context, req *v1.FindGameRequest) (*v1.FindGameReply, error) {
    // ... business logic call ...

    rp := &v1.FindGameReply{
        Page: &params.Page{},
        List: make([]*v1.Game, 0, len(list)),
    }
    _ = copierx.Copy(&rp.Page, &r.Page)
    for i := range list {
        rp.List = append(rp.List, gameToPB(&list[i]))
    }
    return rp, nil
}
```

## 7. Method Naming Convention for Data Access

In `service`, `biz`, and `data` layers, use consistent prefixes for methods that access data:

- **`Find`**: Returns a **list** of records (zero or more items)
- **`Get`**: Returns a **single** record (exactly one item)

```go
// Good - Find returns list
func (ro gameRepo) Find(ctx context.Context, condition *biz.FindGame) []biz.Game
func (ro gameRepo) FindByIDs(ctx context.Context, ids []int64) []*biz.Game
func (ro gameRepo) FindByCode(ctx context.Context, codes string) []biz.Game

// Good - Get returns single item
func (ro gameRepo) GetByID(ctx context.Context, id int64) (*biz.Game, error)
func (ro gameRepo) GetByUsername(ctx context.Context, username string) (*biz.User, error)

// Good - Find returns list (even for tree structures)
func (ro actionRepo) FindTree(ctx context.Context) ([]*biz.Action, error)
func (ro permissionRepo) FindUserPermissions(ctx context.Context, userID int64) ([]*biz.Action, error)

// Bad - inconsistent naming
func (ro gameRepo) GetByIDs(ctx context.Context, ids []int64) []*biz.Game  // Should be FindByIDs
func (ro gameRepo) FindByID(ctx context.Context, id int64) (*biz.Game, error)  // Should be GetByID
```

## 8. Data Layer CRUD with gorm.G Generic API

Data layer methods **must** use `gorm.G[model.XXX]` generic API for database operations. This provides type safety and consistent query building.

### Create

```go
func (ro gameRepo) Create(ctx context.Context, item *biz.CreateGame) (err error) {
    db := gorm.G[model.Game](ro.data.DB(ctx))

    // Check uniqueness for unique fields
    count, err := db.Where("name = ?", item.Name).Count(ctx, "*")
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("check name exists failed")
        return
    }
    if count > 0 {
        return biz.ErrDuplicateField(ctx, "name", item.Name)
    }

    // Generate ID if not provided
    if item.ID == 0 {
        item.ID = ro.data.ID(ctx)
    }

    // Copy to model and create
    var m model.Game
    copierx.Copy(&m, item)
    err = db.Create(ctx, &m)
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("create game failed")
    }
    return
}
```

### Get (single record)

```go
func (ro gameRepo) Get(ctx context.Context, id int64) (item *biz.Game, err error) {
    db := gorm.G[model.Game](ro.data.DB(ctx))
    item = &biz.Game{}

    m, err := db.Where("id = ?", id).First(ctx)
    if err == gorm.ErrRecordNotFound {
        return item, biz.ErrRecordNotFound(ctx)
    }
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("get game failed")
        return
    }
    copierx.Copy(item, m)
    return
}
```

### Find (list with pagination)

```go
func (ro gameRepo) Find(ctx context.Context, condition *biz.FindGame) (rp []biz.Game) {
    rp = make([]biz.Game, 0)
    db := gorm.G[model.Game](ro.data.DB(ctx))

    // Apply filters
    if condition.Name != nil {
        db.Where("name LIKE ?", "%"+*condition.Name+"%")
    }

    // Count total before pagination
    if !condition.Page.Disable {
        count, err := db.Count(ctx, "*")
        if err != nil {
            log.WithContext(ctx).WithError(err).Error("count game failed")
            return
        }
        condition.Page.Total = count
        if count == 0 {
            return
        }
    }

    // Apply ordering and pagination
    db.Order("id DESC")
    if !condition.Page.Disable {
        limit, offset := condition.Page.Limit()
        db.Limit(int(limit)).Offset(int(offset))
    }

    // Execute query
    list, err := db.Find(ctx)
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("find game failed")
        return
    }

    copierx.Copy(&rp, list)
    return
}
```

### Update (with diff detection)

```go
func (ro gameRepo) Update(ctx context.Context, item *biz.UpdateGame) (err error) {
    db := gorm.G[model.Game](ro.data.DB(ctx))

    // Get existing record
    m, err := db.Where("id = ?", item.ID).First(ctx)
    if err == gorm.ErrRecordNotFound {
        return biz.ErrRecordNotFound(ctx)
    }
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("get game failed")
        return
    }

    // Detect changes
    change := make(map[string]interface{})
    utils.CompareDiff(m, item, &change)
    if len(change) == 0 {
        return biz.ErrDataNotChange(ctx)
    }

    // Check uniqueness if unique field is being updated
    if item.Name != nil && (m.Name == nil || *item.Name != *m.Name) {
        count, err := db.Where("name = ? AND id != ?", *item.Name, item.ID).Count(ctx, "*")
        if err != nil {
            log.WithContext(ctx).WithError(err).Error("check name uniqueness failed")
            return err
        }
        if count > 0 {
            return biz.ErrDuplicateField(ctx, "name", *item.Name)
        }
    }

    // IMPORTANT: Use native DB.Updates for map updates (gorm.G.Updates expects struct type)
    err = ro.data.DB(ctx).Model(&model.Game{}).Where("id = ?", item.ID).Updates(change).Error
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("update game failed")
    }
    return
}
```

### Delete

```go
func (ro gameRepo) Delete(ctx context.Context, ids ...int64) (err error) {
    db := gorm.G[model.Game](ro.data.DB(ctx))

    _, err = db.Where("id IN ?", ids).Delete(ctx)
    if err != nil {
        log.WithContext(ctx).WithError(err).Error("delete game failed")
    }
    return
}
```

**Key Points:**
- Always use `gorm.G[model.XXX](ro.data.DB(ctx))` to create typed query builder
- Use `copierx.Copy` for model ↔ biz struct conversion
- Use `utils.CompareDiff` to detect changes in Update
- For Update with map: use **native `ro.data.DB(ctx).Model().Where().Updates(change)`** (not `gorm.G.Updates`)
- Check `gorm.ErrRecordNotFound` for Get/Update operations
- Apply pagination with `condition.Page.Limit()` and early return when count is 0

## 9. Repo Interface Location

All `Repo` interfaces in the `biz` layer **must** be defined in `biz.go`, not scattered across individual entity files.

```go
// Good - all Repo interfaces in biz/biz.go
// biz/biz.go
package biz

type UserRepo interface {
    Create(ctx context.Context, item *User) error
    Find(ctx context.Context, condition *FindUser) []User
    // ...
}

type RoleRepo interface {
    Create(ctx context.Context, item *Role) error
    Find(ctx context.Context, condition *FindRole) []Role
    // ...
}

type ActionRepo interface {
    Create(ctx context.Context, item *Action) error
    Find(ctx context.Context, p *page.Page, filter *Action) ([]*Action, int64, error)
    // ...
}

// Bad - Repo interface in individual entity file
// biz/user.go
package biz

type UserRepo interface {  // Should be in biz.go
    Create(ctx context.Context, item *User) error
    // ...
}
```

**Benefits:**
- Centralized dependency interface definitions
- Easier to review all data layer contracts in one place
- Cleaner separation: entity files contain only structs and use cases

## 10. Proto Field Naming Convention

In `api/*.proto` files, use **snake_case** for field names, not camelCase. Protobuf's standard convention uses snake_case.

```protobuf
// Good - snake_case field names
message CreateUserRequest {
    string user_name = 1;
    string role_id = 2;
    string created_at = 3;
    bool is_locked = 4;
}

message FindUserRequest {
    params.Page page = 1;
    string start_created_at = 2;
    string end_created_at = 3;
}

// Bad - camelCase field names
message CreateUserRequest {
    string userName = 1;    // Should be user_name
    string roleId = 2;      // Should be role_id
    string createdAt = 3;   // Should be created_at
    bool isLocked = 4;      // Should be is_locked
}
```

**Note:** The generated Go code will automatically convert snake_case to CamelCase (e.g., `user_name` → `UserName`).

## 11. Config YAML Field Naming Convention

In `configs/*.yaml` files, use **camelCase** for field names, not snake_case. This matches the Go struct field naming convention.

```yaml
# Good - camelCase field names
server:
  http:
    addr: 0.0.0.0:8080
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9090
    timeout: 1s
  jwt:
    key: your-secret-key
    expires: 24h

data:
  database:
    driver: mysql
    dsn: root:password@tcp(127.0.0.1:3306)/dbname
    maxIdleConns: 10
    maxOpenConns: 100
    connMaxLifetime: 1h

# Bad - snake_case field names
data:
  database:
    max_idle_conns: 10    # Should be maxIdleConns
    max_open_conns: 100   # Should be maxOpenConns
    conn_max_lifetime: 1h # Should be connMaxLifetime
```

**Note:** YAML field names should match the Go struct tags (typically `json` or `yaml` tags with camelCase).

## 12. Service Layer Unit Tests

Every public method in the `service` layer **must** have at least one corresponding unit test. Tests are located in `internal/tests/service/` directory.

### Test File Naming

- Service file: `internal/service/game.go`
- Test file: `internal/tests/service/game_test.go`

### Test Method Naming

Test method names must follow the pattern `Test<MethodName>`:

| Service Method | Test Method |
|----------------|-------------|
| `CreateGame` | `TestCreateGame` |
| `FindGame` | `TestFindGame` |
| `UpdateGame` | `TestUpdateGame` |
| `DeleteGame` | `TestDeleteGame` |

### Test Implementation Pattern

All service tests **must** use `mock.XxxService()` to create the service instance:

```go
// Good - use mock.XxxService() to create service instance
func TestFindGame(t *testing.T) {
    s := mock.GameService()
    ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

    req := &v1.FindGameRequest{
        Page: &params.Page{
            Num:  1,
            Size: 10,
        },
    }

    rp, err := s.FindGame(ctx, req)
    assert.NoError(t, err)
    assert.NotNil(t, rp)
    assert.NotNil(t, rp.Page)
}

func TestCreateGame(t *testing.T) {
    s := mock.GameService()
    ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

    name := "test-game"
    req := &v1.CreateGameRequest{
        Name: &name,
    }

    _, err := s.CreateGame(ctx, req)
    assert.NoError(t, err)
}

// Bad - creating service instance manually
func TestFindGame(t *testing.T) {
    // Don't do this - always use mock.XxxService()
    s := &service.GameService{}
    // ...
}
```

### Test File Structure

```go
package service

import (
    "context"
    "testing"

    "github.com/go-cinch/common/page/params"
    v1 "auth/api/auth"
    "auth/internal/tests/mock"
    "github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {
    s := mock.GameService()
    ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")
    // test implementation...
}

func TestFindGame(t *testing.T) {
    s := mock.GameService()
    ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")
    // test implementation...
}

func TestUpdateGame(t *testing.T) {
    s := mock.GameService()
    ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")
    // test implementation...
}

func TestDeleteGame(t *testing.T) {
    s := mock.GameService()
    ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")
    // test implementation...
}
```

**Key Points:**
- Every service method must have a corresponding `Test<MethodName>` function
- Always use `mock.XxxService()` to create the service instance
- Always use `mock.NewContextWithUserId()` to create context with tenant ID
- Use `github.com/stretchr/testify/assert` for assertions
- Test files are in `internal/tests/service/`, not alongside service files

## 13. Database Naming Convention

Table names **must** use `t_` prefix with singular form. Index names **do not** use `t_` prefix. Column names **must not** use database reserved keywords.

```sql
-- Table naming
CREATE TABLE t_order (...)    -- Good
CREATE TABLE t_user (...)     -- Good
CREATE TABLE order (...)      -- Bad: reserved keyword
CREATE TABLE orders (...)     -- Bad: plural form

-- Index naming (no t_ prefix)
CREATE INDEX idx_user_name ON t_user(name);           -- Good
CREATE UNIQUE INDEX idx_user_word ON t_user(word);    -- Good
CREATE INDEX idx_t_user_name ON t_user(name);         -- Bad: unnecessary t_ prefix

-- Column naming
order_type VARCHAR(50)        -- Good
user_group VARCHAR(100)       -- Good
type VARCHAR(50)              -- Bad: reserved keyword
group VARCHAR(100)            -- Bad: reserved keyword
```

**Common reserved keywords to avoid in columns:**
- `order`, `group`, `select`, `table`, `index`, `key`, `type`, `status`, `value`
