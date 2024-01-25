--Entities:
--customer, payment, order, kitchenOrder, orderItem, table, staff

CREATE TABLE customer(
    customerID      varchar(5)      PRIMARY KEY,
    customerName    varchar(20),
    instructions    varchar(100),
    allergies       varchar(50)
);

CREATE TABLE staff(
    staffID     varchar(10) PRIMARY KEY,
    staffName   varchar(20),
    staffRole   varchar(15),
    contact     varchar(10)
);

CREATE TABLE payment(
    paymentID   varchar(10)     PRIMARY KEY,
    price       numeric(10,2),
    orderID     varchar(10),
    orderTime   varchar(4)
);

CREATE TABLE menuItem(
    menuItemID  varchar(10),
    itemName    varchar(20),
    price       numeric(10,2),
    calories    float
);

CREATE TABLE restauranttable(
    tableID         varchar(10),
    tableLocation   varchar(10)
);
