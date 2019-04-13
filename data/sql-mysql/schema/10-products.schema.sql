CREATE TABLE `products` (
  `id` INT(10),
  `legal_owner` INT(10) NOT NULL,
  `display_name` VARCHAR(128) NOT NULL,
  `summary` VARCHAR(512) NOT NULL,
-- see ../docs/Relational-Schemas.md#reformatting-data-via-a-trigger
  `support_phone` VARCHAR(12),
  `support_email` VARCHAR(255) NOT NULL,
  `homepage` VARCHAR(255),
  `logo_url` VARCHAR(255),
  `repo_url` VARCHAR(255),
  `issues_url` VARCHAR(255),
  `ontology` ENUM ('TANGIBLE GOOD', 'DIGITAL GOOD', 'SOFTWARE SERVICE', 'CONSULTING SERVICE', 'PHYSICAL SERVICE'),

  CONSTRAINT `products_key` PRIMARY KEY ( `id` ),
  CONSTRAINT `products_ref_entities` FOREIGN KEY ( `id` ) REFERENCES `entities` ( `id` ),
  CONSTRAINT `products_ref_users` FOREIGN KEY ( `legal_owner` ) REFERENCES `users` ( `id` )
);
DELIMITER //
CREATE TRIGGER `products_phone_format`
  BEFORE INSERT ON products FOR EACH ROW
    BEGIN
      SET new.support_phone=(SELECT NUMERIC_ONLY(new.support_phone));
    END;//
DELIMITER ;
