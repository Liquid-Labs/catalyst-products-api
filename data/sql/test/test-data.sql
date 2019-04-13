-- the legal owner
INSERT INTO entities (pub_id) VALUES ('4C2B3954-8D7F-48BA-B720-3B0F15F91BA9');
-- TODO: 'SET' is not ANSI SQL; for this and other reasons, we want to do a
-- replacement scheme. Possibly something like:
-- 1) Name template files with a commen prefix ('.sql.template').
-- 2) Use bash subsitutios, so "VALUES ($JANE_DOE_ID)"
-- 3) Have a 'template.vars' file.
-- 4) source template.vars; for $TEMPLATE in ...; do ...; eval "$(cat "$TEMPLATE")" > $SQL_FILE; done
SET @some_org_id=LAST_INSERT_ID();
INSERT INTO users (id, auth_id, legal_id, legal_id_type, active) VALUES (@some_org_id,'xzy098', '55-5555555', 'EIN', 0);
-- INSERT INTO orgs (id, display_name, summary, phone, email) VALUES (@some_org_id,'Some Org','Builders of things.','5555551111','janedoe@test.com');

INSERT INTO entities (pub_id) VALUES ('D929BEE3-8034-40A9-B33E-E1A28507EE68');
SET @proudct_a=LAST_INSERT_ID();
INSERT INTO products (id, legal_owner, display_name, summary, support_email, homepage, logo_url, repo_url, ontology)
  VALUES (@proudct_a, @some_org_id, 'Bauble', 'A thing for your wall.', 'bauble@foo.com', 'https://foo.com/proudcts/bauble', 'https://foo.com/assets/bauble_logo.svg', 'https://git.foo.com/bauble_repo', 'TANGIBLE GOOD');

INSERT INTO entities (pub_id) VALUES ('016B5F34-D36A-4970-ADC8-4FADC01425D9');
SET @proudct_b=LAST_INSERT_ID();
INSERT INTO products (id, legal_owner, display_name, summary, support_email, homepage, logo_url, repo_url, ontology)
  VALUES (@proudct_b, @some_org_id, 'Blog', 'Online articles.', 'blog@foo.com', 'https://foo.com/sass/blog', 'https://foo.com/assets/blog_logo.svg', 'https://git.foo.com/blog_repo', 'SOFTWARE SERVICE');
