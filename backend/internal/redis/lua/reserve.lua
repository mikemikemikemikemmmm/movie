local value = ARGV[1]
local keyCount = #KEYS
local ttl_ms = 7 * 24 * 60 * 60 * 1000  -- 一週的毫秒數

-- 1️⃣ 檢查所有 key 是否存在且 value 正確
for i, key in ipairs(KEYS) do
    local v = redis.call("GET", key)

    if not v then
        return cjson.encode({status="failed", key=key, value=nil, reason="key does not exist"})
    end

    if v ~= value then
        return cjson.encode({status="failed", key=key, value=v, reason="value mismatch"})
    end
end

-- 2️⃣ 所有 key 都存在且 value 正確 → 設定 TTL 一週
for i, key in ipairs(KEYS) do
    local newValue = value .. ":reserved"
    redis.call("SET", key, newValue)
    redis.call("PEXPIRE", key, ttl_ms)
end

return cjson.encode({status="success"})