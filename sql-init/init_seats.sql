CREATE TABLE IF NOT EXISTS seats (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    x INT NOT NULL,
    y INT NOT NULL,
    user_id INT,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    CONSTRAINT unique_xy UNIQUE (x, y)  -- 確保 (x, y) 組合唯一
);
-- 生成 100 筆假資料，x 與 y 從 1 到 10
INSERT INTO seats (x, y)
SELECT x_series AS x,
       y_series AS y
FROM generate_series(1, 10) AS y_series
CROSS JOIN generate_series(1, 10) AS x_series
WHERE NOT (
    y_series BETWEEN 1 AND 5 
    AND 
    x_series IN (1, 2, 9, 10) -- 排除 x 是 1, 2, 9, 或 10, 且 y 在 1 到 5 之間
);