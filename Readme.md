# Пример Select с указанием полей

> Если нужно сделать множественный GroupBy или OrderBy, то используй **MultiGroupBy** и **Orders**

```go
builder := Select{
    From: "users",
    As:   "u",
    Fields: []any{
        // будет преобразовано в {As}.id, если As указан, иначе id
        "id",   
        
        // будет преобразовано в {As}.name, если As указан, иначе name
        "name", 
        
        // Если нужно сделать выборку поля по каким-то правилам или используя функцию, то можно использовать RawField
        RawField{Sql: "IF(u.balance > 0, TRUE, FALSE) as has_balance"}, // вставляет как есть
        RawField{Sql: "NOW() as current_date"},                         // вставляет как есть
        
        // если указан TableAlias, то будет преобразовано в {TableAlias}.name, иначе если указан As у таблицы в Select, то в {As}.name, иначе в name
        // если указан As у Field, то к полю будет добавлено As {As}
        Field{
            Name:       "name",
            TableAlias: "f",
            As:         "friend_name",
        }, 
        
        // если указан TableAlias, то будет преобразовано в {TableAlias}.last_name, иначе если указан As у таблицы в Select, то в {As}.last_name, иначе в last_name
        // если указан As у Field, то к полю будет добавлено As {As}
        Field{
            Name:       "last_name",
            TableAlias: "f",
            As:         "friend_last_name",
        }, 
    },
    Joins: []Join{
        {
            Table: "friends",
            As:    "f",
            On:    &On{Fields: []string{"u.friend_id", "f.id"}},
        },
    },
    Where:   (&Condition{Field: "f.id", Operator: MoreThen, Arg: 100}).
			And(Condition{Field: "f.ass", Operator: Equal, Arg: "penes"}),
    GroupBy: "u.id",
    Order: Order{
        Direction: DESC,
        Field:     "f.id",
    },
    Limit:  10,
    Offset: 200,
}

query, args := builder.ToSql()

fmt.Println(query, args)
```

# Пример Select с автоматическим парсингом полей из структуры

> Каждое поле, которое должно попасть в запрос должно иметь тег db c название поля в таблице.<br>
> Если мы хотим указать тег db, но не хотим, чтоб поле спарсилось, то указываем тег sb:"skip".<br>
> Таким образом можно комбинировать явно указанные поля и автоматически спарсенные поля.
```go
type UserModel struct {
	Id          uint           `db:"id"`
	Uuid        string         `db:"uuid"`
	CurrencyId  sql.NullInt64  `db:"currency_id"`
	CurrentDate time.Time      `db:"current_date" sb:"skip"`
}


builder := Select{
    From: "users",
    As:   "u",
    StructFields: UserModel{},
    Fields: []any{
        RawField{Sql: "NOW() as current_date"},
    }
}

query, args := builder.ToSql()

fmt.Println(query, args)
```

# Пример Update

```go
builder := Update{
    Table:  "users",
    Values: []SetValue{
        {
            Field: "name",
            Value: "SamOgon",
        },
        {
            Field: "balance",
            Value: Increment{ Value: 134 }, // если нужно сделать balance = balance + ?
        },
    },
}

query, args := builder.ToSql()

fmt.Println(query, args)
```

# Пример Insert с указанием полей

```go
builder := Insert{
    Table:        "users",
    Values:       []Value{
        {
            Field: "name",
            Arg:   "SamOgon",
        },
        {
            Field: "last_seen",
            Fn:    NOW, // если нужно использовать функцию
        },
    },
    Timestamps:   true, // если нужно автоматически заполнить поля created_at и updated_at
}

query, args := builder.ToSql()

fmt.Println(query, args)
```


# Пример Insert с автоматическим парсингом полей из структуры
> Каждое поле, которое должно попасть в запрос должно иметь тег db c название поля в таблице.
```go
type CreateUserDto struct {
    Name     string `db:"name"`
    LastName string `db:"last_name"`
    Age      uint   `db:"age"`
}

dto := CreateUserDto{
    Name:     "Sam",
    LastName: "Ogon",
    Age:      69,
}

builder := Insert{
    Table:        "users",
    StructValues: dto,
}

query, args := builder.ToSql()

fmt.Println(query, args)
```

# Пример Delete
```go
builder := Delete{
    From:  "users",
    Where: &Condition{Field: "name", Operator: semen_builder.Equal, Arg: "SamOgon"},
}

query, args := builder.ToSql()

fmt.Println(query, args)
```

# Пример In
> Чтоб передать несколько значений в условие, используй ***Args***
```go
builder := Select{
    From: "users",
    As:   "u",
    Fields: []any{
        "id",   
        "name",
    },
    Where:   &Condition{Field: "f.id", Operator: In, Args: []any{1,2,3,4,5}},
}

query, args := builder.ToSql()

fmt.Println(query, args)
```

# Если нужно добавить Condition в пустой Where по условию
> Если необходимо добавить ***Condition*** в пустой ***Where*** по условию, то можно использовать функции ***AppendOr*** и ***AppendAnd***.
> Если ***Where*** не пустой, то будет добавлено ***Condition*** с and/or, иначе ***Condition*** стананет корневым.

```go
searchName := ""
searchId := 228

builder := Select{
    From: "users",
    As:   "u",
    Fields: []any{
        "id",
        "name",
    },
}

if searchName != "" {
    builder.Where = &Condition{Field: "name", Operator: semen_builder.Equal, Arg: searchName}
}

if searchId != 0 {
    builder.Where = AppendOr(builder.Where, Condition{Field: "id", Operator: Equal, Arg: searchId})
}

query, args := builder.ToSql()

fmt.Println(query, args)
```