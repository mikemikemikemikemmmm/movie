-- ARGV[1] = userID
local inputVal = ARGV[1]
local allProcessing = true
local allReserved = true

for i, key in ipairs(KEYS) do
    local val = redis.call("GET", key)
        redis.log(redis.LOG_NOTICE, "✅ val: " .. tostring(val))
    if not val then
        return "no val with key:" .. key .. " value:" .. tostring(val)
    end
    if val ~= inputVal then
        allProcessing = false
    end
    if val ~= inputVal .. ":reserved" then
        allReserved = false
    end
end

if allProcessing then
    return "processing"
elseif allReserved then
    return "success"
else
    return "error"
end