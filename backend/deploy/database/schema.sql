CREATE TABLE menuitem(
    menuitemid serial PRIMARY KEY,
    calories int,
    menuitemname varchar(20) unique,
    itemdescription varchar(50),
    price float
);

CREATE TABLE allergens(
    menuitemid  int,
    allergen    varchar(20),
    PRIMARY KEY(menuitemid,allergen),
    FOREIGN KEY (menuitemid) REFERENCES menuitem(menuitemid)
);

CREATE TABLE restaurantorders(
    orderID         int     PRIMARY KEY,
    orderDateAndTime timestamp,
    tableNumber     int,
    amountNumber    float,
    paid            int
);

CREATE TABLE ordermenuitems(
    orderID     int,
    menuItemID  int,
    orderStatus varchar(20),
    PRIMARY KEY(orderID,menuItemID),
    FOREIGN KEY (orderID) REFERENCES restaurantorders(orderID),
    FOREIGN KEY (menuItemID) REFERENCES menuItem(menuItemID)
);

CREATE TABLE staff(
    staffID     int  PRIMARY KEY,
    staffName   varchar(15)
);
