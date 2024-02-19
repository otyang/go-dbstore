# go-dbstore

## Overview

The `dbstore` package offers a streamlined approach to performing database operations in Go applications. It introduces a generic repository pattern, `IRepository`, to efficiently manage data persistence using the bun: [https://github.com/uptrace/bun](https://github.com/uptrace/bun) database library. This pattern simplifies your codebase and eliminates repetitive CRUD implementation.

## Features

* **Versatile Interface:** The `IRepository` interface exposes methods for fundamental database operations:
    * `Create`
    * `CreateBulk`
    * `FindOneByPK`
    * `FindOneWhere`
    * `FindManyWhere`
    * `UpdateOneByPK`
    * `UpdateManyByPK`
    * `UpdateOneWhere`
    * `Upsert`
    * `DeleteByPK`
    * `DeleteWhere`
* **Criteria-Based Operations:**  `SelectCriteria`, `UpdateCriteria`, and `DeleteCriteria` functions for flexible queries.
* **Transaction Support:**
    * `NewWithTx` to inject existing transactions
    * `Transaction` to simplify transaction management and error handling 

 
## Usage Examples

**1. Creating a Repository:**

```go
import (
  "github.com/uptrace/bun"
  "github.com/otyang/dbstore"
  "github.com/otyang/dbstore/obun"
)

type User struct {
  ID       int64  `bun:",pk,autoincrement"`
  Name     string
  Email    string
}

var db *bun.DB // Assumes 'db' is your initialized bun database connection

userRepository := obun.NewRepository(db, (*User)(nil))
```

**2. Inserting a Record:**

```go
user := &User{Name: "Alice", Email: "alice@example.com"}
err = userRepository.Create(ctx, user, false) // Set suppressDuplicate as needed
if err != nil {
    // Handle error
}
```

**3. Finding Records with Criteria:**

```go
var users []*User
query := obun.FindManyWhere(ctx, &users, dbstore.PaginationOption{},
  func(q *bun.SelectQuery) *bun.SelectQuery {
     return q.Where("name LIKE ?", "%Bob%")
  })
if query.Error() != nil {
    // Handle error
}
```

**4. Transactional Operations:** 

```go
err := obun.Transaction(ctx, func(ctx context.Context, tx bun.Tx) error {
   // ...  perform multiple operations within the transaction
   return nil // Commit on success, rollback automatically on error
})
if err != nil {
    // Handle error
}
```
 

## License

The `dbstore` package is released under the MIT License: [https://choosealicense.com/licenses/mit/](https://choosealicense.com/licenses/mit/). 
