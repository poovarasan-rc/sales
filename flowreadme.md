
# 📊 Revenue Reporting Go Service

This Go web server provides endpoints for loading CSV data and retrieving revenue reports, with periodic data refresh and configurable settings via a TOML file.

---

## 🚀 Features

- ✅ Connects to a database on startup.
- 🔄 Periodically refreshes data in the background.
- 📁 Loads data from a CSV via an HTTP endpoint.
- 📊 Serves revenue report data via a separate endpoint.
- ⚙️ Configuration handled through a TOML config file.

---

## 📁 Project Structure

```
main.go              # Entry point
sales/
├── db/              # Database connection logic
├── loadcsvdata/     # Handles loading CSV data
├── refreshdata/     # Periodic refresh implementation
├── report/          # Revenue report logic
└── toml/            # Configuration values in TOML format
```

---

## 🛠️ Installation

```bash
git clone <your-repo>
cd <your-repo>
go mod tidy
go run main.go
```

---

## ⚙️ Configuration

The `toml/` directory contains configuration files (e.g., `config.toml`) used to control values like:

- Database connection strings
- Refresh interval durations
- File paths for CSV imports


---

## 🔌 API Endpoints

### `POST /loadcsv`

This API reads a structured CSV file (./sample.csv), parses the data, checks for integrity, and performs bulk inserts or updates into the database for Products, Customers, and Orders.

---

### 🧭 Flow Overview

1. **HTTP Handler (`LoadCSVData`)**
   - Validates request method (must be `POST`)
   - Calls `ReadConsStore()` to process the data
   - Handles CORS headers
   - Returns:
     - `200 OK` on success
     - `405` on invalid method
     - `501` on internal processing errors

2. **Data Orchestration (`ReadConsStore`)**
   - Prevents concurrent ingestion using `GFlag`
   - Calls `Readfile()` to parse the CSV
   - Transforms data using `ConstructRec()`
   - Passes result to `StoreData()` for DB persistence

3. **CSV File Reader (`Readfile`)**
   - Reads all rows using `encoding/csv`
   - Returns `[][]string` or logs file read errors

4. **Data Construction (`ConstructRec`)**
   - Detects duplicate product/customer names by ID
   - Updates or overrides duplicate order records
   - Aggregates SQL value strings into an `InsStruct`

5. **Product & Customer:** Conflicting names for existing IDs result in a failure.
   **Orders:**
   - If duplicated in the file → the latest is used.
   - If already in DB → the old order is marked for replacement.

6. **Struct Output**
   ```go
   type InsStruct struct {
       PrdIns string // Product INSERT values
       CusIns string // Customer INSERT values
       OrdMod string // Order IDs to modify
       OrdIns string // Order INSERT values
   }
   ```

7. **🧩 SQL Execution Steps**

    7.1. **Begin Transaction**
   - Opens a transaction using `db.GDBCon.Begin()`

    7.2. **Insert Products**
   ```sql
   INSERT INTO product_details (prd_id, name, unitprice, createdby, createddate) VALUES ...
   ```

    7.3. **Insert Customers**
   ```sql
   INSERT INTO customer_details (cus_id, name, email, address, createdby, createddate) VALUES ...
   ```

    7.4. **Delete Orders**
   - Deletes old records when modifying existing orders.
   ```sql
   DELETE FROM sales_details WHERE order_id IN (...)
   ```

    7.5. **Insert Orders**
   ```sql
   INSERT INTO sales_details (
       order_id, product_id, customer_id, category, region, date_of_sale,
       quantity_sold, discount, shipping_cost, payment_method, createdby, createddate
   ) VALUES ...
   ```

    7.6. **Commit Transaction**
   - If all queries succeed, the transaction is committed.
   - If any step fails, a rollback is triggered and the error is logged.

### 🚫 Error Handling

- Every DB operation is checked.
- On failure:
  - Logs the specific step using error codes (`RSD01` to `RSD04`)
  - Rolls back the entire transaction to maintain integrity
- Returns the error back to the caller for HTTP response generation.

### 🛡️ Atomic Operation Guarantee

All inserts and deletes occur in a **single database transaction**. If any insert or delete fails, no data is committed — ensuring all-or-nothing behavior.

---

- Returns:
  - `200 OK` on success
  - `405 Method Not Allowed` for non-POST requests
  - `501 Not Implemented` on processing error

---

### `GET /getrevenue`

Returns revenue details from the processed data.

---

## 🔄 Periodic Refresh

A background goroutine (`refreshdata.PeriodicRefresh`) runs at regular intervals to keep data fresh. The refresh interval is configurable via the TOML file.

---
