# インスタンス生成する
sudo apt update
wget https://golang.org/dl/go1.16.7.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.16.7.linux-amd64.tar.gz
vim .bashrc
   export PATH=$PATH:/usr/local/go/bin
sudo apt install mysql-server mysql-client
sudo systemctl enable mysql
sudo mysql -uroot -p
mysql> ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '';
mysql> FLUSH PRIVILEGES;

git clone https://github.com/hidemikata/go_crypt_auto_trade.git
mv go_crypt_auto_trade btcanallive_refact
cd btcanallive_refact
go get


# config/config.ini
  
[bitflyer]
  
api_key = 
  
secret_key = 
    
# mysql
creata database coin_data;
  
  
CREATE TABLE `btc_jpy_live` (
  `date` datetime DEFAULT NULL,
  `symbol` varchar(10) DEFAULT NULL,
  `open` float DEFAULT NULL,
  `high` float DEFAULT NULL,
  `low` float DEFAULT NULL,
  `close` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
  
CREATE TABLE `btc_jpy_live_position` (
  `date` datetime DEFAULT NULL,
  `buy_or_sell` varchar(10) DEFAULT NULL,
  `price` float DEFAULT NULL,
  `fix_date` datetime DEFAULT NULL,
  `fix_price` float DEFAULT NULL,
  `profit` float DEFAULT NULL,
  `symbol` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `backtest_profit` (
  `date` datetime DEFAULT NULL,
  `total_profit` float DEFAULT NULL,
  `sma_long` int DEFAULT NULL,
  `sma_short` int DEFAULT NULL,
  `sma_min_max_rate` float DEFAULT NULL,
  `rci` int DEFAULT NULL,
  `position_count` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
  
  CREATE TABLE `btc_jpy_live_position_backtest` (
  `date` datetime DEFAULT NULL,
  `buy_or_sell` varchar(10) DEFAULT NULL,
  `price` float DEFAULT NULL,
  `fix_date` datetime DEFAULT NULL,
  `fix_price` float DEFAULT NULL,
  `profit` float DEFAULT NULL,
  `symbol` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#ローソク足データインポート
mysql -uroot -p -D coin_data < btc_jpy_live.sql

# mysqlcnf
/etc/mysql/mysql.conf.d/mysqld.cnf

key_buffer_size         = 1024M
max_allowed_packet      = 16M
thread_stack            = 1M
thread_cache_size       = 8

query_cache_limit       = 2048M
query_cache_size        = 2048M
query_cache_type=1
innodb_buffer_pool_size=1024M
sort_buffer_size=2048M
read_buffer_size=2048M
innodb_thread_concurrency=128

# run
go run main.go


