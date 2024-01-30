create table menuitem (
    menuitemid serial not null primary key,
    itemname varchar(20) not null,
    price numeric(10, 2) not null,
    calories int
);