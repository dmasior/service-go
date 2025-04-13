-- USERS
INSERT INTO "user"
    (id, email, password_hash, last_login_at, created_at)
VALUES
    -- John
    ('U_ZWF85NWFH6L497VV3GQQM62HGGP40J5V', 'john@example.com', '$argon2id$v=19$m=65536,t=3,p=4$xj/jeKPHgItWS/djsur3Sg$9BGTObGgvZ+hC/tXZKGlrH9F/bpu31nCHYzMMBpNmG8', '2025-04-12 12:01:58.516 +0200', '2025-04-12 12:01:58.516 +0200'),
    -- Jane
    ('U_CXRH8EQYB2DCGN79CPH26N4A40WAUYS0', 'jane@example.com', '$argon2id$v=19$m=65536,t=3,p=4$iQoVShRzGyhfYilKJIDYog$qUqG5OHEe3nClXqpHKdFN5VkQKl2r/QAzcDNSS8ABik', '2025-04-12 11:01:58.516 +0200', '2025-04-12 12:01:58.516 +0200')
;
