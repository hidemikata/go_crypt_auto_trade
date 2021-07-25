
#config/config.ini
  
[bitflyer]
  
api_key = 
  
secret_key = 
    
#mysql
  
CREATE TABLE `btc_jpy_live` (
  `date` datetime DEFAULT NULL,
  `symbol` varchar(10) DEFAULT NULL,
  `open` float DEFAULT NULL,
  `high` float DEFAULT NULL,
  `low` float DEFAULT NULL,
  `close` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
  
CREATE TABLE `btc_jpy_live_position` (
  `date` datetime DEFAULT NULL,
  `buy_or_sell` varchar(10) DEFAULT NULL,
  `price` float DEFAULT NULL,
  `fix_date` datetime DEFAULT NULL,
  `fix_price` float DEFAULT NULL,
  `profit` float DEFAULT NULL,
  `symbol` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

  

go run main.go


