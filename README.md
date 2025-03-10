# Gopher loyalty project

# Docs

Планы как сделаю проект:
- Использую DDD подход для организации структуры проекта и кода.
- Выделю Worker как отдельное приложение (если успею, иначе будет внутри API жить), которое будет работать с единой БД Postgres.
- Worker будет заниматься опросом AccrualAPI для получения резульата обработки.
- Для аутентификации применю JWT.
- Для генерации http api применю подход через написание спеки openapi v3 + кодогенерацию + пакет go-chi.
- Для поддержания актуальной документации воспользуюсь пакетом swaggest.
- Работа с БД через database/sql.

Схема модулей проекта

```mermaid
flowchart TD
    GophermartBuyer --> JWTAuth
    JWTAuth --> API

    subgraph GopherLoyalty[Группа GopherLoyalty]
        API[API]
        Worker[Worker]
        Postgres[Postgres]

        API --> Postgres
        Worker --> Postgres
    end

    subgraph External[Внешняя система]
        AccrualAPI[AccrualAPI]
    end

    Worker -- Периодический опрос --> AccrualAPI

    GophermartBuyer[Gophermart Buyer]
    JWTAuth[JWT Auth]
```

Схема БД и сущностей DDD (1 к 1)

```mermaid
erDiagram
    User {
        uuid UserID
        string login
        string password
    }

    LoyaltyAccount {
        uuid LoyaltyAccountID
        int UserID "unique"
        int balance
        int withdrawTotal
    }

    LoyaltyAccountTransaction {
        uuid LoyaltyAccountTransactionID
        datetime CreatedAt
        uuid LoyaltyAccountID
        enum TransactionType
        uuid OrderID
    }

    Order {
        uuid OrderID
        string Number
        datetime UploadedAt
        int UserID
        enum Status
    }

    BonusCalculation {
        uuid BonusCalculationID
        uuid OrderID "unique"
        enum LoyaltyStatus
        int Accrual
    }

    User ||--o{ Order : "Order.UserID -> User.UserID"
    User ||--|| LoyaltyAccount : "LoyaltyAccount.UserID -> User.UserID"
    LoyaltyAccount ||--o{ LoyaltyAccountTransaction : "ID -> LoyaltyAccountID"
    Order ||--o{ LoyaltyAccountTransaction : "ID -> OrderID"
    Order ||--|| BonusCalculation : "ID -> OrderID"
```