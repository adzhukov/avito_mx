CREATE TABLE IF NOT EXISTS offers (
    offer_id int,
    seller_id int,
    price int,
    name text,
    quantity int,
    UNIQUE(offer_id, seller_id)
);
