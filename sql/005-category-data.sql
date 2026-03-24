-- product 1, 4, 7 - Clothing
INSERT INTO categories (product_id, name, code)
SELECT 
    id, 
    'Clothing', 
    CONCAT('CLO-ID-', id) 
FROM products
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

-- product 2, 6 - Shoes
INSERT INTO categories (product_id, name, code)
SELECT 
    id, 
    'Shoes', 
    CONCAT('SHO-ID-', id) 
FROM products
WHERE code IN ('PROD002', 'PROD006');

-- product 3, 5, 8 - Accessories
INSERT INTO categories (product_id, name, code)
SELECT 
    id, 
    'Accessories', 
    CONCAT('ACC-ID-', id) 
FROM products
WHERE code IN ('PROD003', 'PROD005', 'PROD008');