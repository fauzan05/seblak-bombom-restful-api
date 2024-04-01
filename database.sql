create database seblak_bombom;
create database seblak_bombom_test;

use seblak_bombom;
use seblak_bombom_test;

select * from users;
select count(*) from users where first_name != "Rudi";
select * from addresses;
select * from tokens;
select * from products;
delete from tokens where id = 1;
show create table users;
show create table addresses;

drop table users;
drop table addresses;