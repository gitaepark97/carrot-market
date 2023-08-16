DROP TABLE IF EXISTS goods_images;

DROP TABLE IF EXISTS goods_categories;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS goods;

ALTER TABLE sessions ALTER COLUMN user_id SET DEFAULT nextval('sessions_user_id_seq'::regclass);
