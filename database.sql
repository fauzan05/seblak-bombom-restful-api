create database seblak_bombom;
create database seblak_bombom_test;
show tables;

use seblak_bombom;
use seblak_bombom_test;

select * from users;
select count(*) from users where first_name != "Rudi";
select * from addresses;
select * from tokens;
select * from products;
select * from images;
select * from categories;
select * from orders;
select * from order_products;
select * from midtrans_snap_orders;

show create table users;
show create table addresses;

delete from tokens where id = 1;


drop table users;
drop table addresses;