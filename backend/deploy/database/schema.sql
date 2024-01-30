create table menuitem (
    itemid serial not null primary key,
    itemname varchar(20) not null unique,
    price numeric(10, 2) not null,
    calories int
);