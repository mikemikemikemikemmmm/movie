-- ARGV[1] 為鎖的值，ARGV[2] 為 TTL 秒數
local value = ARGV[1]
local ttl_ms = 10000  -- TTL 轉為毫秒
local seat_set_key = "seat_set"
local locked_keys = {}

-- 嘗試鎖定所有鍵
for i, key in ipairs(KEYS) do
    local ok = redis.call("SET", key, value, "NX", "PX", ttl_ms)
    if ok then
        table.insert(locked_keys, key)
        redis.log(redis.LOG_NOTICE, "✅ 鎖定成功: 鍵=" .. key .. " 值=" .. value .. " TTL(ms)=" .. ttl_ms)
    else
        redis.log(redis.LOG_NOTICE, "❌ 鎖定失敗: 鍵=" .. key .. "，開始回滾已鎖定鍵")
        -- 回滾已鎖定的鍵
        for _, lk in ipairs(locked_keys) do
            redis.call("DEL", lk)
            redis.log(redis.LOG_NOTICE, "♻️ 回滾釋放鎖鍵: " .. lk)
        end
        return "failed"
    end
end

-- 所有鍵鎖定成功後，加入集合
for _, key in ipairs(locked_keys) do
    redis.call("SADD", seat_set_key, key)
    redis.log(redis.LOG_NOTICE, "➕ 已加入集合: " .. seat_set_key .. " 鍵=" .. key)
end

redis.log(redis.LOG_NOTICE, "🎉 所有鍵鎖定並加入集合成功！總鎖定鍵數=" .. #locked_keys)
return "success"