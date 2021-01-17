CREATE TABLE IF NOT EXISTS offers (
    offer_id bigint,
    seller_id bigint,
    price bigint,
    name text,
    quantity bigint,
    UNIQUE(offer_id, seller_id)
);
