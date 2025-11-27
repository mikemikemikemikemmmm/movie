local value = ARGV[1]
local key_count = #KEYS
local released_count = 0
local seat_set_key = "seat_set"

-- 遍歷所有鍵，嘗試釋放
for i, key in ipairs(KEYS) do
    local current = redis.call("GET", key)
    if current == value then
        redis.call("DEL", key)
        released_count = released_count + 1
        redis.call("SREM", seat_set_key, key)
        redis.log(redis.LOG_NOTICE, "✅ 已釋放鍵: " .. key)
    else
        redis.log(redis.LOG_NOTICE, "⚠️ 未釋放鍵: " .. key .. "（值不匹配）")
    end
end

redis.log(redis.LOG_NOTICE, "🔓 總共釋放 " .. released_count .. " / " .. key_count .. " 個鍵。")
return "success"