INSERT INTO apps (id, name, secret)
VALUES (2, 'test', 'test_secret')
ON CONFLICT DO NOTHING;