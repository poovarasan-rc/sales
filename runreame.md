
## Setup and Running the Project

## ✅ Prerequisites

Before running the project, ensure you have the following installed:

### 1. Go
- **Version**: 1.18 or higher  
- **Download**: https://golang.org/dl/

### 2. MySQL
- **Version**: 5.7 or later  
- **Download**: https://dev.mysql.com/downloads/

### 3. Clone the repository
```bash
git clone git@github.com:poovarasan-rc/sales.git
cd sales
```

### 4. API endpoinds

---

## 4.1 📡 URL: `http://localhost:8080/loadcsv`

### Description
This endpoint loads CSV data into the database by internally calling the `ReadConsStore()` function.

### ⚠️ Requirements
- No request body is needed.
- Only `POST` method is supported.
- CORS is enabled (allows all origins).

### ✅ Success Response

- **Status Code:** `200 OK`
- **Body:** `Data stored successfully ...`

### ❌ Invalid Method

- **Status Code:** `405 Method Not Allowed`
- **Body:** `Invalid method. Use POST.`

### ❌ Error Response

- **Status Code:** `501 Not Implemented`
- **Body:** `Error - <error message>`

---

## 4.2 📡 URL: `http://localhost:8080/getrevenue`

### Description
This endpoint fetch the revenue details in category, product, region and total revenue for the respective date range.

### ⚠️ Requirements
- Request body is needed.
- Only `PUT` method is supported.
- CORS is enabled (allows all origins).

### Total Revenue: (For a date range)

```bash
Request :
{
    "fdt":"2024-01-03",
    "tdt":"2024-05-23"
}


### ✅ Success Response :
{
    "revn":4448.57,
    "sts":"S",
    "emsg":""
}
```

### Total Revenue by Product:(For a date range)

```bash
Request :
{
    "type":"product",
    "fdt":"2024-01-03",
    "tdt":"2024-05-23"
}


### ✅ Success Response :
{
    "revndata": [
        {
            "name": "iPhone 15 Pro",
            "revn": 3802.1
        },
        {
            "name": "Levi's 501 Jeans",
            "revn": 148.98
        },
        {
            "name": "UltraBoost Running Shoes",
            "revn": 188
        },
        {
            "name": "Sony WH-1000XM5 Headphones",
            "revn": 309.49
        }
    ],
    "sts": "S",
    "emsg": ""
}
```

### Total Revenue by Category:(For a date range)

```bash
Request :
{
    "type":"category",
    "fdt":"2024-01-03",
    "tdt":"2024-05-23"
}


### ✅ Success Response :
{
    "revndata": [
        {
            "name": "Electronics",
            "revn": 4111.59
        },
        {
            "name": "Clothing",
            "revn": 148.98
        },
        {
            "name": "Shoes",
            "revn": 188
        }
    ],
    "sts": "S",
    "emsg": ""
}
```

### Total Revenue by Region:(For a date range)

```bash
Request :
{
    "type":"region",
    "fdt":"2024-01-03",
    "tdt":"2024-05-23"
}


### ✅ Success Response :
{
    "revndata": [
        {
            "name": "Europe",
            "revn": 1314
        },
        {
            "name": "Asia",
            "revn": 2637.08
        },
        {
            "name": "South America",
            "revn": 188
        },
        {
            "name": "North America",
            "revn": 309.49
        }
    ],
    "sts": "S",
    "emsg": ""
}
```

### ❌ Error Response

```bash
{
  "sts": "E",
  "emsg": "RGRD01 : <error message>"
}
```

---

### 🗄️ 5. Database Creation Script (MySQL)

The following SQL script sets up the necessary tables for the application. Ensure that your MySQL server is running and you have access to a database where these tables can be created.

> ⚠️ Note: Replace `createdby` and `createddate` values in actual inserts as per your application logic.

#### 📋 SQL Script

```sql
-- Table: product_details
CREATE TABLE `product_details` (
  `prd_id` varchar(100) NOT NULL,
  `name` varchar(100) NOT NULL,
  `unitprice` decimal(13,2) NOT NULL,
  `createdby` varchar(100) NOT NULL,
  `createddate` datetime NOT NULL,
  PRIMARY KEY (`prd_id`)
);

-- Table: customer_details
CREATE TABLE `customer_details` (
  `cus_id` varchar(100) NOT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `address` varchar(300) NOT NULL,
  `createdby` varchar(100) NOT NULL,
  `createddate` datetime NOT NULL,
  PRIMARY KEY (`cus_id`),
  UNIQUE KEY `email` (`email`)
);

-- Table: sales_details
CREATE TABLE `sales_details` (
  `order_id` int NOT NULL,
  `product_id` varchar(100) NOT NULL,
  `customer_id` varchar(100) NOT NULL,
  `category` varchar(100) DEFAULT NULL,
  `region` varchar(100) DEFAULT NULL,
  `date_of_sale` date DEFAULT NULL,
  `quantity_sold` int DEFAULT NULL,
  `discount` decimal(5,2) DEFAULT NULL,
  `shipping_cost` decimal(10,2) DEFAULT NULL,
  `payment_method` varchar(50) DEFAULT NULL,
  `createdby` varchar(100) NOT NULL,
  `createddate` datetime NOT NULL,
  PRIMARY KEY (`order_id`),
  KEY `sales_details_customer_details_FK` (`customer_id`),
  KEY `sales_details_product_details_FK` (`product_id`),
  KEY `idx_date_of_sale` (`date_of_sale`),
  KEY `idx_category` (`category`),
  KEY `idx_region` (`region`),
  CONSTRAINT `sales_details_customer_details_FK` FOREIGN KEY (`customer_id`) REFERENCES `customer_details` (`cus_id`),
  CONSTRAINT `sales_details_product_details_FK` FOREIGN KEY (`product_id`) REFERENCES `product_details` (`prd_id`)
);
